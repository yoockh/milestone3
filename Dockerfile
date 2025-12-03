# Builder stage (Debian-less official golang tag)
FROM golang:1.25 AS builder
RUN apt-get update && apt-get install -y --no-install-recommends git ca-certificates \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o /app/server ./be/app

# Final stage: small distro image
FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=builder /app/server /app/server
USER nonroot

ENV PORT=8080
EXPOSE 8080
CMD ["/app/server"]
