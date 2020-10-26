.PHONY: test
test: lint
	go test -coverprofile=coverage.txt -covermode=atomic -race ./...

.PHONY: lint
lint: check-format vet

.PHONY: lint-local
lint-local: lint
	golangci-lint run

.PHONY: vet
vet:
	go vet ./...

.PHONY: check-format
check-format:
	@echo "Running gofmt..."
	$(eval unformatted=$(shell find . -name '*.go' | grep -v ./.git | grep -v vendor | xargs gofmt -l))
	$(if $(strip $(unformatted)),\
		$(error $(\n) Some files are not formatted properly! Run: \
			$(foreach file,$(unformatted),$(\n)    gofmt -w $(file))$(\n)),\
		@echo All files are well formatted.\
	)

.PHONY: install-ci
install-ci:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(GOPATH)/bin v1.30.0
