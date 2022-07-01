package recognize

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
	"github.com/fogleman/gg"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
)

var NoFaceFound = errors.New("recognized 0 faces")
var EncodeError = errors.New("encode unsuccessful")
var DecodeError = errors.New("decode unsuccessful")
var ServerError = errors.New("connection error with aws recognition service")

func dragRectangle(drw *gg.Context, x0, y0, x1, y1 float64) {
	drw.DrawRectangle(x0, y0, x1-x0, y1-y0)
	drw.SetLineWidth(2)
	drw.SetHexColor("#f54281")
	drw.StrokePreserve()
	drw.SetRGBA(0, 0, 0, 0.5)
	drw.Fill()
}

func drawAge(drw *gg.Context, ageMin, ageMax int32, x, y, fontsize float64) {
	drw.LoadFontFace("VCR_OSD_MONO_1.001.ttf", fontsize)
	drw.SetHexColor("#02db9a")
	drw.DrawString(fmt.Sprintf("%d-%d", ageMin, ageMax), x, y)
}

var Client *rekognition.Client

func SetupRekognition() error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}
	Client = rekognition.NewFromConfig(cfg)
	return nil

}

func ProcessImage(input io.Reader) (io.ReadCloser, error) {

	// To save images bounds.
	var x, y int
	img, _, err := image.Decode(input)
	if err != nil {
		return io.ReadCloser(nil), DecodeError
	}

	x = img.Bounds().Dx()
	y = img.Bounds().Dy()

	// Image encoding before sent to aws rekognition service
	imgBytes := new(bytes.Buffer)

	// Todo: Quality should be dynamic
	if err := jpeg.Encode(imgBytes, img, &jpeg.Options{Quality: 80}); err != nil {
		return io.ReadCloser(nil), EncodeError
	}

	// Rekognition service's request
	out, err := Client.DetectFaces(context.TODO(), &rekognition.DetectFacesInput{
		Image:      &types.Image{Bytes: imgBytes.Bytes()},
		Attributes: []types.Attribute{types.AttributeAll},
	})

	if err != nil {
		return io.ReadCloser(nil), ServerError
	}

	if len(out.FaceDetails) == 0 {
		return io.ReadCloser(nil), NoFaceFound
	}

	// To draw rectangle and write age as text, used another library currently
	drw := gg.NewContextForImage(img)

	// To processing all faces
	for _, v := range out.FaceDetails {
		bb := v.BoundingBox
		x0 := float64(*bb.Left * float32(x))
		y0 := float64(*bb.Top * float32(y))
		x1 := float64(*bb.Width*float32(x) + *bb.Left*float32(x))
		y1 := float64(*bb.Height*float32(y) + *bb.Top*float32(y))

		dragRectangle(drw, x0, y0, x1, y1)

		// For dynamically font size
		r := x1 - x0
		fs := r / math.Sqrt(1.8*math.Sqrt(r))
		drawAge(drw, *v.AgeRange.Low, *v.AgeRange.High, x0+2, y1-2, fs)

	}

	// Encoding
	buf := new(bytes.Buffer)

	if err := jpeg.Encode(buf, drw.Image(), &jpeg.Options{Quality: 90}); err != nil {
		return io.ReadCloser(nil), DecodeError
	}

	return io.NopCloser(buf), nil
}
