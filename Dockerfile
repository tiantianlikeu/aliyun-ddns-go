FROM golang:1.20.2-alpine3.16 AS builder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /aliyun-ddns

ADD . /aliyun-ddns

RUN cd /aliyun-ddns/ && go build -ldflags="-s -w" -o /aliyun-ddns/main




FROM alpine:3.16
MAINTAINER "ttliku"
VOLUME /tmp
ENV LANG C.UTF-8
#安装时区数据 tzdata
RUN echo -e  "http://mirrors.aliyun.com/alpine/v3.16/main\nhttp://mirrors.aliyun.com/alpine/v3.16/community" >  /etc/apk/repositories \
apk update && apk add tzdata \
&& cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
&& echo "Shanghai/Asia" > /etc/timezone \
&& apk del tzdata

WORKDIR /aliyun-ddns
RUN mkdir -p /aliyun-ddns/conf
COPY --from=builder /aliyun-ddns/main /aliyun-ddns/main
COPY --from=builder /aliyun-ddns/config.yaml /aliyun-ddns/conf/config.yaml
EXPOSE 8080

CMD /aliyun-ddns/main


