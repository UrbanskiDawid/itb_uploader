#!/bin/bash
go get golang.org/x/crypto/ssh
GOOS=windows GOARCH=386 go build -o itb_uploader.exe -v .
