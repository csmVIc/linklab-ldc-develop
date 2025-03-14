FROM golang:1.15.3-alpine AS builder 
ENV GOPROXY=https://goproxy.io
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk update && apk upgrade && \
    apk add --no-cache unzip 
LABEL maintainer="Yang Guang <yangg.china@outlook.com>"
WORKDIR /app
COPY . /app
RUN unzip -d ../base-library base-library.zip
RUN go build -o main .

FROM alpine 
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk update && apk upgrade && \
    apk add --no-cache bash ca-certificates tzdata
ENV TZ="Asia/Shanghai"
LABEL maintainer="Yang Guang <yangg.china@outlook.com>"

# 工作目录准备
WORKDIR /app

# 执行文件拷贝
COPY --from=builder /app/main /app

# 配置文件拷贝
COPY config /app/config

# 创建文件夹
RUN mkdir /app/log

CMD ["./main"]
