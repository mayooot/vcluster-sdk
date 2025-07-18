FROM alpine:3.16

RUN apk add --no-cache ca-certificates

COPY ./stuff /stuff

ENTRYPOINT ["/stuff"]