ARG ALPINE_VERSION=3.12
FROM alpine:${ALPINE_VERSION}
WORKDIR /
RUN apk --no-cache add ca-certificates
COPY smartdevicemanagement_exporter /smartdevicemanagement_exporter
ENTRYPOINT ./smartdevicemanagement_exporter
