package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)

func Run(ctx context.Context) error {
	addr := ":8080"
	handler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello, world!\n")
	}
	http.HandleFunc("/hello", handler)
	srv := &http.Server{
		Addr: addr,
		//BaseContext: func(net.Listener) context.Context { return ctx },
	}

	// Shutdown logic in a separate goroutine.
	shutdownResult := make(chan error)
	go func() {
		// Wait for the context to be closed.
		<-ctx.Done()

		// Shutdown server
		log.Print("Shutdown server")
		if err := srv.Shutdown(context.WithoutCancel(ctx)); err != nil {
			shutdownResult <- fmt.Errorf("server shutdown: %w", err)
			return
		}
		shutdownResult <- nil
	}()

	// Start server
	log.Printf("Webserver is listening on %s\n", addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("listen and serve: %w", err)
	}
	return <-shutdownResult
}
