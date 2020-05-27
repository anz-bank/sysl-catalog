FROM golang:latest
ENV SYSL_PLANTUML="http://localhost:8080"
COPY . /usr/src
WORKDIR /usr/src
RUN make install