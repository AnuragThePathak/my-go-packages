package signals

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Context() context.Context {
	ctx, stopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sig
		log.Println("Shutting down server...")

		// Shutdown signal with grace period of 10 seconds
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		stopCtx()
	}()
	return ctx
}