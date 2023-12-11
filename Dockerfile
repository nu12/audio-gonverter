FROM golang:1.20-alpine as builder

WORKDIR /app

COPY . .

RUN go build -o gonverter cmd/main.go  cmd/helpers.go cmd/routes.go cmd/handlers.go

FROM alpine:3.19.0
LABEL org.opencontainers.image.source https://github.com/nu12/audio-gonverter

ARG commit
ENV COMMIT=${commit}

WORKDIR /app

RUN apk add --no-cache ffmpeg

COPY --from=builder /app/gonverter /app/gonverter
COPY --from=builder /app/cmd/templates /app/cmd/templates
COPY --from=builder /app/cmd/static /app/cmd/static

VOLUME [ "/tmp/original", "/tmp/converted" ]

EXPOSE 8080 9000

ENTRYPOINT ["./gonverter"]