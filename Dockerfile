FROM golang:latest AS builder
COPY . /usr/src
WORKDIR /usr/src
RUN go build -o sysl-catalog .

FROM openjdk:14-jdk-alpine3.10
WORKDIR /usr/
ENV PLANTUML_VERSION=1.2019.10
ENV LANG en_US.UTF-8
RUN \
  apk add --no-cache graphviz wget ca-certificates && \
  apk add --no-cache graphviz wget ca-certificates ttf-dejavu fontconfig && \
  apk add readline && \
  apk del wget ca-certificates
RUN apk add bash
COPY scripts/ /usr/scripts
COPY --from=builder /usr/src/sysl-catalog .
COPY java/plantuml.jar .
COPY java/nailgun-server-1.0.0-SNAPSHOT.jar .
RUN mkdir -p /out
ENV SYSL_PLANTUML=plantuml.jar
ENTRYPOINT ["./sysl-catalog", "-o", "/out/"]