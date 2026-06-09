# ========== Build Stage ==========
FROM golang:1.23-alpine AS builder

WORKDIR /build

# 先复制依赖文件，利用缓存加速
COPY server-go/go.mod server-go/go.sum ./
RUN go mod download

# 复制全部源码
COPY server-go/ .

# 编译（静态链接，适合 alpine 镜像）
RUN CGO_ENABLED=0 go build -o suisui .

# ========== Run Stage ==========
FROM alpine:latest

# 安装 ca-certificates（用于 HTTPS 请求）
RUN apk --no-cache add ca-certificates

WORKDIR /data

# 从构建阶段复制二进制
COPY --from=builder /build/suisui .

# 创建上传目录
RUN mkdir -p uploads

EXPOSE 3001

VOLUME ["/data/suisui.db", "/data/uploads"]

CMD ["./suisui"]
