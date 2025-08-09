CONFIG_FILE := ./config/config.json
CONFIG_DIR := $(dir $(CONFIG_FILE))

VERSION := $(shell jq -r '.references[].path' $(CONFIG_FILE) | xargs -I{} jq -r 'select(has("version")) | .version' $(CONFIG_DIR){} )
APP_NAME := $(shell jq -r '.references[].path' $(CONFIG_FILE) | xargs -I{} jq -r 'select(has("title")) | .title' $(CONFIG_DIR){} )
DESCRIPTION := $(shell jq -r '.references[].path' $(CONFIG_FILE) | xargs -I{} jq -r 'select(has("description")) | .description' $(CONFIG_DIR){} )
DEV_PORT := $(shell jq -r '.references[].path' $(CONFIG_FILE) | xargs -I{} jq -r 'select(has("dev_port")) | .dev_port' $(CONFIG_DIR){} )
APP_ROOT := $(shell jq -r '.references[].path' $(CONFIG_FILE) | xargs -I{} jq -r 'select(has("app_root")) | .app_root' $(CONFIG_DIR){} )
HOST := $(shell jq -r '.references[].path' $(CONFIG_FILE) | xargs -I{} jq -r 'select(has("host")) | .host' $(CONFIG_DIR){} )
IS_REVERSE_PROXIED := $(shell \
  jq -r '.references[].path' $(CONFIG_FILE) | \
  xargs -I{} jq -er 'if has("is_behind_reverse_proxy") then .is_behind_reverse_proxy.statement_valid else empty end' $(CONFIG_DIR)/{} \
)
PORT := $(shell \
  jq -r '.references[].path' $(CONFIG_FILE) | \
  xargs -I{} jq -er 'if has("is_behind_reverse_proxy") then .is_behind_reverse_proxy.port_to_use else empty end' $(CONFIG_DIR)/{} \
)

BIN_DIR := ./bin
ENV := production

.PHONY: build dev deploy_prod preview

build:
	(cd client && bun run build)
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=linux go build -ldflags "\
		    -s -w \
		    -X main.Environment=$(ENV) \
		    -X main.AppRoot=$(APP_ROOT) \
		    -X main.ETCDir=/etc/$(APP_NAME) \
		    -X main.Host=$(HOST) \
		    -X main.DevPort=$(DEV_PORT) \
		    -X main.Title=$(APP_NAME) \
		    -X main.Version=$(VERSION) \
		    -X main.ReverseProxy="$(IS_REVERSE_PROXIED)" \
		    -X main.ProxyPort=$(PORT)" \
		    -trimpath -buildvcs=false -o $(BIN_DIR)/$(APP_NAME)-$(VERSION) ./cmd

dev:
	@echo "Starting development servers..."
	@( \
		pids=""; \
		cleanup() { \
			echo "Gracefully shutting down all servers..."; \
			for pid in $$pids; do \
				if kill -0 $$pid 2>/dev/null; then \
					kill -TERM $$pid 2>/dev/null || true; \
				fi; \
			done; \
			sleep 0.5; \
			for pid in $$pids; do \
				if kill -0 $$pid 2>/dev/null; then \
					kill -KILL $$pid 2>/dev/null || true; \
				fi; \
			done; \
			wait 2>/dev/null || true; \
			echo "All servers stopped."; \
			exit 0; \
		}; \
		trap cleanup INT TERM EXIT; \
		echo "Starting backend server..."; \
		air & pids="$$pids $$!"; \
		echo "Starting frontend dev server..."; \
		(cd client && bun run dev) & pids="$$pids $$!"; \
		echo "All development servers started! Press Ctrl+C to stop all servers."; \
		wait; \
	)

# Files and folders to include in the archive
deploy_prod:
	sudo systemctl stop LocalDex
	@make build -B
	sudo cp $(BIN_DIR)/$(APP_NAME)-$(VERSION) /usr/local/bin/LocalDex
	sudo mkdir -p /etc/$(APP_NAME)
	sudo cp -r ./config /etc/$(APP_NAME)
	sudo systemctl restart LocalDex.service

.ONESHELL:

preview:
	# sudo mkdir -p /etc/$(APP_NAME)
	# sudo cp -r ./config /etc/$(APP_NAME)
	# sudo env $(shell grep -v '^#' /etc/$(APP_NAME)/$(APP_NAME).env | xargs) $(BIN_DIR)/$(APP_NAME)-$(VERSION)
	sudo env $(shell grep -v '^#' /etc/$(APP_NAME)/$(APP_NAME).env | xargs -d '\n') $(BIN_DIR)/$(APP_NAME)-$(VERSION)
