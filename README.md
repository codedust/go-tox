gotox
=====

[![GoDoc](https://godoc.org/github.com/codedust/go-tox?status.png)](http://godoc.org/github.com/codedust/go-tox)


gotox is a Go wrapper for the [c-toxcore](https://github.com/TokTok/c-toxcore) library.

Pull requests, bug reportings and feature requests (via github issues) are always welcome. :)

For a list of supported toxcore features see [PROGRESS.md](PROGRESS.md).

## Installation
First, install the [c-toxcore](https://github.com/TokTok/c-toxcore) library.

Next, download `go-tox` using go:
```
go get github.com/codedust/go-tox
```

## License
gotox is licensed under the [GPLv3](COPYING).

## How to use
See [bindings.go](bindings.go) for details about supported API functions and [callbacks.go](callbacks.go) for the supported callbacks.

The best place to get started are the examples in [examples/](examples/).

```
go run examples/example.go
```

Feel free to ask for help in the issue tracker. ;)
