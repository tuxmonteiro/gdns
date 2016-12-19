FROM golang:1.6-wheezy

MAINTAINER Marcelo Teixeira Monteiro

RUN mkdir -p /go/src/github.com/tuxmonteiro
WORKDIR /go/src/github.com/tuxmonteiro/gdns

ADD ./src/github.com/tuxmonteiro/gdns /go/src/github.com/tuxmonteiro/gdns
RUN go-wrapper download
RUN go-wrapper install

# this will ideally be built by the ONBUILD below ;)
CMD ["go-wrapper", "run"]
