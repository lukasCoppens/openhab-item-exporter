FROM alpine:latest
RUN apk update && apk add ca-certificates
ADD openhab-item-exporter /bin/openhab-item-exporter
RUN chmod u+x /bin/openhab-item-exporter
ENTRYPOINT [ "/bin/openhab-item-exporter" ]