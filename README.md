[![Go Report Card Badge]][Go Report Card]
[![GoDoc Badge]][GoDoc]

![Mote](go-mote-logo.png)
Buy the Mote controller & accessories here: https://shop.pimoroni.com/products/mote

This repo contains a port of the [Pimoroni Mote library](https://github.com/pimoroni/mote) from Python to Go.

If you have a PHAT Mote go here: https://github.com/johnmccabe/motephat

## Prerequisites

### OSX

The `go.bug.st/serial/enumerator` library used to detect the Mote port requires the use of cgo on MacOSX in order to access the IOKit Framework so you will need to have XCode installed.

```
xcode-select --install
```

## Using

Import the library via:

```go
import "github.com/johnmccabe/mote"
```

You can refer to the examples in the `example\` directory and the [GoDocs][GoDoc] for information on using the library.

## Examples

You can run the supplied example programs (ported from their Python equivalents) as follows.

```shell
go mod tidy
```

Then running each example as follows.

```shell
go run examples/rgb/rgb.go 255 0 0
go run examples/rainbow/rainbow.go
go run examples/cheerlights/cheerlights.go
```



*The Golang Gopher was created by [Renée French](http://reneefrench.blogspot.co.uk/) and is [Creative Commons Attributions 3.0](https://creativecommons.org/licenses/by/3.0/) licensed.*

[Go Report Card]: https://goreportcard.com/report/github.com/johnmccabe/mote
[Go Report Card Badge]: https://goreportcard.com/badge/github.com/johnmccabe/mote
[GoDoc]: https://pkg.go.dev/github.com/johnmccabe/mote
[GoDoc Badge]: https://pkg.go.dev/badge/github.com/johnmccabe/mote.svg
