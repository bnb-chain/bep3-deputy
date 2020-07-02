# This dockerfile has been optimized to reduce build times.
# This is achieved by downloading and building dependencies first in a way that can be cached despite changes in the deputy code.
# Deputy source files can be modified without triggering a full rebuild of all the dependencies.
# Stage 1 downloads dependencies and creates a list of required go packages.
# Stage 2 copies in the dependency sources and builds the dependencies based on the list. After this source code is copied in and the deputy is built.

FROM golang:1.13-alpine AS depCalculator

RUN apk add --no-cache git

WORKDIR /opt/app

# download dependencies' source code (required for go list to work)
COPY go.mod go.sum ./
RUN go mod download

# calculate the list of dependencies based on the deputy code
COPY . .
# the formatting directive outputs a package's import path plus '@' plus the package's module's version
# these can be passed to `go get` to download and build the package
RUN go list -deps -f '{{ $ip := .Module}}{{if $ip}}{{if not $ip.Main}}{{.ImportPath}}@{{$ip.Version}}{{end}}{{end}}' ./... > deps.txt



FROM golang:1.13-alpine

# Set up apk dependencies
ENV PACKAGES make git libc-dev bash gcc linux-headers eudev-dev curl ca-certificates

# Install minimum necessary dependencies, remove packages
RUN apk add --no-cache $PACKAGES

# Set working directory for the build
WORKDIR /opt/app

# copy in downloaded modules to avoid downloading again
COPY --from=depCalculator $GOPATH $GOPATH
# build dependencies (go.mod needs to be there so go knows it's working inside a module)
COPY --from=depCalculator /opt/app/deps.txt /opt/app/go.mod /opt/app/go.sum ./
RUN cat deps.txt | xargs go get -v

# build the deputy (add the verbose flag to show what packages are being built)
COPY . .
ENV BUILD_FLAGS '-v'
RUN make build


ENV DEPUTY_HOME /deputy
ENV BNB_NETWORK 1
ENV CONFIG_FILE_PATH $DEPUTY_HOME/config/config.json
ENV CONFIG_TYPE "local"
# You need to specify aws s3 config if you want to load config from s3
ENV AWS_REGION ""
ENV AWS_SECRET_KEY ""

# Run as non-root user for security
# USER 1000 # TODO

VOLUME [ $DEPUTY_HOME ]

# Run the app
CMD ./build/deputy --bnb-network $BNB_NETWORK --config-type $CONFIG_TYPE --config-path $CONFIG_FILE_PATH --aws-region $AWS_REGION --aws-secret-key $AWS_SECRET_KEY