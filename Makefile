export ROOT_DIR := $(shell pwd)

.PHONY:all
all: app


.PHONY:build
build:
	@docker-compose build


.PHONY:up
up:
	@docker-compose up -d

.PHONY:down
down:
	@docker-compose down -v


.PHONY:logs
logs:
	@docker-compose logs -f


.PHONY: proto
proto:
	@make -C app/proto


.PHONY: app
app: bin_dir proto
	@make -C app


.PHONY: bin_dir
bin_dir:
	@mkdir -p ${ROOT_DIR}/bin


clean:
	@echo "Performing cleanup..."
	@rm -rf bin
