package main

import (
	stlflag "flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/asticode/go-ftp"
	"github.com/asticode/go-toolkit/context"
	"github.com/asticode/go-toolkit/flag"
	"github.com/molotovtv/go-logger"
	"github.com/rs/xlog"
)

// Flags
var (
	outputPath = stlflag.String("o", "", "the output path")
	inputPath  = stlflag.String("i", "", "the input path")
)

func main() {
	// Get subcommand
	s := flag.Subcommand()
	stlflag.Parse()

	// Init logger
	l := xlog.New(logger.NewConfig(logger.FlagConfig()))

	// Init ftp
	f := ftp.New(ftp.FlagConfig())
	f.Logger = l

	// Log
	l.Debugf("Subcommand is %s", s)

	// Init canceller
	var c = context.NewCanceller()

	// Handle signals
	handleSignals(l, c)

	// Switch on subcommand
	switch s {
	case "download":
		var ctx, _ = c.NewContext()
		if err := f.Download(ctx, *inputPath, *outputPath); err != nil {
			l.Fatal(err)
		}
	case "upload":
		var ctx, _ = c.NewContext()
		if err := f.Upload(ctx, *inputPath, *outputPath); err != nil {
			l.Fatal(err)
		}
	}
}

// handleSignals handles signals
func handleSignals(l xlog.Logger, c *context.Canceller) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		for s := range ch {
			l.Debugf("Received signal %s", s)
			c.Cancel()
		}
	}()
}
