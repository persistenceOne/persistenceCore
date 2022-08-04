FROM informalsystems/hermes:1.0.0-rc.1

# Set up dependencies
ENV PACKAGES curl make git jq ca-certificates

USER root

# Install ca-certificates
RUN set -eux;

# Install minimum necessary dependencies
RUN apt-get update; apt-get install -y $PACKAGES
