.PHONY: build lint clean test

build: lint test contextdir

contextdir: *.go */*.go
	@echo '==> Building $@'
	go build -o "$@"

lint:
	@echo '==> Linting'
	! gofmt -d -e . | grep .
	go vet
	staticcheck

test:
	@echo '==> Testing'
	go test -v

clean:
	@echo '==> Cleaning'
	rm -rf -- contextdir test
