# Base builder cho Go
FROM golang:1.24-alpine AS base-builder

WORKDIR /app

# Copy và download module
COPY go.mod go.sum ./
RUN go mod download

# Install bash and git (Air requires bash)
RUN apk add --no-cache bash git

# Cài đặt tool phục vụ dev mode
# Air: hot reload
# Swag: generate swagger docs
RUN go install github.com/air-verse/air@latest \
    && go install github.com/swaggo/swag/cmd/swag@latest

ENV PATH=$PATH:/go/bin

# Copy toàn bộ source code
COPY . .

WORKDIR /app

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8081/health || exit 1

EXPOSE 8081

# Dev mode sẽ dùng `air` để auto reload khi code thay đổi
ENTRYPOINT  ["air", "-c", ".air.toml"]


