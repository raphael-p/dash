main_package_path = ./cmd/dash
binary_name = dsh
version=$(shell git describe --tags --always --dirty)
release_version ?=


# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty
no-dirty:
	@test -z "$(shell git status --porcelain)"

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: run quality control checks
.PHONY: audit
audit: test
	go mod tidy -diff
	go mod verify
	test -z "$(shell gofmt -l .)" 
	go vet ./...
	go tool staticcheck -checks=all,-ST1000,-U1000 ./...
	go tool govulncheck ./...

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

## upgradeable: list direct dependencies that have upgrades available
.PHONY: upgradeable
upgradeable:
	go run github.com/oligot/go-mod-upgrade@latest

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## tidy: tidy modfiles and modernize and format .go files
.PHONY: tidy
tidy:
	go mod tidy -v
	go fix ./...
	go fmt ./...

## build: build the application
.PHONY: build
build:
	go build -ldflags="-X 'main.AppVersion=$(version)'" -o=/tmp/bin/${binary_name} ${main_package_path}

## run: run the  application
.PHONY: run
run: build
	/tmp/bin/${binary_name}

# ==================================================================================== #
# OPERATIONS
# ==================================================================================== #

## push: push changes to the remote Git repository
.PHONY: push
push: confirm audit no-dirty
	git push

## release: create and push a version tag, then publish with GoReleaser
.PHONY: release
release: release-version-check audit no-dirty
	@test -z "$$(git rev-parse --verify --quiet "refs/tags/v$(release_version)")" || (echo "tag v$(release_version) already exists locally" >&2; exit 1)
	@test -z "$$(git ls-remote --tags origin "refs/tags/v$(release_version)")" || (echo "tag v$(release_version) already exists on origin" >&2; exit 1)
	@gh auth status >/dev/null
	git tag -a "v$(release_version)" -m "Release v$(release_version)"
	git push origin "v$(release_version)"
	GITHUB_TOKEN="$$(gh auth token)" goreleaser release --clean

.PHONY: release-version-check
release-version-check:
	@test -n "$(release_version)" || (echo 'usage: make release release_version=0.1.0' >&2; exit 1)
	@case "$(release_version)" in v*) echo 'release_version must not start with v' >&2; exit 1;; esac
	@git check-ref-format "refs/tags/v$(release_version)"

## production/deploy: deploy the application to production
.PHONY: production/deploy
production/deploy: confirm audit no-dirty
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -X 'main.AppVersion=$(version)'" -o=/tmp/bin/darwin_arm64/${binary_name} ${main_package_path}
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -X 'main.AppVersion=$(version)'" -o=/tmp/bin/darwin_amd64/${binary_name} ${main_package_path}
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -X 'main.AppVersion=$(version)'" -o=/tmp/bin/linux_amd64/${binary_name} ${main_package_path}
