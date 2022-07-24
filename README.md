# koron/thinout

[![PkgGoDev](https://pkg.go.dev/badge/github.com/koron/thinout)](https://pkg.go.dev/github.com/koron/thinout)
[![Actions/Go](https://github.com/koron/thinout/workflows/Go/badge.svg)](https://github.com/koron/thinout/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/koron/thinout)](https://goreportcard.com/report/github.com/koron/thinout)

Tool to randomly thin out lines in a text file.

You can specify...

1. Rate to output (`-r`)
2. Seed for random numbers (`-s`)
2. Line numbers to forced output (`-f`)

## How to install

```console
$ go install github.com/koron/thinout@latest
```

## Usage

```console
$ ./thinout -h
Usage of thinout:
  -f value
        specify fixed line numbers, comma separated
  -r float
        rate to output [0.0,1.0] (default 0.1)
  -s value
        seed for random numbers: "time" or int64 (default: 0)
```
