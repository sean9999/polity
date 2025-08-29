REPO=github.com/sean9999/polity/v2
SEMVER := $$(git tag --sort=-version:refname | head -n 1)
BRANCH := $$(git branch --show-current)
REF := $$(git describe --dirty --tags --always)

info:
	@printf "REPO:\t%s\nSEMVER:\t%s\nBRANCH:\t%s\nREF:\t%s\n" $(REPO) $(SEMVER) $(BRANCH) $(REF)

binaries: bin/polityd bin/polity bin/polityctl

bin:
	mkdir -p bin

bin/polity: bin
	go build -v -o bin/polity -ldflags="-X 'main.Version=$(REF)'" ./cmd/polity/
	env GOOS=linux GOARCH=arm go build -v -o bin/polity_linux_arm -ldflags="-X 'main.Version=$(REF)'" ./cmd/polity/
	env GOOS=linux GOARCH=amd64 go build -v -o bin/polity_linux_amd64 -ldflags="-X 'main.Version=$(REF)'" ./cmd/polity/
	env GOOS=darwin GOARCH=arm go build -v -o bin/polity_darwin_arm -ldflags="-X 'main.Version=$(REF)'" ./cmd/polity/
	env GOOS=freebsd GOARCH=amd64 go build -v -o bin/polity_freebsd_amd64 -ldflags="-X 'main.Version=$(REF)'" ./cmd/polity/

bin/polityd: bin
	go build -v -o bin/polityd -ldflags="-X 'main.Version=$(REF)'" ./cmd/polityd/
	env GOOS=linux GOARCH=arm go build -o bin/polityd_linux_arm -ldflags="-X 'main.Version=$(REF)'" ./cmd/polityd/
	env GOOS=linux GOARCH=amd64 go build -o bin/polityd_linux_amd64 -ldflags="-X 'main.Version=$(REF)'" ./cmd/polityd/
	env GOOS=darwin GOARCH=arm64 go build -o bin/polityd_darwin_arm64 -ldflags="-X 'main.Version=$(REF)'" ./cmd/polityd/
	env GOOS=freebsd GOARCH=amd64 go build -o bin/polityd_freebsd_amd64 -ldflags="-X 'main.Version=$(REF)'" ./cmd/polityd/

bin/polityctl: bin
	go build -v -o bin/polityctl -ldflags="-X 'main.Version=$(REF)'" ./cmd/polityctl/
	env GOOS=linux GOARCH=arm go build -o bin/polityd_linux_arm -ldflags="-X 'main.Version=$(REF)'" ./cmd/polityctl/
	env GOOS=linux GOARCH=amd64 go build -o bin/polityd_linux_amd64 -ldflags="-X 'main.Version=$(REF)'" ./cmd/polityctl/
	env GOOS=darwin GOARCH=arm64 go build -o bin/polityd_darwin_arm64 -ldflags="-X 'main.Version=$(REF)'" ./cmd/polityctl/
	env GOOS=freebsd GOARCH=amd64 go build -o bin/polityd_freebsd_amd64 -ldflags="-X 'main.Version=$(REF)'" ./cmd/polityctl/

tidy:
	go mod tidy

install:
	go install ./cmd/polityd
	mkdir -p ${HOME}/.config/polity

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

.PHONY: test
