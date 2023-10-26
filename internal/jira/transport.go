package jira

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/mk-5/fjira/internal/app"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	defaultHttpTransport = createHttpTransport()
)

func createHttpTransport() *http.Transport {
	t := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	if f := os.Getenv("SSL_CERT_FILE"); f != "" {
		t.TLSClientConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
		data, err := os.ReadFile(f)
		if err != nil {
			app.Error(fmt.Sprintf("Cannot read file for SSL_CERT_FILE. %s", err.Error()))
			app.GetApp().Quit()
			return t
		}
		rootCAs := systemCertPool()
		rootCAs.AppendCertsFromPEM(data)
		t.TLSClientConfig.RootCAs = rootCAs
	}
	return t
}

func systemCertPool() *x509.CertPool {
	pool, err := x509.SystemCertPool()
	if err != nil {
		return x509.NewCertPool()
	}
	return pool
}
