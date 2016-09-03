package main

import (
	"flag"

	"github.com/asticode/go-ftp"
	"github.com/molotovtv/go-logger"
	"github.com/molotovtv/go-toolbox"
	"github.com/rs/xlog"
)

// Flags
var (
	outputPath = flag.String("o", "", "the output path")
	inputPath  = flag.String("i", "", "the input path")
)

func main() {
	// Get subcommand
	s := toolbox.Subcommand()
	flag.Parse()

	// Init logger
	l := xlog.New(logger.NewConfig(logger.FlagConfig()))

	// Init ftp
	f := ftp.NewFromConfig(ftp.FlagConfig())
	f.Logger = l

	// Log
	l.Debugf("Subcommand is %s", s)

	// Switch on subcommand
	switch s {
	case "copy":
		if err := f.Copy(*inputPath, *outputPath); err != nil {
			l.Fatal(err)
		}
	default:

	}
}
