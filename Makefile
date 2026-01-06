ROOT := .
TMP_DIR := tmp
BIN := $(TMP_DIR)/main

.PHONY: build run clean

build:
	@echo "Building..."
	@mkdir -p $(TMP_DIR)
	go build -o $(BIN) ./cmd/instay

run: build
	@echo "Running..."
	@$(BIN)

clean:
	@echo "Cleaning..."
	@rm -rf $(TMP_DIR)