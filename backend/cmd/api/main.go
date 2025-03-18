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

	slog := NewSlogger()

	db, err := database.New(*dbUrl)
	if err != nil {
		slog.Level(FATAL).Fatal(err)
	}

	a := &Api{
		db:   db,
		l:    slog,
		eHub: editor.NewHub(),
    cHub: chat.NewHub(),
	}

	go a.eHub.Run()
	go a.cHub.Run()

	server := &http.Server{
		Addr:         *addr,
		Handler:      a.routes(),
		ErrorLog:     slog.Level(ERROR),
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
			slog.Level(FATAL).Fatal("Server shutdown failed")
		}

		a.db.Close()
		close(waitForShutdown)
	}()

	slog.Level(INFO).Printf("Starting the server on %s", *addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Level(FATAL).Fatalf("Server shutdown failed:%s", err)
	}

	<-waitForShutdown
	slog.Level(INFO).Printf("Server shut down successfully!")
}
