package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/molotovtv/go-logger"
	"github.com/molotovtv/go-astitools/context"
	"github.com/molotovtv/go-astitools/flag"
	"github.com/molotovtv/go-ftp"
)

// Flags
var (
	outputPath = flag.String("o", "", "the output path")
	inputPath  = flag.String("i", "", "the input path")
)

func main() {
	log.Setup("go-ftp")
	// Get subcommand
	s := astiflag.Subcommand()
	flag.Parse()
	// Init logger
	log.FlagInit()

	// Init ftp
	f := ftp.New(ftp.FlagConfig())

	// Log
	log.Debugf("Subcommand is %s", s)

	// Init canceller
	var c = asticontext.NewCanceller()

	// Handle signals
	handleSignals(c)

	// Switch on subcommand
	switch s {
	case "download":
		var ctx, _ = c.NewContext()
		if err := f.Download(ctx, *inputPath, *outputPath); err != nil {
			log.Fatal(err)
		}
	case "upload":
		var ctx, _ = c.NewContext()
		if err := f.Upload(ctx, *inputPath, *outputPath); err != nil {
			log.Fatal(err)
		}
	}
}

// handleSignals handles signals
func handleSignals(c *asticontext.Canceller) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		for s := range ch {
			log.Debugf("Received signal %s", s)
			c.Cancel()
		}
	}()
}
