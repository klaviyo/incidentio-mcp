package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/incident-io/incidentio-mcp-golang/internal/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received interrupt signal, shutting down gracefully...")
		cancel()
	}()

	// Use the server from internal/server package - single source of truth
	srv := server.New()
	if err := srv.Start(ctx); err != nil {
		log.Printf("Server error: %v", err)
	}
}
