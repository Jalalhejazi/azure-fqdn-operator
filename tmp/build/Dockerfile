FROM alpine:3.6

RUN apk update && apk add ca-certificates

RUN adduser -D azure-fqdn-operator
USER azure-fqdn-operator

ADD tmp/_output/bin/azure-fqdn-operator /usr/local/bin/azure-fqdn-operator
