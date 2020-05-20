FROM golang:alpine

COPY . /app

WORKDIR /app

RUN go build -o interview cmd/haha/main.go

CMD ["/app/interview"]