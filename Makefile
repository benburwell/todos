SRC_DIR=./src
BUILD_DIR=./build
SOURCES=$(wildcard *.go src/*.go src/*/*.go)
BINARY=todos

$(BUILD_DIR)/$(BINARY): $(SOURCES) vendor
	mkdir -p $(BUILD_DIR)
	go build -o $@ $(SRC_DIR)

install: $(SOURCES) vendor
	go install $(SRC_DIR)

clean:
	rm -rf $(BUILD_DIR)

vendor:
	dep ensure

.PHONY: install clean
