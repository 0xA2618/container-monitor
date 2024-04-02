FROM golang:1.22.1

WORKDIR /app

COPY src/ .

RUN go env -w GOPROXY='https://goproxy.cn,direct' && \
    go build cmd/main.go
CMD ["./main"]