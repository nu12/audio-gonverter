FROM golang:1.20-alpine as builder

WORKDIR /app

COPY . .

RUN go build -o gonverter cmd/main.go  cmd/helpers.go cmd/routes.go

FROM alpine:latest
LABEL org.opencontainers.image.source https://github.com/nu12/audio-gonverter

ARG commit
ENV COMMIT=${commit}

WORKDIR /app

COPY --from=builder /app/gonverter /app/gonverter
COPY --from=builder /app/cmd/templates /app/cmd/templates
COPY --from=builder /app/cmd/static /app/cmd/static

ENTRYPOINT ["./gonverter"]