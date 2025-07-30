PROJECT_NAME=ponghub

BINARY=bin/$(PROJECT_NAME)
SRC=cmd/$(PROJECT_NAME)/*.go

.PHONY: all build run clean

all: build

build:
	go build -o $(BINARY) $(SRC)

run: build
	$(BINARY)

clean:
	del $(BINARY)
