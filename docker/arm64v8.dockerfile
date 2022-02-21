#BASE
FROM arm64v8/golang:1.16-alpine AS base

WORKDIR /

RUN apk update && apk add --no-cache git curl openssh-client gcc g++ musl-dev

RUN mkdir -p /src

COPY ./ /src/

RUN cd /src && go get ./...

RUN rm -rf /src

#BUILD
FROM base AS build

WORKDIR /

RUN apk update && apk add --no-cache git curl openssh-client gcc g++ musl-dev

RUN mkdir -p /src

COPY ./ /src/

RUN cd /src && go get ./... && go build -o ranker

#FINAL
FROM arm64v8/alpine

ARG version # you could give this a default value as well

RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true
        
COPY --from=build               /src/ranker        / 

ENV PORT=80

ENTRYPOINT ["/ranker"]