FROM golang:alpine AS builder
ADD . /src
RUN cd /src && go build -o sysl-catalog

FROM alpine:3.10 AS nailgun
WORKDIR /usr
RUN apk add build-base
ADD https://raw.githubusercontent.com/facebook/nailgun/master/nailgun-client/c/ng.c .
RUN gcc -Wall ng.c -o ng

FROM openjdk:14-jdk-alpine3.10
WORKDIR /usr
ENV LANG en_US.UTF-8

ADD java/plantuml.jar .
ADD java/nailgun-server-1.0.0-SNAPSHOT.jar .
COPY --from=nailgun /usr/ng .
COPY --from=builder /src/sysl-catalog /bin/
COPY scripts/nailgun.sh .
RUN apk add --no-cache --upgrade bash
RUN apk add --no-cache graphviz wget ca-certificates && \
      apk add --no-cache graphviz wget ca-certificates ttf-dejavu fontconfig
RUN apk add git
RUN mkdir -p /out
RUN apk add --no-cache --upgrade git
ENV SYSL_PLANTUML=plantuml.jar
ENTRYPOINT ["./nailgun.sh"]