.PHONY: build env osqueryd osqueryi deps

all: build

deps:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure -vendor-only

build:
	echo "$(shell pwd)/build/s3-config-extension.ext" > /tmp/extensions.load
	go build -i -o build/s3-config-extension.ext .

osqueryd: build
	osqueryd \
		--extensions_autoload=/tmp/extensions.load \
		--pidfile=/tmp/osquery.pid \
		--database_path=/tmp/osquery.db \
		--extensions_socket=/tmp/osquery.sock \
		--config_refresh=60 \
		--config_plugin=s3

osqueryi: build
	osqueryi --extension=./build/s3-config-extension.ext

clean:
	rm -rf /tmp/extensions.load
	rm -rf /tmp/osquery.*
	rm -rf build
