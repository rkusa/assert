deps:
	go get -t -v ./...

test:
	go vet ./...
	go test -cover -short ./...

.PHONY: deps, test