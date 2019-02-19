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
)

func main() {
	args := unzip4win.ParseArgs()

	logger, err := unzip4win.AppLogger(args.IsDebug)
	if err != nil {
		log.Fatalf("Can't initialize logger: %v", err)
	}

	logger.Debug("argument is",
		zap.String("config", args.ConfigFile),
		zap.String("zip", args.ZipFile))

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

	fmt.Print("Input any for exit...")
	bufio.NewScanner(os.Stdin).Scan()
}
