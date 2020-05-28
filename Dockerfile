#FROM golang:latest AS builder
#COPY example /usr/src
#WORKDIR /usr/src
#RUN go build -o sysl-catalog .

# Compile stage
#FROM golang:latest AS builder
#
## Build Delve
#RUN go get github.com/go-delve/delve/cmd/dlv
#
#ADD . /dockerdev
#WORKDIR /dockerdev
#
#RUN go build -gcflags="all=-N -l" -o sysl-catalog

# Final stage
#FROM debian:buster
#
#EXPOSE 8000 40000
#
#WORKDIR /
#COPY --from=builder /go/bin/dlv /
#COPY --from=builder /dockerdev/sysl-catalog /
#
#CMD ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/sysl-catalog"]
#

FROM openjdk:14-jdk-alpine3.10
MAINTAINER think@hotmail.de
WORKDIR /usr/src
ENV PLANTUML_VERSION=1.2019.11
ENV SYSL_PLANTUML=https://plantuml.com/plantuml
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
#COPY --from=builder /dockerdev/sysl-catalog /usr/src/sysl-catalog
COPY sysl-catalog /usr/src/sysl-catalog
#ENTRYPOINT ["./usr/src/sysl-catalog"]
CMD ./sysl-catalog demo/simple2.sysl -o /out/
#ENTRYPOINT ["/bin/sh ", "./sysl-catalog demo/simple2.sysl -o /out/"]
#CMD ["--help"]
#E
#ENTRYPOINT ["java", "-Djava.awt.headless=true", "-jar", "plantuml.jar", "-p"]
#CMD ["-tsvg"]

#java -Djava.awt.headless=true -jar plantuml.jar -version