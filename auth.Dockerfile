FROM golang:alpine

COPY . /app

WORKDIR /app

RUN go build -o auth cmd/auth/main.go

CMD ["/app/auth"]