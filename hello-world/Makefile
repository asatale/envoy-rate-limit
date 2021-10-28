# Build protobuf artifacts

build-proto:
	@echo "Compiling hello-world.proto"
	@docker run -v `pwd`:/defs namely/protoc-all -f proto/hello-world.proto -o . -l go #or go, csharp, etc
	@mv hello-world/hello-world.pb.go .
	@rm -rf hello-world

clean:
	@echo "Cleaning protobuf artifacts"
	@rm -rf hello-world.pb.go
