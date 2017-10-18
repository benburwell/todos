build: vendor
	go build .

install: vendor
	go install .

clean:
	rm -f todos

vendor:
	dep ensure

.PHONY: build install clean
