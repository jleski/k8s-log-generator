BIN_NAME := loggen
BIN_NAME_LN := loggen_linux
SRC := $(wildcard *.go)

.PHONY: build_and_run

build_and_run: $(BIN_NAME) run

$(BIN_NAME): $(SRC)
	go build -o $(BIN_NAME) .

run: $(BIN_NAME)
	./$(BIN_NAME) -interval 2

$(BIN_NAME_LN): $(SRC)
	GOOS=linux GOARCH=amd64 go build -o $(BIN_NAME_LN) .
