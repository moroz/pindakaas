package httpserver

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"strings"

	"github.com/moroz/pindakaas/config"
	"github.com/moroz/pindakaas/registry"
)

type HTTPServer struct {
	connRegistry *registry.Registry
}

func New(connRegistry *registry.Registry) *HTTPServer {
	return &HTTPServer{
		connRegistry: connRegistry,
	}
}

func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		subdomain, _, _ := strings.Cut(r.Host, ".")
		conn, ok := s.connRegistry.GetConnectionForSubdomain(subdomain)
		if !ok {
			http.Error(w, "Bad Gateway", http.StatusBadGateway)
			return
		}
		_ = conn

		slog.Info("Handling incoming request", "method", r.Method, "path", r.URL.Path)
		fmt.Fprintf(w, "Hello!")
	})

	mux.ServeHTTP(w, r)
}

func (s *HTTPServer) ListenAndServe(ctx context.Context, port uint16) error {
	listenOn := config.FormatHostPort(port)
	listener, err := net.Listen("tcp", listenOn)
	if err != nil {
		return fmt.Errorf("Failed to bind on port %v: %w", port, err)
	}

	log.Printf("HTTP server listening on %s", listener.Addr())

	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	return http.Serve(listener, s)
}

func (s *HTTPServer) ListenAndServeTLS(ctx context.Context, port uint16, certFile, keyFile string) error {
	listenOn := config.FormatHostPort(port)
	listener, err := net.Listen("tcp", listenOn)
	if err != nil {
		return fmt.Errorf("Failed to bind on port %v: %w", port, err)
	}

	log.Printf("HTTPS server listening on %s", listener.Addr())

	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	return http.ServeTLS(listener, s, certFile, keyFile)
}
