# 基础镜像
FROM golang:1.20.5 As goBuild

# 工作目录
WORKDIR /app

# 复制源代码到容器中
COPY . .

ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64

# 编译应用程序
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod tidy
RUN go build -o app .


FROM node:14.12.0 As nodeBuild

# 工作目录
WORKDIR /app

COPY --from=goBuild /app/vue /app

RUN npm config set registry https://registry.npm.taobao.org/
RUN npm config set sass_binary_site https://npm.taobao.org/mirrors/node-sass/
RUN npm install
RUN npm run build

FROM ubuntu:latest

ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update
RUN apt-get install -y sqlite3 libsqlite3-dev

RUN apt install -y tzdata --no-install-recommends \
  &&  echo "Asia/Shanghai" > /etc/timezone  \
  && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
  && dpkg-reconfigure -f noninteractive tzdata

COPY --from=goBuild /app/app /app
COPY --from=goBuild /app/em.db /em.db
COPY --from=goBuild /app/config /config
COPY --from=goBuild /app/selfpara /selfpara
COPY --from=nodeBuild /app/webroot /webroot

VOLUME /selfpara

# 暴露端口
EXPOSE 7070

# 启动应用程序
ENTRYPOINT  ["./app"]