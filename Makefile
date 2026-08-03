APP_NAME := shipdeck
DEV_AUTHORIZED_KEYS ?= data/authorized_keys
DEV_PUBLIC_KEY ?= $(HOME)/.ssh/id_ed25519.pub
AUTHORIZED_KEYS_PATH ?= $(DEV_AUTHORIZED_KEYS)

.PHONY: build run test lint dev dev-auth clean

build:
	go build -o bin/$(APP_NAME) ./cmd/shipdeck

run: dev-auth
	AUTHORIZED_KEYS_PATH=$(AUTHORIZED_KEYS_PATH) air

test:
	go test ./...

lint:
	gofmt -w ./cmd ./internal

dev: dev-auth
	AUTHORIZED_KEYS_PATH=$(AUTHORIZED_KEYS_PATH) air

dev-auth:
	@if [ "$(AUTHORIZED_KEYS_PATH)" != "$(DEV_AUTHORIZED_KEYS)" ]; then exit 0; fi
	@mkdir -p $(dir $(DEV_AUTHORIZED_KEYS))
	@if [ ! -f "$(DEV_AUTHORIZED_KEYS)" ]; then \
		if [ -f "$(DEV_PUBLIC_KEY)" ]; then \
			cp "$(DEV_PUBLIC_KEY)" "$(DEV_AUTHORIZED_KEYS)"; \
			chmod 600 "$(DEV_AUTHORIZED_KEYS)"; \
			printf 'created %s from %s\n' "$(DEV_AUTHORIZED_KEYS)" "$(DEV_PUBLIC_KEY)"; \
		else \
			printf 'missing dev public key: %s\n' "$(DEV_PUBLIC_KEY)"; \
			printf 'set DEV_PUBLIC_KEY=/path/to/key.pub or AUTHORIZED_KEYS_PATH=/path/to/authorized_keys\n'; \
			exit 1; \
		fi \
	fi

clean:
	rm -rf bin
