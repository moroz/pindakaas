package httpserver

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
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
		slog.Info("Handling incoming request", "method", r.Method, "path", r.URL.Path)

		subdomain, _, _ := strings.Cut(r.Host, ".")
		conn, ok := s.connRegistry.GetTunnelForSubdomain(subdomain)
		if !ok {
			http.Error(w, "Bad Gateway", http.StatusBadGateway)
			return
		}
		_ = conn

		host, port, err := parseRemoteAddrPort(r.RemoteAddr)
		if err != nil {
			log.Print("Failed to parse remote address and port: ", err)
		}

		channel, _, err := conn.OpenForwardingChannel(host, port)
		if err != nil {
			log.Print("Failed to open forwarding channel: ", err)
			http.Error(w, "Bad Gateway", http.StatusBadGateway)
			return
		}
		defer channel.Close()

		err = r.Write(channel)
		if err != nil {
			// TODO: Bad Gateway?
			log.Print("Failed to forward request: ", err)
			http.Error(w, "Bad Gateway", http.StatusBadGateway)
			return
		}

		resp, err := http.ReadResponse(bufio.NewReader(channel), r)
		if err != nil {
			// TODO: Internal Server Error? Transport error?
			log.Print("Failed to read response to forwarded request: ", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		for key, values := range resp.Header {
			for _, value := range values {
				resp.Header.Add(key, value)
			}
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
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
