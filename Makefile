PROTOC_GEN_GO := $(shell go env GOPATH)/bin/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(shell go env GOPATH)/bin/protoc-gen-go-grpc
PROTOC_GEN_GRPC_GATEWAY := $(shell go env GOPATH)/bin/protoc-gen-grpc-gateway
PROTOC_GEN_OPENAPIV2 := $(shell go env GOPATH)/bin/protoc-gen-openapiv2
GOOGLEAPIS := $(shell go list -f '{{ .Dir }}' -m github.com/googleapis/googleapis)
GRPC_GATEWAY := $(shell go list -f '{{ .Dir }}' -m github.com/grpc-ecosystem/grpc-gateway/v2)

PROTO_SRC := api/hello/v1/hello.proto
OUT_DIR := api/hello/v1
# Target to generate the Go code and Swagger file
gen:
	protoc -I . \
		-I $(GOOGLEAPIS) \
		-I $(GRPC_GATEWAY) \
		--go_out=$(OUT_DIR) --go-grpc_out=$(OUT_DIR) \
		--grpc-gateway_out=$(OUT_DIR) \
		--openapiv2_out=$(OUT_DIR) \
		$(PROTO_SRC)

	@echo "Generated Go files and Swagger JSON!"
