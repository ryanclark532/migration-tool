# Makefile for building and running a Go application

# Variables
BINARY_NAME=myapp
TESTING_BINARY_NAME=testapp
SRC=./cmd/main/main.go
TEST_SRC=./cmd/test/main.go

# Targets
.PHONY: all build clean run

all: build

build:
	go build -o ./bin/$(BINARY_NAME) $(SRC)

test: 
	go build -o ./bin/$(TESTING_BINARY_NAME) $(TEST_SRC) && ./bin/$(TESTING_BINARY_NAME)
