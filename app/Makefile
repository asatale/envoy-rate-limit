
ifdef ROOT_DIR
export BIN_DIR := $(ROOT_DIR)/bin
else
$(shell mkdir -p bin)
export BIN_DIR := $(shell pwd)/bin
endif


.PHONY: all
all: server client

.PHONY: go-tidy
go-tidy:
	go mod tidy

.PHONY: server
server: go-tidy
	@make -C server

.PHONY: client
client: go-tidy
	@make -C client
