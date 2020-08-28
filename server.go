package gutils

import (
	"crypto/tls"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
)

func CreateTLSServer(domainToAllow string) *http.Server {
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domainToAllow), //Your domain here
		Cache:      autocert.DirCache("certs"),            //Folder for storing certificates
		Email: "noreply@example.com",
	}

	server := &http.Server{
		Addr: ":https",
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
			ServerName: domainToAllow,
		},
	}

	return server
}
