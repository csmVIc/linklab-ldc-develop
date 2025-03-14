FROM golang:1.15.3-buster AS builder 
# ENV GOPROXY=https://goproxy.io
ENV GOPROXY=https://goproxy.cn,direct

RUN sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list
RUN sed -i 's|security.debian.org/debian-security|mirrors.ustc.edu.cn/debian-security|g' /etc/apt/sources.list
RUN apt-get update && apt-get install -y unzip 
LABEL maintainer="Yang Guang <yangg.china@outlook.com>"
WORKDIR /app
COPY . /app
RUN unzip -d ../base-library base-library.zip
RUN go build -o main .

FROM debian:buster-slim
USER root
RUN sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list
RUN sed -i 's|security.debian.org/debian-security|mirrors.ustc.edu.cn/debian-security|g' /etc/apt/sources.list
RUN apt-get clean && apt-get update && apt-get install -y locales && rm -rf /var/lib/apt/lists/* \
    && localedef -i en_US -c -f UTF-8 -A /usr/share/locale/locale.alias en_US.UTF-8
ENV LANG en_US.utf8
RUN apt-get update && \
    apt-get install -y bash ca-certificates tzdata
ENV TZ="Asia/Shanghai"
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
LABEL maintainer="Yang Guang <yangg.china@outlook.com>"

# 工作目录准备
WORKDIR /app

# 执行文件拷贝
COPY --from=builder /app/main /app

# 配置文件拷贝
COPY config /app/config
COPY yaml /app/yaml

# 创建文件夹
RUN mkdir /app/tmp
RUN mkdir /app/log

# 启动命令
CMD ["./main"]

