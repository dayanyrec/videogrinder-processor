FROM golang:1.21-alpine AS tools

# Install only essential packages for development tools
RUN apk add --no-cache \
    git \
    make \
    bash \
    nodejs \
    npm \
    ca-certificates

# Install Go tools
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1 && \
    go install golang.org/x/tools/cmd/goimports@v0.21.0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Install Node.js dependencies for linting
RUN cd web && npm install

FROM golang:1.21-alpine AS development

RUN apk add --no-cache \
    ffmpeg \
    git \
    make \
    bash \
    ca-certificates \
    curl

RUN go install github.com/cosmtrek/air@v1.49.0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p uploads outputs temp

EXPOSE 8080

FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -buildvcs=false \
    -a -installsuffix cgo \
    -ldflags='-w -s -extldflags "-static"' \
    -o web/cmd/main \
    ./web/cmd

RUN CGO_ENABLED=0 GOOS=linux go build \
    -buildvcs=false \
    -a -installsuffix cgo \
    -ldflags='-w -s -extldflags "-static"' \
    -o api/cmd/main \
    ./api/cmd

RUN CGO_ENABLED=0 GOOS=linux go build \
    -buildvcs=false \
    -a -installsuffix cgo \
    -ldflags='-w -s -extldflags "-static"' \
    -o processor/cmd/main \
    ./processor/cmd

FROM alpine:3.18 AS production

RUN apk add --no-cache \
    ffmpeg \
    ca-certificates \
    tzdata \
    wget \
    curl && \
    addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/web/cmd/main ./web/cmd/main
COPY --from=builder /app/api/cmd/main ./api/cmd/main
COPY --from=builder /app/processor/cmd/main ./processor/cmd/main

COPY --from=builder /app/web/static ./web/static

RUN mkdir -p uploads outputs temp && \
    chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./web/cmd/main"]
