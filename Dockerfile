FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o generator ./cmd/generator
RUN CGO_ENABLED=0 go build -o server    ./cmd/server

# Stage 2: runtime
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/generator /app/generator
COPY --from=builder /app/server    /app/server
COPY --from=builder /app/internal/web ./internal/web
# default to server
CMD ["/app/server"]
