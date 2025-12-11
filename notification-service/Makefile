# Где лежат .proto
PROTO_DIR = proto
# Куда генерируем код
PROTO_GEN_DIR = proto

# Все .proto файлы
PROTO_FILES := $(wildcard $(PROTO_DIR)/**/*.proto)

.PHONY: proto
proto:
	protoc \
		-I $(PROTO_DIR) \
		--go_out=$(PROTO_GEN_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_GEN_DIR) \
        --go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)