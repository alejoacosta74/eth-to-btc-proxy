.DEFAULT_GOAL := help
.PHONY: all build run test clean help regtest


## This help screen
help:
	@printf "Available targets:\n\n"
	@awk '/^[a-zA-Z\-\_0-9%:\\]+/ { \
	helpMessage = match(lastLine, /^## (.*)/); \
	if (helpMessage) { \
	helpCommand = $$1; \
	helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
	gsub("\\\\", "", helpCommand); \
	gsub(":+$$", "", helpCommand); \
	printf "  \x1b[32;01m%-35s\x1b[0m %s\n", helpCommand, helpMessage; \
	} \
        } \
        { lastLine = $$0 }' $(MAKEFILE_LIST) | sort -u
	@printf "\n"

# check if golang is installed
go/check:
	@ which go > /dev/null || (echo "golang is not installed" && exit 1)


## Run integration tests against a docker qtum regtest node
integration-test: go/check
	@ printf "\nRunning integration tests...\n\n"
	@ ./test/integration.sh
## Run unit tests
unit-test: go/check
	@ printf "\nRunning tests...\n\n"
	@ go test -v ./...
## Build qproxy
build: ## build qproxy
	@ printf "\nBuilding qproxy...\n\n"
	@ go build -o bin/qproxy main.go

## Run qproxy
run: go/check ## run qproxy
	@ printf "\nRunning qproxy...\n\n"
	@ ./bin/qproxy

## clean everything
clean: 
	@ printf "\nRemoving files and folders...\n\n"
	@ rm -rf bin > /dev/null 2>&1 || true
	@ docker rm -f regtest > /dev/null 2>&1 || true
	@ rm -rf regtest > /dev/null 2>&1 || true
	@ rm -rf ethcli > /dev/null 2>&1 || true
