ARG GO_VERSION=1.25
ARG ALPINE_VERSION=3.23

FROM golang:${GO_VERSION}-alpine AS builder

ARG VERSION=v1.0.0
ARG COMMIT_SHA=abc123xyz789
ARG BUILD_DATE=2026-01-06

RUN apk add --no-cache tzdata

WORKDIR /app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" 
    -trimpath \
    -o healthcheck ./cmd/healthcheck/healthcheck.go

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=${VERSION} -X main.CommitSHA=${COMMIT_SHA} -X main.BuildDate=${BUILD_DATE}" \
    -trimpath \
    -o main ./cmd/instay/main.go

FROM gcr.io/distroless/static-debian12:nonroot AS production

WORKDIR /app 

COPY --from=builder --chown=nonroot:nonroot /app/healthcheck .

COPY --from=builder --chown=nonroot:nonroot /app/main .

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["./healthcheck"]

USER nonroot

ENV GIN_MODE=release

EXPOSE 8080

CMD ["./main"]