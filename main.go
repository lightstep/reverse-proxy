package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

const (
	timeout             = 30 * time.Second
	defaultSatelliteURL = "https://ingest.lightstep.com/"
	defaultPort         = "8126"
)

func main() {
	forwardURL := flag.String("forward-url", defaultSatelliteURL, "Satellite address")
	proxyPort := flag.String("port", defaultPort, "Port for the proxy server")
	caFile := flag.String("ca", "", "Custom CA certificate")
	flag.Parse()

	satelliteURL, err := url.Parse(*forwardURL)
	if err != nil {
		log.Fatalf("Invalid forwarding address: %s, err: %s", forwardURL, err)
	}

	tlsConfig, err := getTLSConfig(*caFile)
	if err != nil {
		log.Fatal(err)
	}
	transport := createTransport(tlsConfig)
	proxy := httputil.NewSingleHostReverseProxy(satelliteURL)
	proxy.Transport = transport

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})
	log.Println("Listening for traffic")
	log.Fatal(http.ListenAndServe(":"+*proxyPort, nil))
}

// getTLSConfig returns a *tls.Config according to whether a user has supplied a customCACertFile. If they have,
// we return a TLSConfig that uses the custom CA cert as the lone Root CA. If not, we return nil which http.Transport
// will interpret as the default system defined Root CAs.
func getTLSConfig(customCACertFile string) (*tls.Config, error) {
	if len(customCACertFile) == 0 {
		return nil, nil
	}

	caCerts := x509.NewCertPool()
	cert, err := ioutil.ReadFile(customCACertFile)
	if err != nil {
		return nil, err
	}

	if !caCerts.AppendCertsFromPEM(cert) {
		return nil, fmt.Errorf("credentials: failed to append certificate")
	}

	return &tls.Config{RootCAs: caCerts}, nil
}

func createTransport(tlsClientConfig *tls.Config) *http.Transport {
	// Use a transport independent from http.DefaultTransport to provide sane
	// ch defaults that make sense in the context of the lightstep client. The
	// differences are mostly on setting timeouts based on the report timeout
	// and period.
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   timeout / 2,
			DualStack: true,
		}).DialContext,
		// The collector responses are very small, there is no point asking for
		// a compressed payload, explicitly disabling it.
		DisableCompression:     true,
		IdleConnTimeout:        4 * timeout,
		TLSHandshakeTimeout:    timeout / 2,
		ResponseHeaderTimeout:  timeout,
		ExpectContinueTimeout:  timeout,
		MaxResponseHeaderBytes: 64 * 1024, // 64 KB, just a safeguard
		TLSClientConfig:        tlsClientConfig,
	}
}
