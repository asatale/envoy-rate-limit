ifndef BIN_DIR
BIN_DIR := $(shell pwd)
endif

.PHONY: all
all: client.go
	@echo "Building client..."
	@CGO_ENABLED=0 go build -o $(BIN_DIR)/client client.go
