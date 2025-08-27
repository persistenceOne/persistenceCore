FROM golang:1.23.6-alpine3.21 AS go-builder

# Set up dependencies
ENV PACKAGES curl make git libusb-dev libc-dev bash gcc linux-headers eudev-dev python3

# Install ca-certificates
RUN set -eux; apk add --no-cache ca-certificates build-base;

# Install minimum necessary dependencies
RUN apk add --no-cache $PACKAGES

# See https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v2.2.4/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
ADD https://github.com/CosmWasm/wasmvm/releases/download/v2.2.4/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a
RUN sha256sum /lib/libwasmvm_muslc.aarch64.a | grep 27fb13821dbc519119f4f98c30a42cb32429b111b0fdc883686c34a41777488f
RUN sha256sum /lib/libwasmvm_muslc.x86_64.a | grep 70c989684d2b48ca17bbd55bb694bbb136d75c393c067ef3bdbca31d2b23b578

# Copy the library you want to the final location that will be found by the linker flag `-lwasmvm_muslc`
RUN cp /lib/libwasmvm_muslc.$(uname -m).a /lib/libwasmvm_muslc.a

# Set working directory for the build
WORKDIR /usr/local/app

# Add source files
COPY . .

# Force it to use static lib (from above) not standard libgo_cosmwasm.so file
RUN LEDGER_ENABLED=true BUILD_TAGS="muslc linkstatic" LDFLAGS='-linkmode external -extldflags "-static"' make build
RUN echo "Ensuring binary is statically linked ..." \
  && (file /usr/local/app/bin/persistenceCore | grep "statically linked")

FROM alpine:3.16

COPY --from=go-builder /usr/local/app/bin/persistenceCore /usr/bin/persistenceCore

COPY contrib/local/ /opt/
RUN chmod +x /opt/*.sh

# Set up dependencies
ENV PACKAGES curl make bash jq
ENV CHAIN_BIN /usr/bin/persistenceCore

# Install minimum necessary dependencies
RUN apk add --no-cache $PACKAGES

WORKDIR /opt

# rest server, p2p, rpc, grpc
EXPOSE 1317 26656 26657 9090

# Run persistenceCore by default, omit entrypoint to ease using container with cli
CMD ["/usr/bin/persistenceCore", "version", "--long"]
