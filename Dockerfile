FROM golang:1.22-alpine 

# 添加调试工具和时区数据
RUN apk add --no-cache \
    curl \
    tzdata \
    bash

# 设置时区
ENV TZ=Asia/Shanghai

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download 

COPY . .

# 显示构建环境信息
RUN echo "Starting build process..." && \
    pwd && \
    ls -la && \
    echo "Go version:" && \
    go version

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/main.go 
RUN echo "Build completed"

EXPOSE 6061

# 使用 bash 来运行应用，添加更多调试信息
CMD ["/bin/bash", "-c", "\
    echo '=== Environment Variables ===' && \
    printenv && \
    echo '=== Current Directory ===' && \
    pwd && \
    ls -la && \
    echo '=== Starting Application ===' && \
    /app/main"]