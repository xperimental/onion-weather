#!/bin/bash -e

GOOS=linux GOARCH=mips go build -v -ldflags="-w" .

