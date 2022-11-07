FROM golang:1.1-alpine3.16 AS go-builder
ARG GAIA_TAG_NAME

# Set up dependencies
ENV PACKAGES curl make git libusb-dev libc-dev bash gcc linux-headers eudev-dev python3

# Install ca-certificates
RUN set -eux; apk add --no-cache ca-certificates build-base;

# Install minimum necessary dependencies
RUN apk add --no-cache $PACKAGES

# Set working directory for the build
WORKDIR /usr/local/app

# Add source files
ADD https://github.com/cosmos/gaia/archive/refs/tags/$GAIA_TAG_NAME.zip /tmp
RUN cd /tmp && unzip $GAIA_TAG_NAME.zip && mv /tmp/gaia-${GAIA_TAG_NAME#?}/* /usr/local/app

# Force it to use static lib (from above) not standard libgo_cosmwasm.so file
RUN LEDGER_ENABLED=false make build

FROM alpine:3.16

COPY --from=go-builder /usr/local/app/build/gaiad /usr/bin/gaiad

# Set up dependencies
ENV PACKAGES curl make bash jq
ENV CHAIN_BIN /usr/bin/gaiad

# Install minimum necessary dependencies
RUN apk add --no-cache $PACKAGES

WORKDIR /opt

# rest server, p2p, rpc, grpc
EXPOSE 1317 26656 26657 9090

# Run persistenceCore by default, omit entrypoint to ease using container with cli
CMD ["/usr/bin/gaiad", "version", "--long"]
