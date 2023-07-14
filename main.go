package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

func main() {

	// Define flag variables
	var certFile string
	var keyFile string
	var showHelp bool

	// Define flags and usage
	flag.StringVar(&certFile, "cert", "", "Path to the TLS certificate file")
	flag.StringVar(&keyFile, "key", "", "Path to the TLS private key file")
	flag.BoolVar(&showHelp, "help", false, "Show help message")

	// Set custom usage function
	flag.Usage = func() {
		flag.PrintDefaults()
	}

	// Parse command-line arguments
	flag.Parse()

	// Check if help flag is provided
	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	// Check if no flags were provided
	if flag.NFlag() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	// Create HTTP reverse proxy
	httpProxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			// Set the target URL to the original request URL
			req.URL.Scheme = "http"
			req.URL.Host = req.Host
		},
		ErrorHandler: func(rw http.ResponseWriter, req *http.Request, err error) {
			log.Println("Reverse proxy error:", err)
			http.Error(rw, "Oops! Something went wrong. Inspect server logs.", http.StatusInternalServerError)
		},
	}

	// Create HTTPS reverse proxy
	httpsProxy := &httputil.ReverseProxy{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Director: func(req *http.Request) {
			// Set the target URL to the original request URL
			req.URL.Scheme = "https"
			req.URL.Host = req.Host
		},
		ErrorHandler: func(rw http.ResponseWriter, req *http.Request, err error) {
			log.Println("Reverse proxy error:", err)
			http.Error(rw, "Oops! Something went wrong. Inspect server logs.", http.StatusInternalServerError)
		},
	}

	ingressHttpServer := &http.Server{
		Addr:    ":80",
		Handler: httpProxy,
	}

	apiServer := &http.Server{
		Addr: ":6443",
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{
				loadTLSCertificate(certFile, keyFile),
			},
		},
		Handler: httpsProxy,
	}

	// Configure the HTTPS server
	ingressHttpsServer := &http.Server{
		Addr: ":443",
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{
				loadTLSCertificate(certFile, keyFile),
			},
		},
		Handler: httpsProxy,
	}

	// Start the HTTPS server on port 6443
	go func() {
		log.Println("Starting reverse proxy server on port 6443 ...")
		err := apiServer.ListenAndServeTLS("", "")
		if err != nil {
			log.Fatal("Error starting reverse proxy server:", err)
		}
	}()

	// Start the HTTPS server on port 80
	go func() {
		log.Println("Starting reverse proxy server on port 80 ...")
		err := ingressHttpServer.ListenAndServe()
		if err != nil {
			log.Fatal("Error starting reverse proxy server:", err)
		}
	}()

	// Start the HTTPS server on port 443
	go func() {
		log.Println("Starting reverse proxy server on port 443...")
		err := ingressHttpsServer.ListenAndServeTLS("", "")
		if err != nil {
			log.Fatal("Error starting reverse proxy server:", err)
		}
	}()

	// Wait indefinitely to keep the program running
	select {}
}

// LoadTLSKeyPair loads a TLS certificate and private key from files and returns a tls.Certificate.
func loadTLSCertificate(certFile, keyFile string) tls.Certificate {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal("Error loading TLS certificate:", err)
	}
	return cert
}
