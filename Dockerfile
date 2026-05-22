FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build -o /deploy-agent ./cmd/server

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /deploy-agent /usr/local/bin/deploy-agent
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/deploy-agent"]
