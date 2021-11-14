
all: hello-world.proto
	@echo "Building protobuf artifacts..."
	@docker run -v ${PWD}:/defs namely/protoc-all -f hello-world.proto -o go -l go
	@docker run -v ${PWD}:/defs namely/protoc-all -f hello-world.proto -o python -l python
