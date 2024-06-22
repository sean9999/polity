REPO=github.com/sean9999/polity
SEMVER := $$(git tag --sort=-version:refname | head -n 1)
BRANCH := $$(git branch --show-current)
REF := $$(git describe --dirty --tags --always)

info:
	@printf "REPO:\t%s\nSEMVER:\t%s\nBRANCH:\t%s\nREF:\t%s\n" $(REPO) $(SEMVER) $(BRANCH) $(REF)

binaries: bin/polityd bin/polity
	mkdir -p bin

bin/polity:
	go build -v -o bin/polity -ldflags="-X 'main.Version=$(REF)'" cmd/polity/**.go
	
bin/polityd:	
	go build -v -o bin/polityd -ldflags="-X 'main.Version=$(REF)'" cmd/polityd/**.go

tidy:
	go mod tidy

install:
	go install ./cmd/polity
	go install ./cmd/polityd
	mkdir -p ${HOME}/.config/polity
	touch ${HOME}/.config/polity/conf.toml

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
