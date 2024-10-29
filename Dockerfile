FROM golang:1.22.1 AS builder
WORKDIR /app
COPY src/ .
RUN go env -w GOPROXY='https://goproxy.cn,direct' 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main cmd/main.go

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
CMD ["./main"]
