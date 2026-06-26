package httpserver

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"strconv"
)

type HTTPServer struct {
	Port    uint16
	Context context.Context
}

func New(ctx context.Context, port uint16) *HTTPServer {
	return &HTTPServer{port, ctx}
}

func (s *HTTPServer) Serve() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Handling incoming request", "method", r.Method, "path", r.URL.Path)
		fmt.Fprintf(w, "Hello!")
	})

	listenOn := net.JoinHostPort("0.0.0.0", strconv.Itoa(int(s.Port)))
	listener, err := net.Listen("tcp", listenOn)
	if err != nil {
		log.Print(err)
		return fmt.Errorf("Failed to bind on port %v: %w", s.Port, err)
	}

	log.Printf("HTTP server listening on %s", listener.Addr())

	go func() {
		<-s.Context.Done()
		listener.Close()
	}()

	return http.Serve(listener, mux)
}
