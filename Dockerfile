FROM golang:1.22-alpine 

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download 

COPY . .

# 添加一些基本的调试输出
RUN echo "Starting build..."
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/main.go 
RUN echo "Build completed"

EXPOSE 6061

# 添加日志输出
CMD echo "Starting application..." && /app/main