	package web

	import (
		"context"
		"log"
		"net/http"
		"os"
		"os/signal"
		"syscall"
		"time"
	)

	type Server struct {
		httpServer *http.Server
	}

	func NewServer(handler http.Handler) *Server {
		srv := &http.Server{
			Addr:         getAddr(),
			Handler:      handler,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
		}

		return &Server{
			httpServer: srv,
		}
	}

	func (s *Server) Start() error {
		// канал для ошибок сервера
		errCh := make(chan error, 1)

		// запуск сервера в горутине
		go func() {
			log.Printf("🚀 server started on http://%s\n", s.httpServer.Addr)
			errCh <- s.httpServer.ListenAndServe()
		}()

		// канал для graceful shutdown
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

		select {
		case err := <-errCh:
			return err

		case <-stop:
			log.Println("🛑 shutdown signal received")

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			return s.httpServer.Shutdown(ctx)
		}
	}

	func getAddr() string {
		if port := os.Getenv("PORT"); port != "" {
			return ":" + port
		}
		return ":8080"
	}