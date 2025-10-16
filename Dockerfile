# 多阶段构建 - 构建阶段
FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS builder

# 声明构建参数
ARG TARGETOS
ARG TARGETARCH

# 设置工作目录
WORKDIR /app

# 安装必要的工具
RUN apk add --no-cache git ca-certificates tzdata

# 复制go模块文件和源代码
COPY go.mod *.go ./

# 构建应用程序
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -installsuffix cgo -o infopush .

# 运行阶段
FROM alpine:latest

# 安装ca-certificates和tzdata
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/infopush .

# 复制 data 目录
COPY data/ /app/data/

# 修改文件权限
RUN chmod +x infopush

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["./infopush"] 