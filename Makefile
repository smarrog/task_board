.DEFAULT_GOAL := help

.PHONY: help
.PHONY: proto-shared
.PHONY: build-api build-auth build-core build-notification build-all

help:
	@echo ""
	@echo "Root targets:"
	@echo "  build-all     - build all services"
	@echo "  proto-shared  - generate proto in module \"shared\""

build-api:
	$(MAKE) -C api-service build

build-auth:
	$(MAKE) -C auth-service build

build-core:
	$(MAKE) -C core-service build

build-notification:
	$(MAKE) -C notification-service build

build-all: build-api build-auth build-core build-notification

proto-shared:
	$(MAKE) -C shared proto