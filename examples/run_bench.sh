#!/bin/bash
gomp < bench1.go > bench.go
go build bench.go
./bench
rm -rf bench.go bench
