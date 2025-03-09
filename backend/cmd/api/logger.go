package main

import (
	"log"
	"os"
)

type LogLevel uint

const (
	FATAL LogLevel = iota
	ERROR
	WARN
	INFO
	DEBUG
	TRACE
)

type Logger interface {
	Level(LogLevel) *log.Logger
}

// Simple logger
type Slog struct {
	infoLogger  *log.Logger
	errLogger   *log.Logger
	debugLogger *log.Logger
	fatalLogger *log.Logger
}

func (sl *Slog) Level(lvl LogLevel) *log.Logger {
	switch lvl {
	case FATAL:
		return sl.fatalLogger
	case ERROR:
		return sl.errLogger
	case INFO:
		return sl.infoLogger
	case DEBUG:
		return sl.debugLogger
	}
	return nil
}

func NewSlogger() *Slog {
	return &Slog{
		debugLogger: log.New(os.Stdout, "DEBUG\t", log.Ldate|log.Ltime|log.Lshortfile),
		infoLogger:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errLogger:   log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		fatalLogger: log.New(os.Stderr, "FATAL\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
