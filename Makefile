export ROOT_DIR := $(shell pwd)

FRAMEWORK ?= python


.PHONY:all
all: app


.PHONY:build
build:
ifeq ($(FRAMEWORK), go)
	@echo "DOCKERFILE=Dockerfile.golang" > .env
	@docker-compose build
else
ifeq ($(FRAMEWORK), python)
	@echo "DOCKERFILE=Dockerfile.python" > .env
	@docker-compose build
else
	$(error "Unsupported framework $(FRAMEWORK)")
endif
endif

.PHONY:up
up:
	@docker-compose up -d

.PHONY:down
down:
	@docker-compose down -v


.PHONY:logs
logs:
	@docker-compose logs -f

