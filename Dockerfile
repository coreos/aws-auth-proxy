FROM golang:1.5.2
MAINTAINER colin.hom@coreos.com

ENV GO15VENDOREXPERIMENT=1

WORKDIR $GOPATH/src/github.com/coreos
ADD . ./aws-auth-proxy
WORKDIR ./aws-auth-proxy
RUN go install github.com/coreos/aws-auth-proxy

ENTRYPOINT ["aws-auth-proxy"]