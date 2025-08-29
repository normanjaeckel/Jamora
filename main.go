package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"golang.org/x/sys/unix"

	"github.com/normanjaeckel/Jamora/pkg/server"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}

func run() error {
	ctx, cancel := interruptContext()
	defer cancel()

	// Load config and database, then start the webserver

	if err := server.Run(ctx); err != nil {
		return fmt.Errorf("running server: %w", err)
	}
	return nil
}

// interruptContext listens on SIGINT and SIGTERM. If one of theses signals is
// received for the first time, the context is closed. If one is received for
// the second time, the process is killed with status code 1.
func interruptContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sigs := make(chan os.Signal, 1)
		// TODO: Use syscall.SIG... oder unix.SIG...
		signal.Notify(sigs, unix.SIGINT, unix.SIGTERM)
		<-sigs
		cancel()

		// If the signal is received for the second time, make a hard cut.
		<-sigs
		log.Fatal("Programm interrupted")
	}()
	return ctx, cancel
}
