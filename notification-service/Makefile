APP_NAME = notification-service
CMD_DIR = ./cmd/notification-service

PROTO_DIR = proto
PROTO_GEN_DIR = proto
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto) $(wildcard $(PROTO_DIR)/*/*.proto)

.PHONY: proto
proto:
	@echo "Generating protobuf files..."
	protoc \
		-I $(PROTO_DIR) \
		--go_out=$(PROTO_GEN_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_GEN_DIR) \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)

.PHONY: build
build: proto
	@echo "Building $(APP_NAME)..."
	go build -o bin/$(APP_NAME) $(CMD_DIR)

.PHONY: run
run: build
	./bin/$(APP_NAME)

.PHONY: clean
clean:
	rm -rf bin/