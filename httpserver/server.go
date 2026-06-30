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
	"github.com/moroz/pindakaas/types"
	"github.com/moroz/pindakaas/web/handlers"
)

type HTTPServer struct {
	connRegistry types.TunnelRegistry
	adminRouter  http.Handler
}

func New(props *handlers.RouterProps) *HTTPServer {
	return &HTTPServer{
		connRegistry: props.TunnelRegistry,
		adminRouter:  handlers.Router(props),
	}
}

func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := r.Host

	if h, _, err := net.SplitHostPort(r.Host); err == nil {
		host = h
	}

	switch host {
	case config.BaseDomain, "localhost":
		s.adminRouter.ServeHTTP(w, r)
	default:
		s.ServeReverseProxy(w, r)
	}
}

func (s *HTTPServer) ServeReverseProxy(w http.ResponseWriter, r *http.Request) {
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
