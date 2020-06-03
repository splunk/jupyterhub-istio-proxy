.PHONY: test
test: lint
	go test -coverprofile=coverage.txt -covermode=atomic -race ./...

.PHONY: lint
lint: vet check-format
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

.PHONY: publish-coverage
publish-coverage:
	# publishes results from coverage.txt
	curl -s https://codecov.io/bash | bash
