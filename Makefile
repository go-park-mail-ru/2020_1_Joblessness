build-app:
	go build -o app cmd/haha/main.go

build-search:
	go build -o search cmd/search/main.go

build-interview:
	go build -o interview cmd/interview/main.go

build-auth:
	go build -o auth cmd/auth/main.go