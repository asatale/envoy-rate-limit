
export ROOT_DIR := $(shell pwd)

.PHONY:all
all: app

.PHONY: bin_dir
bin_dir:
	@mkdir -p ${ROOT_DIR}/bin

.PHONY: proto
proto:
	@make -C app/proto

.PHONY: app
app: bin_dir proto
	@make -C app

clean:
	@echo "Performing cleanup..."
	@rm -rf bin
