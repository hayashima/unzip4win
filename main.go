package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/yeka/zip"
)

func main() {
	configFile := flag.String("config", "config.toml", "")
	flag.Parse()

	fmt.Printf("%v\n", flag.Args())
	fmt.Println(*configFile)

	zipPath := flag.Arg(0)
	fmt.Println(zipPath)

	reader, err := zip.OpenReader(zipPath)
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
