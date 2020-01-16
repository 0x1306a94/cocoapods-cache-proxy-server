FROM alpine
LABEL maintainer="0x1306a94 0x1306a94@gmail.com"

ENV COCOAPODS_CACHE_PROXY_USER=admin \
    COCOAPODS_CACHE_PROXY_PASSWORD=admin11 \
    COCOAPODS_CACHE_PROXY_CACHE_DIR=/app/cococapods/cache

RUN mkdir -p /app/bin /app/cococapods/cache

ADD https://github.com/0x1306a94/cocoapods-cache-proxy-server/releases/download/v2/cocoapods-cache-proxy-server-linux-amd64 /app/bin

WORKDIR /app

EXPOSE 9898

CMD ["./bin/cocoapods-cache-proxy-server-linux-amd64"]