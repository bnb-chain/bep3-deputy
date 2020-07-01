# Build the app in two stages to optimize build times.
# Stage 1 creates a list of go packages to pass onto the next stage
# Stage 2 downloads and builds the dependencies based on the list, then builds the app.
# This process allows app source files to be modified without triggering a full rebuild of all the dependencies.
FROM golang:1.13-alpine AS dependency_calculator

RUN apk add --no-cache git

WORKDIR /opt/app

# first download dependencies, without source files present to avoid invalidating this layer
COPY go.mod go.sum ./
RUN go mod download

# calculate the list of dependencies based on the source
COPY . .
RUN go list -deps -f '{{ $ip := .Module}}{{if $ip}}{{if not $ip.Main}}{{.ImportPath}}@{{$ip.Version}}{{end}}{{end}}' ./... > deps.txt



FROM golang:1.13-alpine

# Set up apk dependencies
ENV PACKAGES make git libc-dev bash gcc linux-headers eudev-dev curl ca-certificates

# Install minimum necessary dependencies, remove packages
RUN apk add --no-cache $PACKAGES

# Set working directory for the build
# TODO why is this here? what does it do?
WORKDIR /opt/app

# build dependencies
COPY --from=dependency_calculator /opt/app/deps.txt /opt/app/go.mod /opt/app/go.sum ./
# TODO copy in downloaded sources to avoid downloading again
RUN cat deps.txt | xargs go get -v

# build the app
COPY . .
# add the verbose flag to show what packages are being built
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