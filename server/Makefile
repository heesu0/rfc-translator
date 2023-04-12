BINARY_NAME=rfc-translator
BUILD_FOLDER=build

.PHONY: build
build:
	@go build -o ./${BUILD_FOLDER}/${BINARY_NAME} ./cmd/main.go

.PHONY: clean
clean:
	@rm -rf ./${BUILD_FOLDER}/*

.PHONY: lint
lint:
	@golangci-lint run