
all: hello-world.proto
	@echo "Building protobuf artifacts..."
	@docker run -v ${ROOT_DIR}/app/proto:/defs namely/protoc-all -f hello-world.proto -o . -l go
