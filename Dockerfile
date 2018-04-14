FROM golang:1.9

ADD . /go/src/github.com/walesey/dicelang
WORKDIR /go/src/github.com/walesey/dicelang

RUN go build website/server.go

CMD ./server