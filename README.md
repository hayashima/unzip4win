# unzip4win

[![CircleCI](https://circleci.com/gh/ryosms/unzip4win.svg?style=svg)](https://circleci.com/gh/ryosms/unzip4win)

ref: [kazuhisa/zip4win](https://github.com/kazuhisa/zip4win)

## Usage

```bash
$ unzip4win -h
Usage of Unzip4win:
  unzip4win [OPITIONS] <zip-file-path>
Options
  -config string
        Set path to customized config.toml.
  -debug
        If this flag is settle, output debug log!
```

## for Developer

### Requirement

This app using `go module`.
If you clone into `GOPATH`, set ENVIRONMENT `GO111MODULE=on`.

#### global golang commands

install commands below with `go get`

* go-assets-builder

```bash
$ go get -u \
    github.com/jessevdk/go-assets-builder
    
```

### Build

##### build with make

If you can use `make` command, you can build using `make`.

1. clone this repository
1. edit config.toml as default parameters
1. make

```bash
$ mkdir -p ${GOPATH}/src/github.com/ryosms
$ cd ${GOPATH}/src/github.com/ryosms
$ git clone https://github.com/ryosms/unzip4win.git
$ cd unzip4win
$ cp config.toml.sample config.toml
$ vi config.toml
# edit for your environment
$ make
```

##### build with docker

The repository contains `docker-compose.yml` for building binaries.
If you install `docker` and `docker-compose`, you can build binaries with docker.

1. clone this repository
1. edit config.toml as default parameters
1. docker-compose up

```bash
$ mkdir -p ${GOPATH}/src/github.com/ryosms
$ cd ${GOPATH}/src/github.com/ryosms
$ git clone https://github.com/ryosms/unzip4win.git
$ cd unzip4win
$ cp config.toml.sample config.toml
$ vi config.toml
# edit for your environment
$ docker-compose up
```

##### other

1. clone this repository
1. install dependencies
1. edit config.toml as default parameters
1. create assets
1. go build!

```bash
$ mkdir -p ${GOPATH}/src/github.com/ryosms
$ cd ${GOPATH}/src/github.com/ryosms
$ git clone https://github.com/ryosms/unzip4win.git
$ cd unzip4win
$ go mod download
$ cp config.toml.sample config.toml
$ vi config.toml
# edit for your environment
$ go generate
$ go build -o unzip4win main.go
```
