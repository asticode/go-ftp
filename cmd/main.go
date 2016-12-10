package main

import (
	stlflag "flag"

	"time"

	"github.com/asticode/go-ftp"
	"github.com/asticode/go-toolkit/flag"
	"github.com/molotovtv/go-logger"
	"github.com/rs/xlog"
	"golang.org/x/net/context"
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

	// Switch on subcommand
	switch s {
	case "download":
		var ctx, _ = context.WithTimeout(context.Background(), 10*time.Minute)
		if err := f.Download(*inputPath, *outputPath, ctx); err != nil {
			l.Fatal(err)
		}
	default:

	}
}
