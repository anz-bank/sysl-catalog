FROM plantuml/plantuml-server:latest
FROM golang:latest
ENV SYSL_PLANTUML="http://localhost:8080"
COPY . /usr/src
WORKDIR /usr/src
RUN go build .
CMD ["./sysl-catalog", "demo/simple.sysl", "--serve"]