FROM golang:latest AS builder
COPY . /usr/src
WORKDIR /usr/src
RUN go build -o sysl-catalog .

FROM openjdk:14-jdk-alpine3.10
WORKDIR /in
ENV PLANTUML_VERSION=1.2019.10
COPY sysl-catalog .
COPY plantuml.jar .
ENV LANG en_US.UTF-8
RUN \
  apk add --no-cache graphviz wget ca-certificates && \
  apk add --no-cache graphviz wget ca-certificates ttf-dejavu fontconfig && \
  wget "http://downloads.sourceforge.net/project/plantuml/${PLANTUML_VERSION}/plantuml.${PLANTUML_VERSION}.jar" -O plantuml.jar && \
  apk del wget ca-certificates
RUN apk add bash
RUN ["java", "-Djava.awt.headless=true", "-jar", "plantuml.jar", "-version"]
RUN ["dot", "-version"]
RUN mkdir -p /out
ENV SYSL_PLANTUML=plantuml.jar
ENTRYPOINT ["./sysl-catalog", "-o", "/out/"]
