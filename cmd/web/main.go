package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/edkadigital/startmeup/pkg/handlers"
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

	// Build the router.
	if err := handlers.BuildRouter(c); err != nil {
		fatal("failed to build the router", err)
	}

	// Register all task queues to allow web server to enqueue tasks.
	tasks.Register(c)

	// Start the server.
	go func() {
		srv := http.Server{
			Addr:         fmt.Sprintf("%s:%d", c.Config.HTTP.Hostname, c.Config.HTTP.Port),
			Handler:      c.Web,
			ReadTimeout:  c.Config.HTTP.ReadTimeout,
			WriteTimeout: c.Config.HTTP.WriteTimeout,
			IdleTimeout:  c.Config.HTTP.IdleTimeout,
		}

		if c.Config.HTTP.TLS.Enabled {
			certs, err := tls.LoadX509KeyPair(c.Config.HTTP.TLS.Certificate, c.Config.HTTP.TLS.Key)
			fatal("cannot load TLS certificate", err)

			srv.TLSConfig = &tls.Config{
				Certificates: []tls.Certificate{certs},
			}
		}

		if err := c.Web.StartServer(&srv); errors.Is(err, http.ErrServerClosed) {
			fatal("shutting down the server", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the web server and task runner.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)
	<-quit
}

// fatal logs an error and terminates the application, if the error is not nil.
func fatal(msg string, err error) {
	if err != nil {
		log.Default().Error(msg, "error", err)
		os.Exit(1)
	}
}
