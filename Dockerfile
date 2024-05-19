FROM golang:alpine

# 为镜像设置环境变量
ENV GO111MODULE=ON \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOARCH=amd64

# 移动到工作路径
WORKDIR /build

COPY . .
# 编译代码
RUN go build -o app .

#移动到用户存放生成的二进制文件
WORKDIR /dist

#复制二进制文件
RUN cp /build/app .

#暴露端口
EXPOSE 8082

#启动容器时运行的命令
CMD ["/dist/app"]