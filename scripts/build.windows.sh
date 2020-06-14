#!/bin/bash
GOOS=windows GOARCH=386 go build -o itb_uploader.exe -v .
