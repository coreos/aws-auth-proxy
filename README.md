#aws-auth-proxy

[![Docker Repository on Quay](https://quay.io/repository/coreos/aws-auth-proxy/status "Docker Repository on Quay")](https://quay.io/repository/coreos/aws-auth-proxy)

##Installation

```sh
#requires go1.5
export GO15VENDOREXPERIMENT=1
cd ./aws-auth-proxy
go build github.com/coreos/aws-auth-proxy
```
##Example

```sh
# aws elasticsearch example
./aws-auth-proxy \
-access-key=xxx \
-secret-key=xxxx \
-service-name=es \
-region-name=<your-aws-region> \
-upstream-host=<your-aws-elastic-search-endpoint> \
-upstream-scheme=https \
-listen-address=":9200"
```

Your proxied elasticsearch endpoint is now here: [http://localhost:9200](http://localhost:9200)


No more securing elastic search endpoints with IP addresses!



