FROM golang:1.20 as builder

WORKDIR /apps

COPY ./ /apps
RUN export GOPROXY=https://goproxy.cn \
    && go build  -ldflags "-s -w" -o ocserv_exporter  main.go \
    && chmod +x ocserv_exporter

FROM alpine
LABEL maintainer="tchua"
COPY --from=builder /apps/ocserv_exporter  /apps/
COPY --from=builder /apps/etc/  /apps/etc/
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN echo -e  "http://mirrors.aliyun.com/alpine/v3.15/main\nhttp://mirrors.aliyun.com/alpine/v3.15/community" >  /etc/apk/repositories \
&& apk update && apk add tzdata nmap-ncat \
&& cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
&& echo "Shanghai/Asia" > /etc/timezone \
&& apk del tzdata

WORKDIR /apps

EXPOSE 18086

CMD ["./ocserv_exporter"]