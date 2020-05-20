FROM golang:alpine

COPY . /app

WORKDIR /app

RUN go build -o app cmd/haha/main.go

CMD ["/app/app"]