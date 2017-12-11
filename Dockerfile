FROM golang:alpine

ENV GOPATH=/go

RUN apk add --no-cache curl git

RUN go get github.com/tools/godep

RUN mkdir -p /go/src/github.com/BlooperDB/API
WORKDIR /go/src/github.com/BlooperDB/API

COPY Gopkg.lock /go/src/github.com/BlooperDB/API/
COPY Gopkg.toml /go/src/github.com/BlooperDB/API/
RUN godep ensure

COPY . /go/src/github.com/BlooperDB/API/

WORKDIR /go

CMD ["go", "run", "src/github.com/BlooperDB/API/cmd/blooperapi/main.go"]