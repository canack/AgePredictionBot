package handler

import (
	"errors"
	"fmt"
	"github.com/canack/yasbot/recognize"
	tele "gopkg.in/telebot.v3"
	"log"
)

type MessageHandler struct {
}

func Welcome(c tele.Context) error {
	c.Notify(tele.Typing)
	c.Send(`Hi!
	I can try predict your age range with a photo
	Let's send a photo or reply a photo with /predict tag.

	If you liked this bot, feel free to contribute!
	Github: github.com/canack/AgePredictionBot`)
	return nil
}

func Predict(c tele.Context) error {
	if !c.Message().IsReply() {
		return c.Reply("You should reply an image with this command")
	}

	return ProcessImage(c)

}

// Returns message type
func processMessageOwner(c tele.Context) *tele.Message {
	if c.Message().IsReply() {
		return c.Message().ReplyTo
	}
	return c.Message()
}

func ProcessImage(c tele.Context) error {
	// Processing original message or replied message
	message := processMessageOwner(c)

	// Maximum size as KB
	maxDocSize := 10240
	//
	maxDocSize *= 1024

	// A little limit to prevent memory overflow
	docSize := message.Media().MediaFile().FileSize
	if docSize > maxDocSize {
		return c.Reply(fmt.Sprintf("Sorry, i can handle only maximum %dMB image", maxDocSize/1024/1024))
	}

	// Grabbing requested file
	requestFile, errFile := c.Bot().File(message.Media().MediaFile())
	defer requestFile.Close()

	if errFile != nil {
		return c.Reply("I can't reach this photo. Please check permissions")
	}

	c.Notify(tele.Typing)

	// Processing image
	reader, err := recognize.ProcessImage(requestFile)

	if errors.Is(err, recognize.NoFaceFound) {
		log.Println("No faces detected")
		return c.Reply("No faces detected")
	} else if errors.Is(err, recognize.ServerError) {
		log.Println("An error occurred with recognition service.", err)
		return c.Reply("An error occurred on recognition service. Please try later or try send a different photo.")
	} else if errors.Is(err, recognize.EncodeError) || errors.Is(err, recognize.DecodeError) {
		log.Println("Unsupported file type")
		return c.Reply("Unsupported file type or broken image uploaded. Please try with a different photo.")
	}

	c.Notify(tele.UploadingPhoto)

	// If requested image sent as document (without compress), response will be sent as a document (no compressing)
	if message.Media().MediaType() == "document" {
		document := &tele.Document{File: tele.FromReader(reader), FileName: "processed.jpg"}
		return c.Reply(document)
	}

	// Already commented above =)
	photo := &tele.Photo{File: tele.FromReader(reader)}
	return c.Reply(photo)
}

func ProcessText(c tele.Context) error {
	c.Notify(tele.Typing)

	return c.Send("I didn't catch that! Please /start how to use.")
}
