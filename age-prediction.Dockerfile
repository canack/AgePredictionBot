FROM golang:1.18.3 as build

COPY . /app
WORKDIR /app/cmd/agebot

RUN go mod tidy && \
    CGO_ENABLED=0 go build -ldflags "-s -w" .

FROM gcr.io/distroless/static
COPY --from=build /app/cmd/* /app

WORKDIR /app
CMD ["/app/agebot"]
