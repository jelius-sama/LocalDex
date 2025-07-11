VERSION := $(shell jq -r '.version' ./config.json)
APP_NAME := $(shell jq -r '.title' ./config.json)
PORT := $(shell jq -r '.port' ./config.json)
BIN_DIR := ./bin

all:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w -X main.Environment=production -X main.Port=$(PORT) -X main.Version=$(VERSION)" -trimpath -buildvcs=false -o $(BIN_DIR)/$(APP_NAME)-$(VERSION) ./cmd
