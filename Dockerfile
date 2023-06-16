# 基础镜像
FROM golang:1.19 As build

# 工作目录
WORKDIR /app

# 复制源代码到容器中
COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=arm

# 编译应用程序
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod tidy
RUN go build -o app .

FROM scratch

COPY --from=build /app/app /app
COPY --from=build /app/application.yml /

# 暴露端口
EXPOSE 8181

# 启动应用程序
ENTRYPOINT  ["./app"]