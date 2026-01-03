REPO=github.com/sean9999/polity/v3
SEMVER := $$(git tag --sort=-version:refname | head -n 1)
BRANCH := $$(git branch --show-current)
REF := $$(git describe --dirty --tags --always)

info:
	@printf "REPO:\t%s\nSEMVER:\t%s\nBRANCH:\t%s\nREF:\t%s\n" $(REPO) $(SEMVER) $(BRANCH) $(REF)

container:
	docker build -t polityd -t polityd:$(BRANCH) -t polityd:$(REF) -t polityd:$(SEMVER) .

bin:
	mkdir -p bin

binaries: bin/polityd bin/polity

bin/polity: bin/polity/polity-linux-amd64 bin/polity/polity-linux-arm64 bin/polity/polity-darwin-amd64 bin/polity/polity-darwin-arm64 bin/polity/polity-windows-amd64

bin/polity/polity-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o bin/polity-linux-amd64 -ldflags="-X 'main.Version=$(REF)' -X 'main.Branch=$(BRANCH)'" ./cmd/polity;

bin/polity/polity-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -o bin/polity-linux-arm64 -ldflags="-X 'main.Version=$(REF)' -X 'main.Branch=$(BRANCH)'" ./cmd/polity;

bin/polity/polity-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -o bin/polity-darwin-amd64 -ldflags="-X 'main.Version=$(REF)' -X 'main.Branch=$(BRANCH)'" ./cmd/polity;

bin/polity/polity-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -o bin/polity-darwin-arm64 -ldflags="-X 'main.Version=$(REF)' -X 'main.Branch=$(BRANCH)'" ./cmd/polity;

bin/polity/polity-windows-amd64:
	GOOS=windows GOARCH=amd64 go build -o bin/polity-windows-amd64 -ldflags="-X 'main.Version=$(REF)' -X 'main.Branch=$(BRANCH)'" ./cmd/polity;


bin/polityd: bin/polityd/polityd-linux-amd64 bin/polityd/polityd-linux-arm64 bin/bin/polityd-darwin-amd64 bin/polityd-darwin-arm64 bin/polityd-windows-arm64.exe bin/polityd-windows-386.exe

bin/polityd/polityd-linux-amd64: bin
	GOOS=linux   GOARCH=amd64 go build -o bin/polityd-linux-amd64 -ldflags="-X 'main.Version=$(REF)' -X 'main.Branch=$(BRANCH)'" ./cmd/polityd

bin/polityd/polityd-linux-arm64: bin
	GOOS=linux   GOARCH=arm64 go build -o bin/polityd-linux-arm64 -ldflags="-X 'main.Version=$(REF)' -X 'main.Branch=$(BRANCH)'" ./cmd/polityd

bin/bin/polityd-darwin-amd64: bin
	GOOS=darwin  GOARCH=amd64 go build -o bin/polityd-darwin-amd64 -ldflags="-X 'main.Version=$(REF)' -X 'main.Branch=$(BRANCH)'" ./cmd/polityd

bin/polityd-darwin-arm64: bin
	GOOS=darwin  GOARCH=arm64 go build -o bin/polityd-darwin-arm64 -ldflags="-X 'main.Version=$(REF)' -X 'main.Branch=$(BRANCH)'" ./cmd/polityd

bin/polityd-windows-arm64.exe: bin
	GOOS=windows GOARCH=amd64 go build -o bin/polityd-windows-arm64.exe -ldflags="-X 'main.Version=$(REF)' -X 'main.Branch=$(BRANCH)'" ./cmd/polityd;

bin/polityd-windows-386.exe: bin
	GOOS=windows GOARCH=386 go build -o bin/polityd-windows-386.exe -ldflags="-X 'main.Version=$(REF)' -X 'main.Branch=$(BRANCH)'" ./cmd/polityd;

tidy:
	go mod tidy

install:
	go install ./cmd/polityd
	mkdir -p ${HOME}/.config/polity
	touch ${HOME}/.config/polity/config.json

clean:
	go clean
	rm bin/*

upgrade:
	go get -u
	go mod tidy

pkgsite:
	@if [ -z "$$(go tool -n golang.org/x/pkgsite/cmd/pkgsite)" ]; then go get -tool golang.org/x/pkgsite/cmd/pkgsite@latest; fi

docs: pkgsite
	go tool golang.org/x/pkgsite/cmd/pkgsite -open .

publish:
	GOPROXY=https://proxy.golang.org,direct go list -m ${REPO}@${SEMVER}

test:
	go test -coverprofile=coverage.out ./...

coverage:
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated:"
	@echo "  - Text summary: go tool cover -func=coverage.out"
	@echo "  - HTML report: open coverage.html"
	@COVERAGE=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	echo "Total coverage: $$COVERAGE%"; \
	if [ "$$(echo "$$COVERAGE < 85" | bc -l)" -eq 1 ]; then \
		echo "❌ Coverage below 85% threshold"; \
		exit 1; \
	else \
		echo "✅ Coverage meets 85% threshold"; \
	fi
