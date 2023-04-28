.PHONY: all build osqueryd osqueryi vendor tidy clean
GOCMD=go
GOBUILD=$(GOCMD) build
GOMOD=$(GOCMD) mod

all: build

build:
	echo "$(shell pwd)/build/s3-config-extension.ext" > /tmp/extensions.load
	$(GOBUILD) -o build/s3-config-extension.ext .

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

vendor:
	$(GOMOD) vendor

tidy:
	$(GOMOD) tidy

clean:
	rm -rf /tmp/extensions.load
	rm -rf /tmp/osquery.*
	rm -rf build
