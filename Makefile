REPO=github.com/sean9999/polity/v3
SEMVER := $$(git tag --sort=-version:refname | head -n 1)
BRANCH := $$(git branch --show-current)
REF := $$(git describe --dirty --tags --always)

info:
	@printf "REPO:\t%s\nSEMVER:\t%s\nBRANCH:\t%s\nREF:\t%s\n" $(REPO) $(SEMVER) $(BRANCH) $(REF)

binaries: bin/polityd bin/polity
	mkdir -p bin

bin/polity:
	go build -o bin/polity -ldflags="-X 'main.Version=$(REF)' -X 'main.Branch=$(BRANCH)'" v3/cmd/polity
	GOOS=linux  GOARCH=amd64 go build -o bin/polity-linux-amd64  v3/cmd/polity;
	GOOS=linux  GOARCH=arm64 go build -o bin/polity-linux-arm64  v3/cmd/polity;
	GOOS=darwin GOARCH=amd64 go build -o bin/polity-darwin-amd64 v3/cmd/polity;
	GOOS=darwin GOARCH=arm64 go build -o bin/polity-darwin-arm64 v3/cmd/polity;

bin/polityd:	
	go build -o bin/polityd -ldflags="-X 'main.Version=$(REF)' -X 'main.Branch=$(BRANCH)'" v3/cmd/polityd
	GOOS=linux  GOARCH=amd64 go build -o bin/polityd-linux-amd64  v3/cmd/polityd;
	GOOS=linux  GOARCH=arm64 go build -o bin/polityd-linux-arm64  v3/cmd/polityd;
	GOOS=darwin GOARCH=amd64 go build -o bin/polityd-darwin-amd64 v3/cmd/polityd;
	GOOS=darwin GOARCH=arm64 go build -o bin/polityd-darwin-arm64 v3/cmd/polityd;

tidy:
	go mod tidy

install:
	cd ./v2 && go install ./cmd/polityd
	mkdir -p ${HOME}/.config/polity
	touch ${HOME}/.config/polity/config.json

clean:
	go clean
	go clean -modcache
	rm bin/*

pkgsite:
	if [ -z "$$(command -v pkgsite)" ]; then go install golang.org/x/pkgsite/cmd/pkgsite@latest; fi

docs: pkgsite
	pkgsite -open .

publish:
	GOPROXY=https://proxy.golang.org,direct go list -m ${REPO}@${SEMVER}

test:
	go test ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated:"
	@echo "  - Text summary: go tool cover -func=coverage.out"
	@echo "  - HTML report: open coverage.html"

coverage-summary:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

coverage-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "HTML coverage report generated: coverage.html"

coverage-check:
	@go test -coverprofile=coverage.out ./... > /dev/null
	@COVERAGE=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	echo "Total coverage: $$COVERAGE%"; \
	if [ "$$(echo "$$COVERAGE < 70" | bc -l)" -eq 1 ]; then \
		echo "❌ Coverage below 70% threshold"; \
		exit 1; \
	else \
		echo "✅ Coverage meets 70% threshold"; \
	fi

.PHONY: test coverage coverage-summary coverage-html coverage-check
