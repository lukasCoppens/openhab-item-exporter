FROM golang:1.17 AS build

WORKDIR /go/src/github.com/lukasCoppens/openhab-item-exporter
COPY . .
RUN go mod vendor && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v .
RUN ls -l && pwd

FROM alpine:latest AS final
COPY --from=0 /go/src/github.com/lukasCoppens/openhab-item-exporter/openhab-item-exporter /lcoppens/
RUN chmod u+x /lcoppens/openhab-item-exporter
ENTRYPOINT ["/lcoppens/openhab-item-exporter"]