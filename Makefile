# Makefile for ponghub Go project

BINARY=bin/ponghub.exe
SRC=cmd/main.go

.PHONY: all build run clean

all: build

build:
	go build -o $(BINARY) $(SRC)

run: build
	$(BINARY)

clean:
	del $(BINARY)
