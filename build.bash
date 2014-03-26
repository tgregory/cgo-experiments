#!/bin/bash
gcc -c -o sample.o sample.c
ar rcs libsample.a sample.o
rm sample.o
go build without_pinning.go
