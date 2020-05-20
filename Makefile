build-app:
	go build -o app cmd/haha/main.go

build-search:
	go build -o search cmd/api/main.go

build-interview:
	go build -o interview cmd/api/main.go

build-auth:
	go build -o auth cmd/api/main.go