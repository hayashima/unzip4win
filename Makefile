.PHONY: build rebuild deps clean assets build-win build-mac

build: deps assets build-win build-mac;
rebuild: clean build;

deps:
	go get -u github.com/jessevdk/go-assets-builder
	go mod download

clean:
	rm -rf ./out/*

assets:
	touch config.toml
	go generate -v

build-win:
	GOOS=windows GOARCH=amd64 go build -ldflags="-w -s -X main.version=`date "+%Y%m%d%H%M%S"`" -v -o out/unzip4win/unzip4win.exe main.go

build-mac:
	GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s -X main.version=`date "+%Y%m%d%H%M%S"`" -v -o out/unzip4win/unzip4win main.go
