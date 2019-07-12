# Usage :
# docker build -t ara .
# docker run -it --add-host db:172.17.0.1 -e EDWIG_DB_NAME=edwig -e EDWIG_DB_USER=edwig -e EDWIG_DB_PASSWORD=edwig -e EDWIG_API_KEY=secret -e EDWIG_DEBUG=true -p 8080:8080 ara

FROM golang:1.12 AS builder

ENV DEV_PACKAGES="libxml2-dev"
RUN apt-get update && apt-get -y install --no-install-recommends $DEV_PACKAGES

WORKDIR /go/src/github.com/af83/edwig
COPY . .

ENV GO111MODULE=on
RUN go get -d -v ./...
RUN go install -v ./...

FROM debian:latest

ENV RUN_PACKAGES="libxml2"

RUN apt-get update && apt-get -y dist-upgrade && apt-get -y install --no-install-recommends $RUN_PACKAGES && \
    apt-get clean && apt-get -y autoremove && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /go/bin/edwig ./
ADD docker-entrypoint.sh ./
COPY db/migrations ./db/migrations
RUN chmod +x ./edwig ./docker-entrypoint.sh && mkdir ./config

ENV EDWIG_CONFIG=./config EDWIG_ENV=production
EXPOSE 8080

ENTRYPOINT ["./docker-entrypoint.sh"]
CMD ["api"]
