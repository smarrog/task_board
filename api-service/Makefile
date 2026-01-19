APP_NAME = api-service
CMD_DIR = ./cmd/api-service

.DEFAULT_GOAL := build

.PHONY: fmt vet build run clean

fmt:
	@echo "fmt $(APP_NAME)..."
	go fmt ./...

vet: fmt
	@echo "vet $(APP_NAME)..."
	go vet ./...

build: vet
	@echo "Building $(APP_NAME)..."
	go build -o bin/$(APP_NAME) $(CMD_DIR)

run: build
	@echo "Launching $(APP_NAME)..."
	./bin/$(APP_NAME)

clean:
	@echo "Cleaning $(APP_NAME)..."
	rm -rf bin/