.PHONY: clean

BIN_DIR := bin
BIN_FILE := $(BIN_DIR)/dg
SRC_MAIN := cmd/dg.go
SRC_FILES := $(SRC_MAIN) $(wildcard config/*.go) $(wildcard utils/*.go)

$(BIN_FILE): $(BIN_DIR) $(SRC_FILES)
	go build -o $@ $(SRC_MAIN)

$(BIN_DIR):
	mkdir -p $@

clean:
	rm -Rf $(BIN_DIR)