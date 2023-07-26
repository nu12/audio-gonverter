FROM golang:1.20
LABEL org.opencontainers.image.source https://github.com/nu12/audio-gonverter

RUN go install github.com/caarlos0/svu@v1.10.2

RUN curl https://get.docker.com | bash

RUN wget https://github.com/aquasecurity/trivy/releases/download/v0.43.1/trivy_0.43.1_Linux-64bit.deb \
 && dpkg -i trivy_0.43.1_Linux-64bit.deb \
 && rm trivy_0.43.1_Linux-64bit.deb
