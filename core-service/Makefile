APP_NAME = core-service
CMD_DIR = ./cmd/core-service

.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	go build -o bin/$(APP_NAME) $(CMD_DIR)

.PHONY: run
run: build
	./bin/$(APP_NAME)

.PHONY: clean
clean:
	rm -rf bin/

.PHONY: migrate-up
migrate-up:
	go run ./cmd/migrate up

.PHONY: migrate-down
migrate-down:
	go run ./cmd/migrate down