.PHONY: tidy
tidy:
	go mod tidy -v
	# go fmt ./...

.PHONY: build
build:
	go build -o=/tmp/bin/traintrack ./cmd/api
	
.PHONY: run
run: build
	/tmp/bin/traintrack
