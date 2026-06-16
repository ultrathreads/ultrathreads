package server

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Config struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	MaxHeaderBytes  int
}

type Server struct {
	httpServer      *http.Server
	shutdownTimeout time.Duration
}

func NewServer(cfg Config, handler http.Handler) *Server {
	if cfg.ShutdownTimeout == 0 {
		cfg.ShutdownTimeout = 10 * time.Second
	}
	if cfg.MaxHeaderBytes == 0 {
		cfg.MaxHeaderBytes = 1 << 20
	}
	return &Server{
		shutdownTimeout: cfg.ShutdownTimeout,
		httpServer: &http.Server{
			Addr:           ":" + cfg.Port,
			Handler:        handler,
			ReadTimeout:    cfg.ReadTimeout,
			WriteTimeout:   cfg.WriteTimeout,
			MaxHeaderBytes: cfg.MaxHeaderBytes,
		},
	}
}

// Start launches the server in a background goroutine.
// Returns an error channel that receives fatal listen errors.
func (s *Server) Start() <-chan error {
	errCh := make(chan error, 1)
	go func() {
		fmt.Printf("🚀 UltraThreads starting on %s\n", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("server listen failed: %w", err)
		}
		close(errCh)
	}()
	return errCh
}

// Stop gracefully shuts down the server within the configured timeout.
func (s *Server) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.shutdownTimeout)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		fmt.Printf("❌ HTTP server forced to shutdown: %v\n", err)
		return err
	}
	fmt.Println("✅ HTTP server stopped gracefully")
	return nil
}

// ShutdownTimeout returns the configured shutdown duration for external use.
func (s *Server) ShutdownTimeout() time.Duration {
	return s.shutdownTimeout
}