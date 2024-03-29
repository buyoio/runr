package state

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/buyoio/goodies/streams"

	"gopkg.in/yaml.v3"
)

func (l *Runr) State(file ...bool) *State {
	if len(file) > 0 && file[0] {
		return l.fileState
	}
	return l.aggrState
}

func (l *Runr) Path() string {
	return getDefaultConfigPath()
}

func (l *Runr) Marshal() error {
	yamlFile, err := yaml.Marshal(l.fileState)
	if err != nil {
		return err
	}
	return os.WriteFile(l.Path(), yamlFile, 0644)
}

func (runr *Runr) Logger() *slog.Logger {
	if runr.logger != nil {
		return runr.logger
	}
	// ensure logs path is correct
	// and reset it after logger creation
	if runr.fileState.Logs.Path == nil {
		runr.logger, _ = runr.fileState.Logs.Logger()
		return runr.logger
	}
	// todo don't like this at all
	p := runr.fileState.Logs.Path
	*runr.fileState.Logs.Path = resolvePath(*runr.fileState.Logs.Path)
	runr.logger, runr.closer = runr.fileState.Logs.Logger()
	runr.fileState.Logs.Path = p
	return runr.logger
}

func (runr *Runr) IO() *streams.IO {
	return runr.io
}

// https://github.com/tilt-dev/tilt/blob/master/internal/cli/cli.go#L115C6-L155
func (runr *Runr) GetContext() context.Context {
	// logger := runr.Logger()

	var ctx context.Context
	if runr.cmd != nil {
		ctx = runr.cmd.Context()
	} else {
		ctx = context.Background()
		// logger.Debug("No context found; using background context.")
	}

	// SIGNAL TRAPPING
	ctx, cancel := context.WithCancel(ctx)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		<-sigs
		cancel()

		// If we get another signal, OR it takes too long
		// to exit after canceling context, just exit
		select {
		case <-sigs:
			runr.Logger().Debug("force quitting...")
			os.Exit(1)
		case <-time.After(2 * time.Second):
			runr.Logger().Debug("Context canceled but app still running; forcibly exiting.")
			os.Exit(1)
		}
	}()

	return ctx
}
