.PHONY: test
test:
	go test -v ./...

.PHONY: release
release:
	goreleaser release --snapshot --clean

.PHONY: release-check
release-check:
	goreleaser check
