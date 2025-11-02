# Multi-stage build for polityd

FROM golang:1.25-alpine AS builder
WORKDIR /src
# Install build dependencies (git for go mod) and CA certs for module downloads
#RUN apk add --no-cache git ca-certificates

# Copy go.mod and go.sum first for better layer caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the rest of the source
COPY . .

# Build a static binary
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -trimpath -ldflags="-s -w" -o /out/polityd ./cmd/polityd

# 2) Runtime stage
FROM gcr.io/distroless/static:nonroot
# Copy CA certificates in case outbound TLS is ever needed
#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /out/polityd /polityd

USER nonroot:nonroot
ENTRYPOINT ["/polityd"]
