BINARY_NAME=nullbeacon
TEAMSERVER=./cmd/teamserver
AGENT=./cmd/agent

.PHONY: build build-agent clean fmt lint release

build:
	go build -ldflags "-s -w" -o $(BINARY_NAME) $(TEAMSERVER)

build-agent:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o $(BINARY_NAME)-linux-amd64 $(AGENT)
	GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o $(BINARY_NAME)-linux-arm64 $(AGENT)

clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME)-linux-* teamserver.log

fmt:
	gofmt -w .

lint:
	golangci-lint run ./...

release:
	@read -p "Version tag (e.g. v0.1.4): " tag; \
	git tag $$tag && git push origin $$tag
