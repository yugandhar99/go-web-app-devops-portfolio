# syntax=docker/dockerfile:1.7

ARG GO_VERSION=1.23

# Base stage used for dependency download.
FROM golang:${GO_VERSION}-alpine AS base
WORKDIR /src
COPY go.mod ./
RUN go mod download

# Development stage with live reload for local Docker Compose usage.
FROM base AS dev
RUN go install -mod=mod github.com/githubnemo/CompileDaemon@v1.4.0
COPY . .
EXPOSE 8080
ENTRYPOINT ["CompileDaemon", "--build=go build -o main .", "--command=./main", "-polling=true"]

# Test stage used by CI or local validation.
FROM base AS test
COPY . .
RUN go test -v ./...

# Build a small static binary.
FROM base AS builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/go-web-app .

# Minimal production runtime. Docker multi-stage builds keep build tools out of the final image.
FROM scratch AS prod
WORKDIR /app
COPY --from=builder /out/go-web-app /app/go-web-app
COPY --from=builder /src/static /app/static
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/app/go-web-app"]
