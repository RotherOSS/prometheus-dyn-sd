FROM golang:alpine3.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o prometheus-dyn-sd main.go

FROM alpine:latest  
WORKDIR /app
COPY --from=builder /app/prometheus-dyn-sd .
CMD ["./prometheus-dyn-sd"]