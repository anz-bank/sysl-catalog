#FROM golang:latest AS builder
#COPY . /usr/src
#WORKDIR /usr/src
#RUN go build -o sysl-catalog .

FROM openjdk:14-jdk-alpine3.10
MAINTAINER think@hotmail.de
WORKDIR /usr/src
ENV PLANTUML_VERSION=1.2019.11
COPY sysl-catalog .
ENV LANG en_US.UTF-8
RUN \
  apk add --no-cache graphviz wget ca-certificates && \
  apk add --no-cache graphviz wget ca-certificates ttf-dejavu fontconfig && \
  wget "http://downloads.sourceforge.net/project/plantuml/1.2020.10/plantuml.1.2020.10.jar" -O plantuml.jar && \
  apk del wget ca-certificates
RUN apk add bash
#RUN apk add ls
RUN ["java", "-Djava.awt.headless=true", "-jar", "plantuml.jar", "-version"]
RUN ["dot", "-version"]
RUN mkdir -p /out
#COPY --from=builder /usr/src/sysl-catalog /usr/src/sysl-catalog

ENV SYSL_PLANTUML=java
ENTRYPOINT ["./sysl-catalog", "-o", "/out/"]
#CMD ls -a demo/