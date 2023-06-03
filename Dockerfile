FROM docker.io/debian:bookworm
RUN sed -i 's/deb.debian.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apt/sources.list.d/debian.sources && \
    apt update && \
    apt install apt-transport-https ca-certificates wget unzip openssl tar curl -y && \
    sed -i 's/http/https/g' /etc/apt/sources.list.d/debian.sources && \
    mkdir /work
COPY ./BlrecLivePostBot /work
WORKDIR /work
ENTRYPOINT /work/BlrecLivePostBot
HEALTHCHECK --interval=30s --timeout=3s CMD curl -fs http://127.0.0.1:2023/api/v1/system/healthcheck | exit 1