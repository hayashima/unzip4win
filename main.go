package main

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"github.com/ryosms/unzip4win/lib"
	"go.uber.org/zap"
	"log"
	"os"
)

//go:generate go-assets-builder -p unzip4win -o lib/assets.go config.toml

func main() {
	args := unzip4win.ParseArgs()

	logger, err := unzip4win.AppLogger(args.IsDebug)
	if err != nil {
		log.Printf("Can't initialize logger:\n %#v", errors.WithStack(err))
		exitWith(1)
	}

	logger.Debug("argument is",
		zap.String("config", args.ConfigFile),
		zap.String("zip", args.ZipFile))

	config, err := unzip4win.LoadConfig(args.ConfigFile)
	if err != nil {
		logger.Error("Can't parse config file.", zap.Error(errors.WithStack(err)))
		exitWith(1)
	}
	logger.Debug(fmt.Sprintf("config: %v", *config))
	unzip4win.LibLogger = logger

	if !unzip4win.IsLookLikeZipFile(args.ZipFile) {
		logger.Error("It seems to be not a zip file.", zap.String("zipPath", args.ZipFile))
		exitWith(1)
	}
	err = unzip4win.Unzip(args.ZipFile, config)
	if err != nil {
		logger.Error("Failed unzip",
			zap.String("zipPath", args.ZipFile),
			zap.Error(errors.WithStack(err)))
		exitWith(1)
	}

	if args.IsDebug {
		exitWith(0)
	}
}

func exitWith(exitCode int) {
	fmt.Print("Input any for exit...")
	bufio.NewScanner(os.Stdin).Scan()
	os.Exit(exitCode)
}
