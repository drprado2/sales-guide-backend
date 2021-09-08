FROM kong/kong-gateway:2.5.0.0-alpine

USER root

ENV LUA_VERSION 5.1

RUN apk update

RUN apk add --no-cache git

RUN apk add lua${LUA_VERSION}
RUN apk add lua${LUA_VERSION}-dev

RUN apk add bash build-base curl git openssl unzip

RUN cd /tmp && \
    git clone https://github.com/keplerproject/luarocks.git && \
    cd luarocks && \
    sh ./configure && \
    make build install && \
    cd && \
    rm -rf /tmp/luarocks

RUN luarocks install busted
RUN luarocks install luacheck
RUN luarocks install luacov

RUN luarocks install kong-jwt2header