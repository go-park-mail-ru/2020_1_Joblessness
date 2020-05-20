FROM golang:alpine

COPY . /app

WORKDIR /app

RUN go build -o search cmd/search/main.go

CMD ["/app/search"]