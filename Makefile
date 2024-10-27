run: build
	./tmp/bin/url-shortener
	

build:
	go fmt ./...
	go build -o=./tmp/bin/url-shortener cmd/url-shortener/main.go
