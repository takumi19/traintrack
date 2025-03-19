package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"traintrack/internal/chat"
	"traintrack/internal/database"
	"traintrack/internal/editor"
	"traintrack/internal/middleware"
)

const (
	defaultIdleTimeout    = 5 * time.Second
	defaultReadTimeout    = 3 * time.Second
	defaultWriteTimeout   = 5 * time.Second
	defaultShutdownPeriod = 30 * time.Second
)

func main() {
	addr := flag.String("addr", ":8090", "HTTP network address")
	dbUrl := flag.String("dsn", "postgres://takumi@localhost:5432/traintrackdb2", "Data source name")
	flag.Parse()

	logger := NewSlogger()

	db, err := database.New(*dbUrl, false)
	if err != nil {
		logger.Level(FATAL).Fatal(err)
	}

	a := &Api{
		db:   db,
		l:    logger,
		eHub: editor.NewHub(),
		cHub: chat.NewHub(),
	}

	go a.eHub.Run()
	go a.cHub.Run()

	server := &http.Server{
		Addr:         *addr,
		Handler:      middleware.LogRequests(a.authenticate(a.routes())),
		ErrorLog:     logger.Level(ERROR),
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		IdleTimeout:  defaultIdleTimeout,
	}

	// Graceful shutdown
	waitForShutdown := make(chan struct{})
	go func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sc

		if err := server.Shutdown(context.Background()); err != nil {
			logger.Level(FATAL).Fatal("Server shutdown failed")
		}

		a.db.Close()
		close(waitForShutdown)
	}()

	logger.Level(INFO).Printf("Starting the server on %s", *addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Level(FATAL).Fatalf("Server shutdown failed:%s", err)
	}

	<-waitForShutdown
	logger.Level(INFO).Printf("Server shut down successfully!")
}
