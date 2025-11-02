FROM golang:1.25-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -trimpath -ldflags="-s -w" -o /bin/polityd ./cmd/polityd

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /bin/polityd /bin/polityd

USER nonroot:nonroot
ENTRYPOINT ["/bin/polityd"]
