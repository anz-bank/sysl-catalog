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
RUN wget "http://downloads.sourceforge.net/project/plantuml/${PLANTUML_VERSION}/plantuml.${PLANTUML_VERSION}.jar" -O plantuml.jar
RUN wget https://github.com/facebook/nailgun/releases/download/nailgun-all-v1.0.0/nailgun-server-1.0.0-SNAPSHOT.jar
RUN apk add bash
COPY scripts/ /usr/scripts
COPY --from=builder /usr/src/sysl-catalog .
RUN mkdir -p /out
ENV SYSL_PLANTUML=plantuml.jar
ENTRYPOINT ["./scripts/nailgun.sh"]