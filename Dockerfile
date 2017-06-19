FROM golang:alpine

ENV GOPATH=/go

RUN apk add --no-cache curl git

RUN curl https://glide.sh/get | sh

RUN mkdir -p /go/src/github.com/BlooperDB/API
WORKDIR /go/src/github.com/BlooperDB/API

COPY glide.yaml /go/src/github.com/BlooperDB/API/
COPY glide.lock /go/src/github.com/BlooperDB/API/
RUN glide install

COPY . /go/src/github.com/BlooperDB/API/

WORKDIR /go

CMD ["go", "run", "src/github.com/BlooperDB/API/cmd/blooperapi/main.go"]