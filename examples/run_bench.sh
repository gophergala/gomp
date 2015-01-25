#!/bin/bash
gomp < bench.go > bench1.go
go build bench1.go
./bench1
rm -rf bench1.go bench1
