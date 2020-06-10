FROM golang:alpine AS builder
ADD . /src
RUN cd /src && go build -o sysl-catalog

FROM gcc:latest AS nailgun
WORKDIR /usr
ADD https://raw.githubusercontent.com/facebook/nailgun/master/nailgun-client/c/ng.c .
RUN gcc -Wall ng.c -o ng

FROM openjdk:14-jdk-alpine3.10
WORKDIR /usr
ENV LANG en_US.UTF-8
RUN apk add --no-cache graphviz wget ca-certificates && \
      apk add --no-cache graphviz wget ca-certificates ttf-dejavu fontconfig
ADD java/plantuml.jar .
ADD java/nailgun-server-1.0.0-SNAPSHOT.jar .
COPY scripts/nailgun.sh .
COPY --from=nailgun /usr/ng .
COPY --from=builder /src/sysl-catalog .
RUN mkdir -p /out
ENV SYSL_PLANTUML=plantuml.jar
ENTRYPOINT ["./nailgun.sh"]