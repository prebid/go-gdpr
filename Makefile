.PHONY: info fmt goimports gofumpt lint go_fix go_vet golangci test coverage clean

info:
	go version

fmt: goimports gofumpt
	$(info === format done)

goimports:
	goimports -e -l -w -local github.com/prebid/go-gdpr .

gofumpt:
	gofumpt -l -w -extra .

lint: go.sum go_fix go_vet golangci
	$(info === lint done)

go.mod:
	go mod tidy
	go mod verify

go.sum: go.mod

go_fix:
	go fix ./...

go_vet:
	go vet -all ./...

golangci:
	golangci-lint run ./...

test:
	go test -v ./...

coverage:
	export GOEXPERIMENT="nocoverageredesign"
	go test -v -race -cover -covermode=atomic -coverprofile coverage.out ./...

clean:
	rm -f $(BINARY)
	rm -f coverage.*
	rm -f .test_report.xml
