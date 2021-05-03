FROM golang:1.16-alpine

# Set up dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python3

# Set working directory for the build
WORKDIR /usr/local/app

# Add source files
COPY . .

# Install minimum necessary dependencies, build persistenceCore, remove packages
RUN apk add --no-cache $PACKAGES && make install

# Install ca-certificates
RUN apk add --update ca-certificates

# Run persistenceCore by default, omit entrypoint to ease using container with cli
CMD ["persistenceCore"]
