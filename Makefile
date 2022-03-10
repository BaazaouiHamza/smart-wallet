PROJECT_NAME := "smart-wallet"
PKG := "git.digitus.me/pfe/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY: test coverage coverhtml lint

lint:
	@go get -u github.com/alecthomas/gometalinter
	@gometalinter --install
	@gometalinter ./... --deadline=1m

swagger:
	@swag i --parseDependency --parseDepth 1 -g api/server.go
	@rg --passthru -U '"additionalProperties":\s*\{\s*"type"\s*:\s*"integer"\s*\}' -r '"additionalProperties": {"type": "string", "enum": ["Contributor", "User"]}' docs/docs.go | sponge docs/docs.go
	@rg --passthru -U '"additionalProperties":\s*\{\s*"type"\s*:\s*"integer"\s*\}' -r '"additionalProperties": {"type": "string", "enum": ["Contributor", "User"]}' docs/swagger.json | sponge docs/swagger.json
	@rg --passthru -U 'additionalProperties:\s*type:\s*integer' -r 'additionalProperties: {"type": "string", "enum": ["Contributor", "User"]}' docs/swagger.yaml | sponge docs/swagger.yaml

pre-build:
	$(MAKE) swagger

build:
	@go build -o smart-wallet\
			  -ldflags "-w\
	                    -s\
						-X '${PKG}/version.gitCommit=$$(git rev-parse HEAD)'\
						-X '${PKG}/version.buildTime=$$(date)'\
						-X '${PKG}/version.goVersion=$$(go version)'\
						-X '${PKG}/version.tag=$$(git describe --tags)'"\
			  cmd/smart-wallet/main.go

test:
	@go test -short ${PKG_LIST}

race:
	@go test -race -short ${PKG_LIST}

msan:
	@go test -msan -short ${PKG_LIST}

coverage:
	@go test -cover -coverprofile=cover/coverage.out ${PKG_LIST}

coveragehtml:
	@go tool cover -html=cover/coverage.out -o coverage.html
