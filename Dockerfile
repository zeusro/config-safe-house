FROM golang:1.12.6 AS build-env
ADD . /src/github.com/zeusro/config-safe-house
ENV GOPATH /:/src/github.com/zeusro/config-safe-house/vendor
WORKDIR /src/github.com/zeusro/config-safe-house
RUN apt-get update -y && apt-get install gcc 
RUN make

FROM alpine:3.10
COPY --from=build-env /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=build-env /src/github.com/zeusro/config-safe-house/config-safe-house /
# 中国特色社会主义
RUN echo https://mirrors.ustc.edu.cn/alpine/v3.10/main > /etc/apk/repositories; \
    echo https://mirrors.ustc.edu.cn/alpine/v3.10/community >> /etc/apk/repositories;\
    echo "Asia/Shanghai" > /etc/timezone ;\
    apk add --no-cache bash 
WORKDIR /
ENTRYPOINT ["/config-safe-house"]