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

# GIT配置
RUN apt-get install -y git
COPY ssh /root/.ssh
RUN chmod 400 ~/.ssh/id_ed25519
RUN ssh-keyscan gitee.com >> /root/.ssh/known_hosts

# 编译环境安装
RUN apt-get install git wget flex bison gperf python3 python3-venv python3-setuptools ninja-build ccache libffi-dev libssl-dev dfu-util libusb-1.0-0 -y
RUN mkdir /app/esp 
COPY lib/esp-idf /app/esp/esp-idf
COPY lib/esp-gitee-tools /app/esp/esp-gitee-tools
RUN cd /app/esp/esp-gitee-tools && ./submodule-update.sh ../esp-idf/
COPY lib/.espressif /root/.espressif
COPY lib/cmake-3.24.1.tar.gz /app/cmake-3.24.1.tar.gz
RUN apt-get install build-essential -y && cd /app && tar -zxvf cmake-3.24.1.tar.gz 
RUN cd /app/cmake-3.24.1 && ./bootstrap && make && make install


# 执行文件拷贝
COPY --from=builder /app/main /app

# 配置文件拷贝
COPY config /app/config

# 创建文件夹
RUN mkdir /app/tmp
RUN mkdir /app/log

# 启动命令
CMD ["./main"]



