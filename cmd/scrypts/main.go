package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"scrypts/internal/auth"
	"scrypts/internal/config"
	"scrypts/internal/middleware"
	"scrypts/internal/notes"
	"scrypts/internal/storage"
	"time"
)

func registerHandlers() {
	// Create rate limiter: 10 requests per minute
	rateLimiter := middleware.NewRateLimiter(10, time.Minute)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Scrypts is alive and kicking")
	})

	// Apply rate limiting to authentication endpoints
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		rateLimiter.RateLimit(http.HandlerFunc(auth.RegisterHandler)).ServeHTTP(w, r)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		rateLimiter.RateLimit(http.HandlerFunc(auth.LoginHandler)).ServeHTTP(w, r)
	})

	http.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			notes.CreateNoteHandler(w, r)
		case http.MethodGet:
			notes.GetNotesHandler(w, r)
		case http.MethodPut:
			notes.UpdateNoteHandler(w, r)
		case http.MethodDelete:
			notes.DeleteNoteHandler(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
}

func main() {
	// initialize config (loads secrets)
	config.Init()

	//initialize db

	if err := storage.Init("./scrypts.db"); err != nil {
		fmt.Println("Failed to init db: ", err)
		return
	}
	defer storage.Close()

	registerHandlers()

	certPath := os.Getenv("SCRYPTS_TLS_CERT")
	keyPath := os.Getenv("SCRYPTS_TLS_KEY")
	httpsPort := os.Getenv("SCRYPTS_HTTPS_PORT")
	if httpsPort == "" {
		httpsPort = "8443"
	}
	httpPort := os.Getenv("SCRYPTS_HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	if certPath != "" && keyPath != "" {
		tlsConfig := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
			CurvePreferences:         []tls.CurveID{tls.CurveP256, tls.X25519},
		}

		// Chain middleware: SecurityHeaders -> CORS -> DefaultServeMux
		handler := middleware.SecurityHeaders(middleware.CORS(http.DefaultServeMux))

		httpsSrv := &http.Server{
			Addr:         ":" + httpsPort,
			Handler:      handler,
			TLSConfig:    tlsConfig,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		}

		// start HTTPS server in background
		go func() {
			log.Printf("Starting HTTPS on :%s", httpsPort)
			if err := httpsSrv.ListenAndServeTLS(certPath, keyPath); err != nil && err != http.ErrServerClosed {
				log.Fatalf("HTTPS server failed: %v", err)
			}
		}()

		// redirect handler - sends clients to the HTTPS endpoint (preserves path)
		redirectHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host := r.Host
			if _, _, err := net.SplitHostPort(host); err == nil {
				h, _, _ := net.SplitHostPort(host)
				host = h
			}
			target := "https://" + host
			if httpsPort != "443" {
				target += ":" + httpsPort
			}
			target += r.URL.RequestURI()
			http.Redirect(w, r, target, http.StatusMovedPermanently)
		})

		// run the HTTP redirector in the foreground so the process stays alive
		log.Printf("Starting HTTP redirector on :%s -> https://:%s", httpPort, httpsPort)
		if err := http.ListenAndServe(":"+httpPort, redirectHandler); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Redirect server failed: %v", err)
		}
	}

	// If no TLS cert/key provided we fall back to plain HTTP (blocking)
	log.Printf("Starting server on http://localhost:%s", httpPort)
	if certPath == "" || keyPath == "" {
		// Chain middleware: SecurityHeaders -> CORS -> DefaultServeMux
		handler := middleware.SecurityHeaders(middleware.CORS(http.DefaultServeMux))
		if err := http.ListenAndServe(":"+httpPort, handler); err != nil {
			fmt.Println("Failed to start HTTP server:", err)
		}
	}
}