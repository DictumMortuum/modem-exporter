PREFIX=/usr/local

format:
	gofmt -s -w .

build: format
	go build -trimpath -buildmode=pie -mod=readonly -modcacherw -ldflags="-s -w"

install: build
	mkdir -p $(PREFIX)/bin
	cp -f modem-exporter $(PREFIX)/bin

install-service:
	cp -f assets/prometheus-modem-exporter.service /etc/systemd/system/
