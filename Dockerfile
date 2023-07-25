FROM golang:1.20-alpine as builder

WORKDIR /app

COPY . .

RUN go build -o gonverter cmd/main.go  cmd/helpers.go 

FROM alpine:latest
LABEL org.opencontainers.image.source https://github.com/nu12/audio-gonverter

WORKDIR /app

COPY --from=builder /app/gonverter /app/gonverter

ENTRYPOINT ["./gonverter"]