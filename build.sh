# 在Windows平台下

# 编译Windows amd64
#export GOOS=windows
#export GOARCH=amd64
#
#go build -o cc-win-amd64.exe

### 编译Linux amd64
#export GOOS=linux
#export GOARCH=amd64
#
#go build -o cc-linux-amd64

## 编译Linux arm
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=arm

go build -o ./target/dcs-linux-arm64


