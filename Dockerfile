FROM golang

RUN mkdir -p /go/src/github.com/simcap/ratpbot
ADD . /go/src/github.com/simcap/ratpbot

RUN go install github.com/simcap/ratpbot

WORKDIR /go/src/github.com/simcap/ratpbot

ENTRYPOINT /go/bin/ratpbot

EXPOSE 8080
