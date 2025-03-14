FROM golang:1.15.3-buster AS builder 
ENV GOPROXY=https://goproxy.io
RUN sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list
RUN sed -i 's|security.debian.org/debian-security|mirrors.ustc.edu.cn/debian-security|g' /etc/apt/sources.list
RUN apt-get update && apt-get install -y locales && rm -rf /var/lib/apt/lists/* \
    && localedef -i en_US -c -f UTF-8 -A /usr/share/locale/locale.alias en_US.UTF-8
ENV LANG en_US.utf8
RUN apt-get update && \
    apt-get install -y unzip
LABEL maintainer="Yang Guang <yangg.china@outlook.com>"

# 工作目录准备
WORKDIR /app

# 拷贝基础库
COPY base-library.zip ./
RUN unzip -d ../base-library base-library.zip

# go语言所需库下载
COPY go.mod go.sum ./
RUN go mod download

# 拷贝源文件
COPY main.go .
COPY driver /app/driver

# 编译
RUN go build -o main .

FROM debian:buster-slim
USER root
RUN sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list
RUN sed -i 's|security.debian.org/debian-security|mirrors.ustc.edu.cn/debian-security|g' /etc/apt/sources.list
RUN apt-get update && apt-get install -y locales && rm -rf /var/lib/apt/lists/* \
    && localedef -i en_US -c -f UTF-8 -A /usr/share/locale/locale.alias en_US.UTF-8
ENV LANG en_US.utf8
RUN apt-get update && \
    apt-get install -y bash ca-certificates tzdata unzip
ENV TZ="Asia/Shanghai"
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
LABEL maintainer="Yang Guang <yangg.china@outlook.com>"

# 工作目录准备
WORKDIR /app

# AliOS-Things编译环境安装
RUN apt-get install -y git python python-pip
RUN python -m pip install --trusted-host=mirrors.aliyun.com -i https://mirrors.aliyun.com/pypi/simple/ --upgrade pip
RUN pip install --trusted-host=mirrors.aliyun.com -i https://mirrors.aliyun.com/pypi/simple/ setuptools
RUN pip install --trusted-host=mirrors.aliyun.com -i https://mirrors.aliyun.com/pypi/simple/ wheel
RUN pip install --trusted-host=mirrors.aliyun.com -i https://mirrors.aliyun.com/pypi/simple/ aos-cube

# GIT配置
RUN apt-get install -y git
COPY ssh /root/.ssh
RUN ssh-keyscan gitee.com >> /root/.ssh/known_hosts

# 执行文件拷贝
COPY --from=builder /app/main /app

# 配置文件拷贝
COPY config /app/config

# 创建文件夹
RUN mkdir /app/tmp
RUN mkdir /app/log

# 启动命令
CMD ["./main"]



