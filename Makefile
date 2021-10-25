format:
	gofmt -s -w .

build: format
	go build
