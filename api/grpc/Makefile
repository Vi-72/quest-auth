MODULE_NAME = github.com/Vi-72/quest-auth/api/grpc/sdk/go

.PHONY: install-tools
install-tools:
	$(info installing deps...)
	GOPROXY=https://proxy.golang.org,direct go install github.com/bufbuild/buf/cmd/buf@v1.28.1
	GOPROXY=https://proxy.golang.org,direct go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2
	GOPROXY=https://proxy.golang.org,direct go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.4.0

.PHONY: clean
clean:
	$(info clearing sdk directory...)
	@rm -rf sdk

.PHONY: generate
generate:
	$(info generating go files from proto files...)
	PATH="$(shell go env GOPATH)/bin:$$PATH" buf generate
	@if [ ! -f ./sdk/go/go.mod ]; then \
	  (cd ./sdk/go && go mod init $(MODULE_NAME)); \
	fi
	@cd ./sdk/go && GOPROXY=https://proxy.golang.org,direct go mod tidy


