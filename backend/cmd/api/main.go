package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"traintrack/internal/chat"
	"traintrack/internal/database"
	"traintrack/internal/editor"
	"traintrack/internal/middleware"

	"github.com/lmittmann/tint"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

const (
	defaultIdleTimeout    = 5 * time.Second
	defaultReadTimeout    = 3 * time.Second
	defaultWriteTimeout   = 5 * time.Second
	defaultShutdownPeriod = 30 * time.Second
)

type config struct {
	baseURL   string
	httpPort  int
	basicAuth struct {
		username       string
		hashedPassword string
	}
	db struct {
		dsn         string
		automigrate bool
	}
	jwt struct {
		secretKey string
	}
}

func main() {
	var cfg config

	// flag.StringVar(&cfg.baseURL, "base-url", "http://localhost:8090", "base URL for the application")
	flag.IntVar(&cfg.httpPort, "http-port", 8090, "port to listen on for HTTP requests")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://takumi@localhost:5432/traintrackdb2?sslmode=disable", "Database DSN")
	flag.BoolVar(&cfg.db.automigrate, "db-automigrate", false, "run migrations on startup")
	flag.StringVar(&cfg.jwt.secretKey, "jwt-secret-key", "to6u2ro7ibzghvsp5h32ihoyi7v3oizk", "secret key for JWT authentication")

	flag.Parse()

	logger := NewSlogger()

	db, err := database.New(cfg.db.dsn, cfg.db.automigrate)
	if err != nil {
		logger.Level(FATAL).Fatal(err)
	}

	a := &Api{
		db:   db,
		l:    logger,
		c:    cfg,
		eHub: editor.NewHub(),
		cHub: chat.NewHub(),
	}

	go a.eHub.Run()
	go a.cHub.Run()

	sl := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))

	authDeps := middleware.AuthDeps{
		DB:        db,
		Logger:    sl,
		JwtSecret: cfg.jwt.secretKey,
	}

	middlewares := middleware.Chain(
		middleware.LogRequests,
		middleware.AuthMiddleware(authDeps),
	)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", a.c.httpPort),
		Handler:      middlewares(a.routes()),
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

	logger.Level(INFO).Printf("Starting the server on %d", a.c.httpPort)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Level(FATAL).Fatalf("Server shutdown failed:%s", err)
	}

	<-waitForShutdown
	logger.Level(INFO).Printf("Server shut down successfully!")
}
