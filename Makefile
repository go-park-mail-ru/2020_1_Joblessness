build-app:
	go build -o app cmd/haha/main.go

build-search:
	go build -o search cmd/search/main.go

build-interview:
	go build -o interview cmd/interview/main.go

build-auth:
	go build -o auth cmd/auth/main.go

test:
	go test ./... -coverprofile cover.out.tmp | grep -v "no test files"

cover:
	cat cover.out.tmp | grep -v "_easyjson.go" > cover.out
	go tool cover -func cover.out | grep "total"
	rm -f cover.out*