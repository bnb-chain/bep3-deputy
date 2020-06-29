FROM golang:1.12-alpine

# Set up apk dependencies
ENV PACKAGES make git libc-dev bash gcc linux-headers eudev-dev curl ca-certificates

# Install minimum necessary dependencies, remove packages
RUN apk add --no-cache $PACKAGES

# Set working directory for the build
# TODO why is this here? what does it do?
WORKDIR /opt/app

# download dependencies first to avoid downloading again when rebuilding
COPY go.mod .
COPY go.sum .
RUN go mod download

# Add source files and build
COPY . .
RUN make build


ENV DEPUTY_HOME /deputy
ENV BNB_NETWORK 1
ENV CONFIG_FILE_PATH $DEPUTY_HOME/config/config.json
ENV CONFIG_TYPE "local"
# You need to specify aws s3 config if you want to load config from s3
ENV AWS_REGION ""
ENV AWS_SECRET_KEY ""

# Run as non-root user for security
USER 1000

VOLUME [ $DEPUTY_HOME ]

# Run the app
CMD ./build/deputy --bnb-network $BNB_NETWORK --config-type $CONFIG_TYPE --config-path $CONFIG_FILE_PATH --aws-region $AWS_REGION --aws-secret-key $AWS_SECRET_KEY