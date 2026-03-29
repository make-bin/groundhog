# @AI_GENERATED
.PHONY: build test lint run migrate-up migrate-down clean generate frontend-dev frontend-build build-full docker-build docker-up docker-down perf-test perf-test-health perf-test-session perf-test-ws

APP_NAME := groundhog
BIN_DIR := bin
CMD_DIR := ./cmd/server/
CONFIG_PATH := configs/config.yaml
MIGRATIONS_DIR := migrations

build:
	go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)

frontend-dev:
	cd web/app && npm install && npm run dev

frontend-build:
	cd web/app && npm install && npm run build

build-full: frontend-build build

test:
	go test ./...

lint:
	golangci-lint run ./...

run: build
	$(BIN_DIR)/$(APP_NAME) gateway run --config $(CONFIG_PATH)

migrate-up:
	$(BIN_DIR)/$(APP_NAME) migrate up --config $(CONFIG_PATH)

migrate-down:
	$(BIN_DIR)/$(APP_NAME) migrate down --config $(CONFIG_PATH)

clean:
	rm -rf $(BIN_DIR)
	go clean

generate:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/channel/v1/channel.proto

docker-build:
	docker build -t groundhog:latest .

docker-up:
	docker compose up -d

docker-down:
	docker compose down

perf-test:
	@echo "Running k6 performance tests..."
	k6 run perf/k6/health.js
	k6 run perf/k6/session.js
	k6 run perf/k6/websocket.js

perf-test-health:
	k6 run perf/k6/health.js

perf-test-session:
	k6 run perf/k6/session.js

perf-test-ws:
	k6 run perf/k6/websocket.js
# @AI_GENERATED: end
