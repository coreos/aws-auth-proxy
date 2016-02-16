package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
	"net/http"
	"net/url"
	"strings"

	"github.com/coreos/pkg/flagutil"
	"github.com/goamz/goamz/aws"
)

func main() {
	var (
		auth          aws.Auth
		targetURL     url.URL
		listenAddress string
		serviceName   string
		region        aws.Region
	)
	const flagEnvPrefix = "AWS_AUTH_PROXY"

	fs := flag.NewFlagSet("aws-auth-proxy", flag.ExitOnError)

	fs.StringVar(&auth.AccessKey, "access-key", "", "aws access key id")
	fs.StringVar(&auth.SecretKey, "secret-key", "", "aws secret access key")
	fs.StringVar(&serviceName, "service-name", "", "aws service name")
	var regionName string
	fs.StringVar(&regionName, "region-name", "", "aws region name")
	fs.StringVar(&targetURL.Host, "upstream-host", "", "host or host:port for upstream endpoint")
	fs.StringVar(&targetURL.Scheme, "upstream-scheme", "https", "scheme for upstream endpoint")
	fs.StringVar(&listenAddress, "listen-address", ":8080", "address for proxy to listen on")

	if err := flagutil.SetFlagsFromEnv(fs, flagEnvPrefix); err != nil {
		log.Fatal(err)
	}

	if len(os.Args) >= 2 && os.Args[1] == "--help" {
		fs.PrintDefaults()
		fmt.Printf("\nflagutil prefix is '%s'\n", flagEnvPrefix)
		fmt.Printf("example:\n\t-access-key=xxx OR export %s_ACCESS_KEY=xxx\n", flagEnvPrefix)
		os.Exit(0)
	}
	fs.Parse(os.Args[1:])

	region = aws.Regions[regionName]

	signer := aws.NewV4Signer(auth, serviceName, region)

	proxyHandler := &AWSProxy{
		TargetURL: &targetURL,
		Signer:    signer,
	}
	fmt.Printf("Listening on %s\n", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, proxyHandler))

}

type AWSProxy struct {
	TargetURL *url.URL
	Signer    *aws.V4Signer
}

func copyHeaders(dst, src http.Header) {
	for k, vals := range src {
		for _, v := range vals {
			dst.Add(k, v)
		}
	}
}

// Escapes URL with building URL lib plus add comma
func EscapeURL(s string) string {
	parts := strings.Split(s,"?")
	path := &url.URL{Path: parts[0]}

	escapedUrl, _ := url.Parse(path.String())
	if (len(parts) > 1) {
		escapedUrl.RawQuery = parts[1]
	}

	return strings.Replace(escapedUrl.String(), ",", "%2C", -1)
}

func (h *AWSProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	respondError := func(err error) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	proxyURL := *r.URL
	proxyURL.Host = h.TargetURL.Host
	proxyURL.Scheme = h.TargetURL.Scheme

  signingURL := EscapeURL(proxyURL.String());
	signedReq, err := http.NewRequest(
		r.Method,
		signingURL,
		r.Body,
	)

	if err != nil {
		respondError(err)
		return
	}

	signedReq.Header.Set("host", proxyURL.Host)
	signedReq.Header.Set("x-amz-date", time.Now().UTC().Format(aws.ISO8601BasicFormat))
	h.Signer.Sign(signedReq)

	signedReq.URL, _ = url.Parse(proxyURL.String())

	resp, err := http.DefaultClient.Do(signedReq)
	if err != nil {
		respondError(err)
		return
	}
	defer resp.Body.Close()

	copyHeaders(w.Header(), resp.Header)

	buf := bytes.Buffer{}
	if _, err := io.Copy(&buf, resp.Body); err != nil {
		respondError(err)
		return
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(buf.Bytes())
}
