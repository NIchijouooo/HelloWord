# 基础镜像
FROM golang:1.19 As build

# 工作目录
WORKDIR /app

# 复制源代码到容器中
COPY . .

ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64

# 编译应用程序
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod tidy
RUN go build -o app .

FROM ubuntu:latest

RUN apt-get update && apt-get install -y sqlite3 libsqlite3-dev

COPY --from=build /app/app /app
COPY --from=build /app/config /config

# 暴露端口
EXPOSE 7070

# 启动应用程序
ENTRYPOINT  ["./app"]