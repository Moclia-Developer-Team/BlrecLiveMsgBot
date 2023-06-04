FROM docker.io/golang:1.20.4-bullseye AS builder
ADD . /src
WORKDIR /src
ENV GOOS=linux GOARCH=amd64
RUN mkdir /src/build && \
    go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://goproxy.cn,direct && \
    go build -v -o /src/build/BlrecLivePostBot && \
    chmod +x /src/build/BlrecLivePostBot

FROM docker.io/debian:bookworm
RUN sed -i 's/deb.debian.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apt/sources.list.d/debian.sources && \
    apt update && \
    apt install apt-transport-https ca-certificates wget unzip openssl tar curl -y && \
    sed -i 's/http/https/g' /etc/apt/sources.list.d/debian.sources && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai  /etc/localtime && \
    mkdir /work
COPY --from=builder /src/build/BlrecLivePostBot /work/BlrecLivePostBot
WORKDIR /work
CMD [ "/work/BlrecLivePostBot" ]
EXPOSE 2023
HEALTHCHECK --interval=30s --timeout=3s CMD curl -fs http://127.0.0.1:2023/api/v1/system/healthcheck | exit 1