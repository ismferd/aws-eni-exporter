FROM golang:alpine as builder
ADD src/ /build/
ENV GOPATH=/build/go
ENV GOBIN=$GOPATH/bin
WORKDIR /build/
RUN env
RUN \
    apk add --no-cache \
        ca-certificates \
        gcc \
        git \
        musl-dev \
    && go get . \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o eni_exporter .

EXPOSE 8888

ENTRYPOINT ["./eni_exporter"]