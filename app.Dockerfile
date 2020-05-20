FROM golang:alpine

COPY . /app

WORKDIR /app

RUN apk add make && make build-app

COPY /etc/letsencrypt /etc/letsencrypt

CMD ["/app/app"]