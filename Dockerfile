ARG BASE_VERSION=1.24-alpine3.23
FROM golang:${BASE_VERSION} AS init-env

RUN apk add --update --no-cache curl make git libc-dev bash gcc linux-headers eudev-dev ncurses-dev

ARG TARGETARCH
ARG BUILDARCH
ARG GITHUB_ORGANIZATION
ARG REPO_HOST
ARG GITHUB_REPO
ARG WASMVM_VERSION

WORKDIR /go/src/${REPO_HOST}/${GITHUB_ORGANIZATION}/${GITHUB_REPO}

# Download CosmWasm libwasmvm if found
RUN set -eux; \
    export ARCH=$(uname -m); \
    if [ ! -z "${WASMVM_VERSION}" ]; then\
      WASMVM_REPO=$(echo $WASMVM_VERSION | awk '{print $1}');\
      WASMVM_VERS=$(echo $WASMVM_VERSION | awk '{print $2}');\
      wget -O /lib/libwasmvm_muslc.a https://${WASMVM_REPO}/releases/download/${WASMVM_VERS}/libwasmvm_muslc.$(uname -m).a;\
      ln /lib/libwasmvm_muslc.a /lib/libwasmvm_muslc.$(uname -m).a;\
    fi;

ARG BUILD_DIR

ADD ${BUILD_DIR}/go.mod ${BUILD_DIR}/go.sum ./

ARG CLONE_KEY

RUN if [ ! -z "${CLONE_KEY}" ]; then\
  mkdir -p ~/.ssh;\
  echo "${CLONE_KEY}" | base64 -d > ~/.ssh/id_ed25519;\
  chmod 600 ~/.ssh/id_ed25519;\
  apk add openssh;\
  git config --global --add url."ssh://git@github.com/".insteadOf "https://github.com/";\
  ssh-keyscan github.com >> ~/.ssh/known_hosts;\
  fi

ARG VENDOR

# Download go mod dependencies
RUN set -eux; \
    if [[ "${BUILD_DIR}" == "." && "${VENDOR}" == "false" ]]; then\
        go mod download;\
    fi;

# Use minimal busybox from infra-toolkit image
FROM ghcr.io/strangelove-ventures/infra-toolkit:v0.1.12 AS infra-toolkit
RUN addgroup --gid 1025 -S heighliner && adduser --uid 1025 -S heighliner -G heighliner

# Use alpine to source CA certificates
FROM alpine:3 as alpine-3

# Build binary from local source
FROM init-env AS build-env

ARG BUILD_TARGET
ARG BUILD_ENV
ARG BUILD_TAGS
ARG PRE_BUILD
ARG BUILD_DIR

ADD . .

RUN set -eux; \
    export CGO_ENABLED=1 LDFLAGS='-linkmode external -extldflags "-static"'; \
    if [ ! -z "$PRE_BUILD" ]; then sh -c "${PRE_BUILD}"; fi; \
    if [ ! -z "$BUILD_TARGET" ]; then \
      if [ ! -z "$BUILD_ENV" ]; then export ${BUILD_ENV}; fi; \
      if [ ! -z "$BUILD_TAGS" ]; then export "${BUILD_TAGS}"; fi; \
      if [ ! -z "$BUILD_DIR" ]; then cd "${BUILD_DIR}"; fi; \
      sh -c "${BUILD_TARGET}"; \
    fi

# Collect binaries
RUN mkdir /root/bin
ARG RACE
ARG BINARIES
ENV BINARIES_ENV ${BINARIES}

RUN bash -c 'set -eux;\
  BINARIES_ARR=();\
  IFS=, read -ra BINARIES_ARR <<< "$BINARIES_ENV";\
  for BINARY in "${BINARIES_ARR[@]}"; do\
    BINSPLIT=();\
    IFS=: read -ra BINSPLIT <<< "$BINARY";\
    BINPATH=${BINSPLIT[1]+"${BINSPLIT[1]}"};\
    BIN="$(eval "echo "${BINSPLIT[0]+"${BINSPLIT[0]}"}"")";\
    if [ ! -z "$RACE" ] && GOVERSIONOUT=$(go version -m $BIN); then\
      if echo $GOVERSIONOUT | grep build | grep "-race=true"; then\
        echo "Race detection is enabled in binary";\
      else\
        echo "Race detection not enabled in binary!";\
        exit 1;\
      fi;\
    fi;\
    if [ ! -z "$BINPATH" ]; then\
      if [[ $BINPATH == *"/"* ]]; then\
        mkdir -p "$(dirname "${BINPATH}")";\
        cp "$BIN" "${BINPATH}";\
      else\
        cp "$BIN" "/root/bin/${BINPATH}";\
      fi;\
    else\
      cp "$BIN" /root/bin/;\
    fi;\
  done'

# Libraries
RUN mkdir -p /root/lib
ARG LIBRARIES
ENV LIBRARIES_ENV ${LIBRARIES}

RUN bash -c 'set -eux;\
  LIBRARIES_ARR=($LIBRARIES_ENV); \
  for LIBRARY in "${LIBRARIES_ARR[@]}"; do cp $LIBRARY /root/lib/; done'

# Directories
RUN mkdir -p /root/dir_abs && touch /root/dir_abs.list
ARG DIRECTORIES
ENV DIRECTORIES_ENV ${DIRECTORIES}

RUN bash -c 'set -eux;\
  DIRECTORIES_ARR=($DIRECTORIES_ENV);\
  i=0;\
  for DIRECTORY in "${DIRECTORIES_ARR[@]}"; do \
    cp -R $DIRECTORY /root/dir_abs/$i;\
    echo $DIRECTORY >> /root/dir_abs.list;\
    ((i = i + 1));\
  done'

# Final image
FROM scratch

LABEL org.opencontainers.image.source="https://github.com/strangelove-ventures/heighliner"

WORKDIR /bin

COPY --from=infra-toolkit /busybox/busybox /bin/sh
COPY --from=infra-toolkit /busybox/busybox /bin/ln
COPY --from=infra-toolkit /usr/local/bin/jq /bin/

RUN for b in \
  cat date df dirname du env grep head less ls md5sum mkdir mv pwd rm sed \
  sha1sum sha256sum sha3sum sha512sum sleep stty tail tar tee tr vi watch which \
  ; do ln ln $b; done; \
  rm -rf sh; \
  ln ln sh;

COPY --from=build-env /root/dir_abs /root/dir_abs
COPY --from=build-env /root/dir_abs.list /root/dir_abs.list

RUN sh -c 'i=0; while read DIR; do\
      PLACEDIR="$(dirname "$DIR")";\
      mkdir -p "$PLACEDIR";\
      mv /root/dir_abs/$i $DIR;\
      i=$((i+1));\
    done < /root/dir_abs.list'

COPY --from=alpine-3 /etc/ssl/cert.pem /etc/ssl/cert.pem

COPY --from=infra-toolkit /etc/passwd /etc/passwd
COPY --from=infra-toolkit --chown=1025:1025 /home/heighliner /home/heighliner
COPY --from=infra-toolkit --chown=1025:1025 /tmp /tmp

COPY --from=build-env /root/bin /bin
COPY --from=build-env /root/lib /lib

WORKDIR /home/heighliner
USER heighliner