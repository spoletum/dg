.PHONY: clean test

BIN_DIR := bin
BIN_FILE := $(BIN_DIR)/dg
SRC_MAIN := cmd/dg.go
SRC_FILES := $(SRC_MAIN) $(wildcard pkg/**/*.go) $(wildcard internal/**/*.go)

all: test $(BIN_FILE)

$(BIN_FILE): $(BIN_DIR) $(SRC_FILES)
	go build -o $@ $(SRC_MAIN)

$(BIN_DIR):
	mkdir -p $@

test:
	go test -cover ./...

clean:
	rm -Rf $(BIN_DIR)