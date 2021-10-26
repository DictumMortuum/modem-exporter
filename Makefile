PREFIX=/usr/local

format:
	gofmt -s -w .

build: format
	go build

install: build
	mkdir -p $(PREFIX)/bin
	cp -f servus $(PREFIX)/bin

install-service:
	cp -f systemd/servus.service /etc/systemd/system/
