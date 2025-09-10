FROM golang:1.24.5-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o server ./cmd

FROM alpine
WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY configs/config.yaml /app/configs/config.yaml
COPY --from=builder --chmod=755 /app/server /app/server
RUN ls -la /app
ENTRYPOINT ["/app/server"]