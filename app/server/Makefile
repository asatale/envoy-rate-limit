ifndef BIN_DIR
BIN_DIR := $(shell pwd)
endif

.PHONY: all
all: server.go
	@echo "Building server..."
	@CGO_ENABLED=0  go build -o $(BIN_DIR)/server server.go
