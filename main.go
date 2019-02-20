package main

import (
	"bufio"
	"fmt"
	"github.com/ryosms/unzip4win/lib"
	"github.com/yeka/zip"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	args := unzip4win.ParseArgs()

	logger, err := unzip4win.AppLogger(args.IsDebug)
	if err != nil {
		log.Printf("Can't initialize logger:\n %v", err)
		exitWith(1)
	}

	logger.Debug("argument is",
		zap.String("config", args.ConfigFile),
		zap.String("zip", args.ZipFile))

	config, err := unzip4win.LoadConfig(args.ConfigFile)
	if err != nil {
		logger.Error("Can't parse config file.", zap.Error(err))
		exitWith(1)
	}
	logger.Debug(fmt.Sprintf("config: %v", *config))

	logger.Debug("extension", zap.String("ext", filepath.Ext(args.ZipFile)))
	if !unzip4win.IsLookLikeZipFile(args.ZipFile) {
		logger.Error("It seems to be not a zip file.", zap.String("zipPath", args.ZipFile))
		exitWith(1)
	}
	err = unzip4win.Unzip(args.ZipFile, config, logger)
	if err != nil {
		logger.Error("Failed unzip",
			zap.String("zipPath", args.ZipFile),
			zap.Error(err))
		exitWith(1)
	}

	reader, err := zip.OpenReader(args.ZipFile)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	for _, f := range reader.File {
		if f.IsEncrypted() {
			f.SetPassword("password")
		}
		r, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}

		buf, err := ioutil.ReadAll(r)
		if err != nil {
			log.Fatal(err)
		}
		defer r.Close()

		fmt.Printf("Size of %v: %v byte(s)\n", f.Name, len(buf))
	}

	exitWith(0)
}

func exitWith(exitCode int) {
	fmt.Print("Input any for exit...")
	bufio.NewScanner(os.Stdin).Scan()
	os.Exit(exitCode)
}
