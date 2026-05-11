ARG VARIANT=1.26-trixie
FROM golang:${VARIANT}

RUN apt-get update && \
    export DEBIAN_FRONTEND=noninteractive && \
    apt-get -y install --no-install-recommends fonts-firacode zstd && \
    rm -rf /var/lib/apt/lists/*
