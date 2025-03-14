FROM golang:1.15.3-alpine AS builder 
ENV GOPROXY=https://goproxy.io
# /etc/apk/repositories - 这是要修改的目标文件,它是 Alpine Linux 的软件源配置文件
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk update && apk upgrade && \
    apk add --no-cache unzip 
LABEL maintainer="Yang Guang <yangg.china@outlook.com>"
WORKDIR /app
COPY . /app
RUN unzip -d ../base-library base-library.zip
# go build构建当前目录的Go应用，生成main
RUN go build -o main .

FROM docker:stable
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk update && apk upgrade && \
    apk add --no-cache bash ca-certificates tzdata
# 设置时区
ENV TZ="Asia/Shanghai"
LABEL maintainer="Yang Guang <yangg.china@outlook.com>"

# 工作目录准备
WORKDIR /app

# 执行文件拷贝，将第一阶段的builder环境中的/app目录下main可执行文件copy到当前目录/app
COPY --from=builder /app/main /app

# 创建文件夹
RUN mkdir /app/tmp
RUN mkdir /app/log
RUN mkdir /app/workspace

# 启动命令
CMD ["./main"]

