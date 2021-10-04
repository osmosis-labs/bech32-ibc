FROM golang:1.16-alpine AS build-env
ARG VERSION

ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev

RUN apk add --no-cache $PACKAGES

WORKDIR /go/src/github.com/osmosis-labs

RUN git clone https://github.com/osmosis-labs/bech32-ibc.git

WORKDIR /go/src/github.com/osmosis-labs/bech32-ibc

RUN git checkout ${VERSION} && make build-linux

FROM alpine:edge

RUN apk add --no-cache ca-certificates
WORKDIR /root

COPY --from=build-env /go/src/github.com/osmosis-labs/bech32-ibc/build/bech32ibcd /usr/bin/bech32ibcd

#USER bech32ibc

WORKDIR /bech32ibc

COPY ./test/setup/bech32ibc-setup.sh .

COPY ./test/setup/valkeys ./setup/valkeys

USER root

RUN chmod -R 777 ./setup

#USER bech32ibc

EXPOSE 26657

ENTRYPOINT [ "./bech32ibc-setup.sh" ]
# NOTE: to run this image, docker run -d -p 26657:26657 ./bech32ibc-setup.sh {{chain_id}} {{genesis_account}} {{seeds}} {{priv_validator_key_path}}