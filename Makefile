generate:
	@echo "Building gRPC and Protobuf in Go format"
	protoc --experimental_allow_proto3_optional --proto_path=. --go_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --go-grpc_out=. ./api/*.proto

clean:
	@echo "Removing compiled protobuf messages"
	rm -R *.pb.go
