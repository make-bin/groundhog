# Stage 1: Frontend build
FROM node:20-alpine AS frontend
WORKDIR /app/web

COPY web/control-ui/package*.json web/control-ui/
RUN cd control-ui && npm ci --prefer-offline

COPY web/control-ui/ control-ui/
RUN cd control-ui && npm run build

COPY web/chat-ui/package*.json web/chat-ui/
RUN cd chat-ui && npm ci --prefer-offline

COPY web/chat-ui/ chat-ui/
RUN cd chat-ui && npm run build

# Stage 2: Go build
FROM golang:1.21-alpine AS builder
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go module files first for layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy frontend build artifacts
COPY --from=frontend /app/web/control-ui/dist web/control-ui/dist
COPY --from=frontend /app/web/chat-ui/dist web/chat-ui/dist

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/groundhog ./cmd/server/

# Stage 3: Runtime image
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/bin/groundhog /usr/local/bin/groundhog
COPY --from=builder /app/configs /etc/groundhog/configs
COPY --from=builder /app/migrations /etc/groundhog/migrations

# Create non-root user
RUN addgroup -S groundhog && adduser -S groundhog -G groundhog
USER groundhog

EXPOSE 8080

ENTRYPOINT ["groundhog", "gateway", "run", "--config", "/etc/groundhog/configs/config.yaml"]
