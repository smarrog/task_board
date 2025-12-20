APP_NAME = notification-service
CMD_DIR = ./cmd/notification-service

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