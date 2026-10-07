# 构建阶段
FROM golang:1.22-alpine AS builder

WORKDIR /app

# 复制后端代码
COPY survey-system/backend ./

# 下载依赖
RUN go mod download

# 编译
RUN go build -o main .

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 从构建阶段复制编译好的二进制
COPY --from=builder /app/main .

# 暴露端口
EXPOSE 8080

# 启动
CMD ["./main"]
