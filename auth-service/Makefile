APP_NAME = auth-service
CMD_DIR = ./cmd/auth-service

.DEFAULT_GOAL := build

.PHONY: fmt vet build run clean
.PHONY: migrate-up migrate-down

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

migrate-up:
	@echo "Migrate up $(APP_NAME)..."
	go run ./cmd/migrate up

migrate-down:
	@echo "Migrate down $(APP_NAME)..."
	go run ./cmd/migrate down