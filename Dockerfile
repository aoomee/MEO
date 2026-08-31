FROM golang:1.27-alpine AS builder

RUN apk add --no-cache ca-certificates git
WORKDIR /src
COPY miaomiaowux/go.mod miaomiaowux/go.sum ./
RUN go mod download
COPY miaomiaowux/ ./
RUN test -f internal/web/dist/index.html && \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/mmwx ./cmd/server

FROM alpine:3.22

RUN apk add --no-cache ca-certificates curl tzdata && \
    addgroup -S mmwx && adduser -S -G mmwx -h /app mmwx && \
    mkdir -p /app/data && chown -R mmwx:mmwx /app
COPY --from=builder /out/mmwx /usr/local/bin/mmwx

USER mmwx
WORKDIR /app
ENV DOCKER=1 \
    PORT=12889 \
    BIND_HOST=0.0.0.0 \
    MMWX_DATA_DIR=/app/data \
    MMWX_UPDATE_REPO=off \
    MMWX_AGENT_GITHUB_REPO=off
VOLUME ["/app/data"]
EXPOSE 12889
HEALTHCHECK --interval=30s --timeout=5s --start-period=30s --retries=3 \
  CMD curl -fsS "http://127.0.0.1:${PORT}/api/setup/status" >/dev/null || exit 1
ENTRYPOINT ["/usr/local/bin/mmwx"]
