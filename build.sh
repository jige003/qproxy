#!/usr/bin/env bash
# coding=utf-8
# author: b0lu
# mail: b0lu_xyz@163.com
GOOS=windows GOARCH=amd64 go build -x -o qproxy_win64.exe qproxy.go
GOOS=windows GOARCH=386 go build -x -o qproxy_win32.exe qproxy.go
go build -o qproxy qproxy.go 
