
ARG REGISTRY=""
ARG GO_BUILD_IMAGE=golang:latest
ARG BASE_IMAGE=think/plantuml:latest
FROM ${REGISTRY:+${REGISTRY}/}${GO_BUILD_IMAGE} as builder
ARG VERSION=dev
ARG REPO_URL=
ARG COMMIT_HASH=
ARG CONTAINER_TAG=
ARG SEMVER=
ARG GOPROXY=""

# Create the user and group files that will be used in the running container to
# run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

WORKDIR /build
COPY go.mod .
COPY go.sum .

# fetch dependencies
ENV GOPROXY=${GOPROXY}
RUN go mod download

COPY . .

# build app
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build \
	-ldflags=" \
		-X main.version=${SEMVER} \
	" \
	-o /build/sysl-catalog \
	./main.go

FROM ${REGISTRY:+${REGISTRY}/}${BASE_IMAGE}
ENV SYSL_PLANTUML="http://www.plantuml.com/plantuml"
WORKDIR /usr/src
COPY --from=builder ./build/sysl-catalog .
RUN ["java", "-Djava.awt.headless=true", "-jar", "/plantuml.jar", "-version"]
RUN ["dot", "-version"]

ENTRYPOINT ["/usr/src/sysl-catalog"]