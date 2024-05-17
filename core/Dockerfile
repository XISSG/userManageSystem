FROM golang:latest

# 设置工作目录
WORKDIR /app

# 复制当前目录下的所有文件到容器中的 /app 目录下
COPY tmp/test .

# 设置环境变量，用于 Go 编译器的工作目录
ENV GOPATH /app

# 指定默认的工作目录
WORKDIR /app
