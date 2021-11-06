FROM golang:1.17 AS build

WORKDIR /go/src/github.com/lukasCoppens/openhab-item-exporter
COPY . .
RUN go build -v ./...

FROM alpine:latest AS final
COPY --from=0 /go/src/github.com/lukasCoppens/openhab-item-exporter/openhab-item-exporter /bin/openhab-item-exporter
ENTRYPOINT ["/bin/openhab-item-exporter"]