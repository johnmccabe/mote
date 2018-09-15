[![Travis Badge]][Travis]
[![Go Report Card Badge]][Go Report Card]
[![GoDoc Badge]][GoDoc]

![Mote](go-mote-logo.png)
Buy the Mote controller & accessories here: https://shop.pimoroni.com/products/mote

This repo contains a port of the [Pimoroni Mote library](https://github.com/pimoroni/mote) from Python to Go. It has been verified to work on OSX Sierra and Raspberry Pi 3 (Raspbian). Do let me know if you run on a different platform.

If you have a PHAT Mote go here: https://github.com/johnmccabe/motephat

# Prerequisites

You should have Go version 1.8+ installed and your `GOPATH` configured.

You will need to have `gcc` available on your `PATH` in order to build with the library due to the `go-serial-native` libraries use of `cgo`.

### Raspberry Pi

Install `gcc` using the `apt` package manager.
```
sudo apt update
sudo apt install gcc
```

### OSX

Install `gcc` via the XCode Command Line tools.
```
xcode-select --install
```

# Installation

Install the mote library with `go get`, like so:

```bash
go get -u github.com/johnmccabe/mote
```
You can of course use your own choice of depedency management tool, Glide, Godep etc.

# Examples

You can run the supplied example programs (ported from their Python equivalents) as follows (installing `glide` first which is used to pull down the examples dependencies).
```
go get github.com/Masterminds/glide
cd $GOPATH/src/github.com/johnmccabe/mote
glide install
```
Then running each example as follows.
```
go run examples/rgb/rgb.go 255 0 0
go run examples/rainbow/rainbow.go
go run examples/cheerlights/cheerlights.go
```



*The Golang Gopher was created by [Renée French](http://reneefrench.blogspot.co.uk/) and is [Creative Commons Attributions 3.0](https://creativecommons.org/licenses/by/3.0/) licensed.*

[Travis]: https://travis-ci.org/johnmccabe/mote
[Travis Badge]: https://travis-ci.org/johnmccabe/mote.svg?branch=master
[Go Report Card]: https://goreportcard.com/report/github.com/johnmccabe/mote
[Go Report Card Badge]: https://goreportcard.com/badge/github.com/johnmccabe/mote
[GoDoc]: https://godoc.org/github.com/johnmccabe/mote
[GoDoc Badge]: https://godoc.org/github.com/johnmccabe/mote?status.svg
