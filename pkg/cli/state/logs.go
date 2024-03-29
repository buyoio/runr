package state

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"

	"github.com/buyoio/goodies/logs"
	"github.com/buyoio/goodies/ptr"
	"github.com/buyoio/goodies/streams"
	"github.com/spf13/pflag"
)

const (
	loggerFileName = "runr.log"
)

func (l *Logs) Logger() (*slog.Logger, func() error) {
	verbose := io.Discard
	if l.Verbose {
		verbose = os.Stdout
	}

	w := io.Discard
	var closer func() error
	if l.Path != nil {
		var file io.Writer
		var err error
		file, closer, err = streams.ToFile(path.Join(*l.Path, loggerFileName))
		if err == nil {
			w = io.MultiWriter(file, verbose)
		} else {
			w = verbose
		}
	}

	var unkown string
	var level slog.Level
	if l.Level == nil {
		level = slog.LevelInfo
	} else {
		switch *l.Level {
		case "trace":
			level = logs.LevelTrace
		case "debug":
			level = slog.LevelDebug
		case "info":
			level = slog.LevelInfo
		case "warn":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		case "fatal":
			level = logs.LevelFatal
		default:
			unkown = *l.Level
			level = slog.LevelInfo
		}
	}

	logger := slog.New(logs.NewHandler(w, slog.HandlerOptions{
		Level: level,
	}))
	if unkown != "" {
		logger.Warn("Unknown log level", "level", unkown)
	}
	return logger, closer
}

func (l *Logs) Raw(name string) (*streams.IO, error) {
	streams := &streams.IO{
		Out:    io.Discard,
		ErrOut: io.Discard,
	}

	if l.Path == nil {
		return streams, nil
	}
	_, err := streams.ToFile(path.Join(*l.Path, fmt.Sprintf("%s.log", name)))
	// l.closer = append(l.closer, &streams.Closer)
	return streams, err
}

func (l *Logs) AddFlags(flags *pflag.FlagSet) {
	flags.BoolVarP(&l.Verbose, "verbose", "v", false, "Logs to stderr")
	if l.Level == nil {
		l.Level = ptr.To[string]("info")
	}
	flags.StringVar(l.Level, "log-level", "info", "Log level (debug, info, warn, error)")
}
