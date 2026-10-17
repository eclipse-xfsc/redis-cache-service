# syntax=docker/dockerfile:1.7

# --------------------------------------------------
# Build stage
# --------------------------------------------------
FROM golang:1.25-bookworm AS builder

WORKDIR /src

# Dependencies
COPY go.mod go.sum ./
COPY vendor/ ./vendor/

# Application sources
COPY . .

# Build static binary
RUN CGO_ENABLED=0 go build \
    -trimpath \
    -mod=vendor \
    -buildvcs=false \
    -ldflags="-s -w" \
    -o /cache \
    ./cmd/cache

# --------------------------------------------------
# Runtime stage
# --------------------------------------------------
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /

COPY --from=builder --chown=65532:65532 /cache /cache

USER 65532:65532

EXPOSE 8080

ENTRYPOINT ["/cache"]