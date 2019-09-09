ARG GOLANG_VERSION=1.13.0
ARG ALPINE_VERSION=3.10

FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS build

RUN apk --no-cache add ca-certificates git

COPY . /github.com/lightstep/reverse-proxy
WORKDIR /github.com/lightstep/reverse-proxy

RUN CGO_ENABLED=0 go build -o reverse-proxy

FROM scratch
EXPOSE 8126
COPY --from=build /github.com/lightstep/reverse-proxy/reverse-proxy reverse-proxy
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT [ "./reverse-proxy" ]
