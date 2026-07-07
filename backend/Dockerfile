FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64
ARG TARGET=api

RUN CGO_ENABLED=${CGO_ENABLED} GOOS=${GOOS} GOARCH=${GOARCH} \
    go build -ldflags="-w -s" -o /app/caiyun-${TARGET} ./cmd/${TARGET}

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata curl wget

RUN addgroup -g 1001 -S appgroup && adduser -u 1001 -S appuser -G appgroup
RUN mkdir -p /app/logs /tmp/caiyun && chown -R appuser:appgroup /app/logs /tmp/caiyun

COPY --from=builder --chown=appuser:appgroup /app/caiyun-${TARGET} /usr/local/bin/caiyun

USER appuser

EXPOSE 8080 8081

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://127.0.0.1:8080/health || wget --no-verbose --tries=1 --spider http://127.0.0.1:8081/health || exit 1

ENTRYPOINT ["caiyun"]
