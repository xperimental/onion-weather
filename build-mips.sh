#!/bin/bash -e

echo "Build binary..."
GOOS=linux GOARCH=mips go build -v -ldflags="-s -w" .

type upx && {
    echo "Packing using upx..."
    upx -9 onion-weather
} || echo "upx not available"
