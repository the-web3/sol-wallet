GITCOMMIT := $(shell git rev-parse HEAD)
GITDATE := $(shell git show -s --format='%ct')

LDFLAGSSTRING +=-X main.GitCommit=$(GITCOMMIT)
LDFLAGSSTRING +=-X main.GitDate=$(GITDATE)
LDFLAGS := -ldflags "$(LDFLAGSSTRING)"

sol-wallet:
	env GO111MODULE=on go build -v $(LDFLAGS) ./cmd/sol-wallet

clean:
	rm sol-wallet

test:
	go test -v ./...

lint:
	golangci-lint run ./...

.PHONY: \
	sol-wallet \
	bindings \
	bindings-scc \
	clean \
	test \
	lint