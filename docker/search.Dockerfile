FROM golang:alpine

COPY . /app

WORKDIR /app

RUN apk add make && make build-search

CMD ["/app/search"]