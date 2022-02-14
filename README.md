# httpservercache - cache of http servers by bind address

[![GoDoc](https://godoc.org/github.com/muir/httpserercache?status.png)](https://pkg.go.dev/github.com/muir/httpserercache)
![unit tests](https://github.com/muir/httpserercache/actions/workflows/go.yml/badge.svg)
[![report card](https://goreportcard.com/badge/github.com/muir/httpserercache)](https://goreportcard.com/report/github.com/muir/httpserercache)
[![codecov](https://codecov.io/gh/muir/httpserercache/branch/main/graph/badge.svg)](https://codecov.io/gh/muir/httpserercache)

Install:

	go get github.com/muir/httpserercache

---

Httpservercache is simply a cache of http servers accessed by
bind address.

For this version, Gorilla Mux is assumed. The API will change with Go 1.18:
the depenency on Gorilla Mux will be replaced with generics.

## Development Status

Development in progress
