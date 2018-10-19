BUILD_DIR := build

build: build-ccacheparser
.PHONY: build

build-%:
	go build -o $(BUILD_DIR)/$* ./cmd/$*

distclean:
	rm -rf build

test:
	go test
