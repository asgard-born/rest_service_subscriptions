package rest_service

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	log.Printf("starting server on port %s", port)
	err := s.httpServer.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		log.Printf("server error: %v", err)
		return err
	}

	log.Println("server stopped gracefully")
	return err
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("shutting down server...")
	return s.httpServer.Shutdown(ctx)
}
