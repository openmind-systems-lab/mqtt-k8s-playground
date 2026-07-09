FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o mqtt-client main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/mqtt-client .
CMD ["./mqtt-client"]
