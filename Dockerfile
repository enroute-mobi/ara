# Usage :
# docker build -t ara --build-arg VERSION="build-999" .
# docker run -it --add-host db:172.17.0.1 -e ARA_DB_NAME=ara -e ARA_DB_USER=ara -e ARA_DB_PASSWORD=ara -e ARA_API_KEY=secret -e ARA_DEBUG=true -p 8080:8080 ara

FROM golang:1.16 AS builder
ARG VERSION=dev

ENV DEV_PACKAGES="libxml2-dev"
RUN apt-get update && apt-get -y install --no-install-recommends $DEV_PACKAGES

WORKDIR /go/src/bitbucket.org/enroute-mobi/ara
COPY . .

RUN go install -v -ldflags "-X bitbucket.org/enroute-mobi/ara/version.value=${VERSION}" ./...

FROM debian:buster

ENV RUN_PACKAGES="libxml2 ca-certificates"

RUN apt-get update && apt-get -y dist-upgrade && apt-get -y install --no-install-recommends $RUN_PACKAGES && \
    apt-get clean && apt-get -y autoremove && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /go/bin/ara ./
COPY docker-entrypoint.sh ./
COPY db/migrations ./db/migrations
COPY siri/templates ./siri/templates
RUN chmod +x ./ara ./docker-entrypoint.sh && mkdir ./config

ENV ARA_CONFIG=./config ARA_ENV=production ARA_ROOT=/app
EXPOSE 8080

ENTRYPOINT ["./docker-entrypoint.sh"]
CMD ["api"]
