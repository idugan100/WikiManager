#! /bin/bash

go fmt *.go
go test -v
go build . && ./WikiManager 