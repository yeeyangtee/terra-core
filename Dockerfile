# docker build . -t cosmwasm/wasmd:latest
# docker run --rm -it cosmwasm/wasmd:latest /bin/sh
FROM golang:1.15-alpine3.12 AS go-builder

# this comes from standard alpine nightly file
#  https://github.com/rust-lang/docker-rust-nightly/blob/master/alpine3.12/Dockerfile
# with some changes to support our toolchain, etc
RUN set -eux; apk add --no-cache ca-certificates build-base;

RUN apk add git
# NOTE: add these to run with LEDGER_ENABLED=true
# RUN apk add libusb-dev linux-headers

WORKDIR /code
COPY . /code/

# See https://github.com/terra-project/go-cosmwasm/releases
ADD https://github.com/terra-project/go-cosmwasm/releases/download/v0.10.5/libgo_cosmwasm_muslc.a /lib/libgo_cosmwasm_muslc.a
RUN sha256sum /lib/libgo_cosmwasm_muslc.a | grep 1f6d76ada5652553bb4a0f6454289b769267e575395c220b65e90d670aa35557

# force it to use static lib (from above) not standard libgo_cosmwasm.so file
RUN LEDGER_ENABLED=false BUILD_TAGS=muslc make update-swagger-docs build

FROM alpine:3.12

WORKDIR /root

COPY --from=go-builder /code/build/terrad /usr/local/bin/terrad
COPY --from=go-builder /code/build/terracli /usr/local/bin/terracli

CMD ["/usr/local/bin/terrad", "version"]
