FROM golang:alpine

ENV GOPATH=/go

RUN apk add --no-cache curl git

RUN go get -u github.com/golang/dep/cmd/dep

RUN mkdir -p /go/src/github.com/BlooperDB/API
WORKDIR /go/src/github.com/BlooperDB/API

COPY . /go/src/github.com/BlooperDB/API/
RUN dep ensure

WORKDIR /go

CMD ["go", "run", "src/github.com/BlooperDB/API/cmd/blooperapi/main.go"]