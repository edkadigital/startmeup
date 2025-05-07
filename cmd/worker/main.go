package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/edkadigital/startmeup/pkg/log"
	"github.com/edkadigital/startmeup/pkg/services"
	"github.com/edkadigital/startmeup/pkg/tasks"
)

func main() {
	// Start a new container.
	c := services.NewContainer()
	defer func() {
		// Gracefully shutdown all services.
		fatal("shutdown failed", c.Shutdown())
	}()

	// Register all task workers.
	tasks.Register(c)

	// Start the worker to process tasks from queues.
	log.Default().Info("Starting task worker")

	// Start the River worker in the foreground
	c.Tasks.Start(context.Background())

	// Wait for interrupt signal to gracefully shut down the task runner.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Default().Info("Shutting down worker")
}

// fatal logs an error and terminates the application, if the error is not nil.
func fatal(msg string, err error) {
	if err != nil {
		log.Default().Error(msg, "error", err)
		os.Exit(1)
	}
}
