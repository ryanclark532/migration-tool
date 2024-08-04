# Makefile for building and running a Go application

# Variables
BINARY_NAME=myapp
TESTING_BINARY_NAME=testapp
SRC=./cmd/main/main.go
TEST_SRC=./cmd/main/main_test.go
TEST_DB_SRC=./cmd/main/server.db

# Targets
.PHONY: all build clean run

all: build

build:
	go build -o ./bin/$(BINARY_NAME).exe $(SRC)

test: 
	go test $(TEST_SRC) && rm $(TEST_DB_SRC)
