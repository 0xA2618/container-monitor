# 阶段一：构建阶段
FROM golang:1.22.1 AS builder

WORKDIR /app

COPY src/ .

RUN go env -w GOPROXY='https://goproxy.cn,direct' && \
    go build cmd/main.go

# 阶段二：最终镜像
FROM golang:1.22.1

WORKDIR /app

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/main .

CMD ["./main"]