package httpserver

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
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

func parseRemoteAddrPort(addr string) (string, uint32, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", 0, err
	}

	parsedPort, err := strconv.ParseUint(port, 10, 32)
	if err != nil {
		return "", 0, err
	}

	return host, uint32(parsedPort), nil
}

func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Handling incoming request", "method", r.Method, "path", r.URL.Path, "host", r.Host)

		subdomain, _, found := strings.Cut(r.Host, ".")
		if !found {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		conn, ok := s.connRegistry.GetTunnelForSubdomain(subdomain)
		if !ok {
			http.Error(w, "Bad Gateway", http.StatusBadGateway)
			return
		}

		proxy := &httputil.ReverseProxy{
			Transport: conn,
			Rewrite: func(pr *httputil.ProxyRequest) {
				addr := net.JoinHostPort(conn.BindAddr, strconv.Itoa(int(conn.BindPort)))
				pr.Out.Host = addr
			},
		}

		proxy.ServeHTTP(w, r)
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
