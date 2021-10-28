
export ROOT_DIR := $(shell pwd)

all: server client

bin_dir:
	@mkdir -p ${ROOT_DIR}/bin

proto:
	@make -C app/proto


server: bin_dir proto
	@make -C app/server

client: bin_dir proto
	@make -C app/client

clean:
	@echo "Performing cleanup..."
	@rm -rf bin
