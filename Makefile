.PHONY: clean e2e

BIN_DIR := bin
BIN_FILE := $(BIN_DIR)/dg
SRC_MAIN := cmd/dg.go
SRC_FILES := $(SRC_MAIN) $(wildcard config/**/*.go) $(wildcard utils/**/*.go)

all: $(COVER_FILE) $(BIN_FILE)

$(BIN_FILE): $(BIN_DIR) $(SRC_FILES)
	go test -cover ./...
	go build -o $@ $(SRC_MAIN)

$(BIN_DIR):
	mkdir -p $@

clean:
	rm -Rf $(BIN_DIR)

e2e:
	go test -v github.com/spoletum/dg/e2e/...  