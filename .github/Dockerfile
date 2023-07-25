FROM golang:1.20
LABEL org.opencontainers.image.source https://github.com/nu12/audio-gonverter

# svu

RUN wget https://github.com/aquasecurity/trivy/releases/download/v0.43.1/trivy_0.43.1_Linux-64bit.deb \
 && dpkg -i trivy_0.43.1_Linux-64bit.deb \
 && rm trivy_0.43.1_Linux-64bit.deb
