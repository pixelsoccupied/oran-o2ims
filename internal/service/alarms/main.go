package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/openshift-kni/oran-o2ims/internal/service/alarms/api/generated"
	"github.com/openshift-kni/oran-o2ims/internal/service/alarms/internal/middleware"
	"github.com/openshift-kni/oran-o2ims/internal/service/alarms/internal/server"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	Serve()
}

// Serve TODO: Call this func using cobra-cli from inside deployment CR.
func Serve() {
	// TODO: Init client-go

	// TODO: Init DB client

	// TODO: Audit and Insert data database

	// TODO: Launch k8s job for DB remove archived data

	// Init server
	r := http.NewServeMux()

	// Register
	alarmServerStrict := generated.NewStrictHandler(server.AlarmsServer{}, nil)

	// Create the handler
	handler := generated.HandlerFromMux(alarmServerStrict, r)

	// Add all the middlewares here
	mwStack := middleware.CreateMwStack(
		middleware.AlarmsOapiValidation(),
		middleware.LogDuration(),
	)

	// Server config
	srv := &http.Server{
		Handler:      mwStack(handler),
		Addr:         net.JoinHostPort("0.0.0.0", "8080"),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		ErrorLog:     log.New(os.Stdout, "ALARMS-SERVER: ", log.Ldate|log.Ltime|log.Lshortfile),
	}

	// Channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)
	// Channel for shutdown signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Start server
	go func() {
		log.Printf("Server is starting on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	// Blocking select
	select {
	case err := <-serverErrors:
		fmt.Printf("error starting server: %v\n", err.Error())
	case sig := <-shutdown:
		fmt.Printf("Shutdown signal received: %v\n", sig)
		if err := gracefulShutdown(srv); err != nil {
			fmt.Printf("Graceful shutdown failed: %v", err)
		}
	}
}

func gracefulShutdown(srv *http.Server) error {
	// Create shutdown context with 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		return err
	}
	log.Println("Server gracefully stopped")
	return nil
}
