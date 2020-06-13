#!/bin/bash
cd src
GOOS=windows GOARCH=386 go build -o ../itb_uploader.exe -v .
