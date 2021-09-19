FROM golang:1.16.6-alpine3.14 AS builder

RUN apk update && \
  apk add --no-cache --update make bash git ca-certificates && \
  update-ca-certificates

WORKDIR /go/src/bankroach

COPY . .

RUN make build

FROM alpine:3.14.0

COPY --from=builder /go/src/bankroach/bin/bankroach /bankroach

CMD [ "/bankroach" ]