FROM golang:1.10.1
MAINTAINER cohom@redhat.com

WORKDIR $GOPATH/src/github.com/coreos
ADD . ./aws-auth-proxy
WORKDIR ./aws-auth-proxy
RUN go install github.com/coreos/aws-auth-proxy

ENTRYPOINT ["aws-auth-proxy"]
