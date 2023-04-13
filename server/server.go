package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type HTTPServer struct {
	srv *http.Server
	mux *http.ServeMux
}

func New(addr string) *HTTPServer {
	mux := http.NewServeMux()

	server := http.Server{Addr: addr, Handler: mux}

	return &HTTPServer{srv: &server, mux: mux}
}

func (s *HTTPServer) Register(endpoint string, handler http.Handler) {
	s.mux.Handle(endpoint, handler)
}

func (s *HTTPServer) Handler() http.Handler {
	return s.srv.Handler
}

func (s *HTTPServer) Start(ctx context.Context) error {
	serverErrors := make(chan error)

	go func() {
		log.Printf("Starting the server on %s\n", s.srv.Addr)
		err := s.srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	// Wait for context cancellation or server error
	select {
	case <-ctx.Done():
		log.Println("Shutting down the server...")
	case err := <-serverErrors:
		log.Fatal("Error starting HTTP server:", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server
	err := s.srv.Shutdown(shutdownCtx)
	if err != nil {
		return fmt.Errorf("could not gracefully shutdown the server: %w", err)
	}
	log.Println("Server stopped")

	return nil
}
