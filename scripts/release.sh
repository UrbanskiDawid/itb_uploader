#!/bin/bash
#
# release: build binary executables for all platforms
#
set -e
NAME=main
WD=OUT/`date +%Y-%m-%d_%H%M`
mkdir -p $WD
echo "linux: $WD/$NAME"
go build -o $WD/$NAME
echo "windows: $WD/$NAME.exe"
(GOOS=windows GOARCH=386 go build -o $WD/$NAME.exe)
