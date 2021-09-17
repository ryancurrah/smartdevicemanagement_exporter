ARG ALPINE_VERSION=3.12

FROM golang:alpine${ALPINE_VERSION} as builder
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN GO111MODULE=on  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/smartdevicemanagement_exporter /app/cmd/smartdevicemanagement_exporter/main.go
#RUN GO111MODULE=on  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/app/main /go/app/cmd/myapp/main.go

FROM alpine:${ALPINE_VERSION}
WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/smartdevicemanagement_exporter /usr/local/bin/smartdevicemanagement_exporter
ENTRYPOINT ["smartdevicemanagement_exporter"]
