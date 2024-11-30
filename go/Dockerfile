# [https://docs.docker.com/develop/develop-images/multistage-build/]
# [https://docs.docker.com/reference/dockerfile]

##
## Build Stage
##

# FROM [--platform=<platform>] <image> [AS <name>]
# The FROM instruction initialises a new build stage and sets the base image for subsequent instructions.
# FROM can appear multiple times within a single Dockerfile to create multiple images or use one build stage as a dependency for another.
# Optionally a name can be given to a new build stage by adding AS name to the FROM instruction. The name can be used in subsequent FROM <name>,
# COPY --from=<name>, and RUN --mount=type=bind,from=<name> instructions to refer to the image built in this stage.
#FROM golang:1.23 AS builder
FROM golang:alpine3.20 AS builder

# RUN [OPTIONS] <command> ...
# The RUN instruction will execute any commands to create a new layer on top of the current image. The added layer is used in the next step in
# the Dockerfile.
RUN apk add --no-cache \
    # Important: required for go-sqlite3
    gcc \
    # Required for Alpine
    musl-dev

# WORKDIR /path/to/workdir
# The WORKDIR instruction sets the working directory for any RUN, CMD, ENTRYPOINT, COPY and ADD instructions that follow it in the Dockerfile.
# If not specified, the default working directory is /.
WORKDIR /workspace

# COPY [OPTIONS] <src> ... <dst>
# The COPY instruction copies new files or directories from <src> and adds them to the filesystem of the image at the path <dest>.
# Files and directories can be copied from the build context, build stage, named context, or an image.
COPY go.mod .
COPY go.sum .
COPY *.go .
COPY api/ api/
COPY config/ config/
COPY server/ server/
COPY store/ store/
COPY www/ /www/
COPY server_cert.pem /
COPY server_privkey.pem /

RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -ldflags='-s -w -extldflags "-static"' -o /microservice

##
## Deploy Stage
##

# gcr.io/distroless/static:nonroot is a Docker image provided by Google as part of their Distroless project.
# gcr.io - is the Google Container Registry, where the image is hosted.
# distroless - signifies that the image is part of the Distroless project, which aims to create minimal container images that contain only
# the necessary components to run an application.
# static - indicates that the image is based on a statically linked libc, making it suitable for applications that are statically compiled.
# nonroot - specifies that the image is designed to run as a non-root user, enhancing security by preventing the application from running
# with elevated priviliges.
FROM gcr.io/distroless/static:nonroot

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /microservice microservice
COPY --from=builder --chown=nonroot:nonroot /www/ www/
COPY --from=builder --chown=nonroot:nonroot /server_cert.pem server_cert.pem
COPY --from=builder --chown=nonroot:nonroot /server_privkey.pem server_privkey.pem

# USER UID[ :GID]
# The USER instruction sets the user name (or UID) and optionally the user group (or GID) to use as the default user and group for the remainder
# of the current stage.
# The specified user is used for RUN instructions and at runtime, runs the relevant ENTRYPOINT and CMD commands.
USER root:root

# ENTRYPOINT ["executable", "param1", "param2"]
# An ENTRYPOINT allows you to configure a container that will run as an executable.
ENTRYPOINT ["/microservice", "-s"]
