#!/bin/bash -e

GOOS=linux GOARCH=mips go build -ldflags="-w" .

