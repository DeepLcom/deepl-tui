GIT_DIR := `git rev-parse --show-toplevel`

MAIN := "."
BIN_NAME := `basename $(git rev-parse --show-toplevel)`
BIN_DIR := "bin"
DIST_DIR := "dist"

# list available recipes
default:
    @just --list

# format code
fmt:
    go fmt ./...

# lint code
lint:
    golangci-lint run ./...

# vet code
vet:
    go vet ./...

# build application
build *args="":
    {{GIT_DIR}}/scripts/build.sh -p {{MAIN}} {{args}}

# create binary distribution
dist *args="":
    {{GIT_DIR}}/scripts/dist.sh -p {{MAIN}} {{args}}

# create a new release
release *args="":
    #!/bin/sh
    export GITHUB_OWNER=DeepLcom
    export GITHUB_REPO=deepl-tui
    {{GIT_DIR}}/scripts/release.sh -p {{MAIN}} {{args}}

changes from="" to="":
    #!/bin/sh
    source {{GIT_DIR}}/scripts/functions.sh
    get_changes {{from}} {{to}}

clean:
    @# build artifacts
    @echo "rm {{BIN_DIR}}/{{BIN_NAME}}"
    @-[ -f {{BIN_DIR}}/{{BIN_NAME}} ] && rm {{BIN_DIR}}/{{BIN_NAME}}
    @-[ -d {{BIN_DIR}} ] && rmdir {{BIN_DIR}}

    @# distribution binaries
    @echo "rm {{DIST_DIR}}/{{BIN_NAME}}_*"
    @rm {{DIST_DIR}}/{{BIN_NAME}}_* 2>/dev/null || true
    @-[ -d {{DIST_DIR}} ] && rmdir {{DIST_DIR}}

