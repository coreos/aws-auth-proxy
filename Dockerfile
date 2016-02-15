FROM golang:1.5.2
MAINTAINER colin.hom@coreos.com

ENV GO15VENDOREXPERIMENT=1

RUN mkdir -p $GOPATH/src/github.com/coreos

WORKDIR $GOPATH/src/github.com/coreos/aws-auth-proxy

ADD . $GOPATH/src/github.com/coreos/aws-auth-proxy
RUN go get github.com/coreos/pkg/flagutil github.com/goamz/goamz/aws
RUN go install github.com/coreos/aws-auth-proxy

ENTRYPOINT ["aws-auth-proxy"]
