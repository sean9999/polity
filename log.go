package main

import (
	"io"

	"github.com/charmbracelet/log"
)

func Slog(w io.Writer) *log.Logger {
	logger := log.New(w)
	logger.SetReportTimestamp(false)
	logger.SetReportCaller(false)
	logger.SetLevel(log.DebugLevel)
	return logger.WithPrefix("Barnyard")
}
