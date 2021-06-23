##@ Test
lint: ## Run linters only.
	@echo -e "\033[2m→ Running linters...\033[0m"
	golangci-lint run --config .golangci.yml

test: ## Run go tests for files with tests.
	@echo -e "\033[2m→ Run tests for all files...\033[0m"
	@if [ $$(cat fixtures/hashes.txt | sed -n 1p) = $$(cat fixtures/film.xml | go run main.go processor.go | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 2p) = $$(cat fixtures/test.xml | go run main.go processor.go | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi


check: lint test ## Run full check: lint and test.

##@ Deploy
build: ## Build library.
	@echo -e "\033[2m→ Building library...\033[0m"
	go build

install: ## Install library.
	@echo -e "\033[2m→ Installing library...\033[0m"
	go install

##@ Generate
hash-gen: ## Gererates result files for tests
	@echo -e "\033[2m→ Generating result files for tests...\033[0m"
	cat fixtures/film.xml | go run main.go processor.go | md5sum | awk '{print $$1}' >  fixtures/hashes.txt
	cat fixtures/test.xml | go run main.go processor.go | md5sum | awk '{print $$1}' >> fixtures/hashes.txt

##@ Other
#------------------------------------------------------------------------------
help:  ## Display help
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
#------------- <https://suva.sh/posts/well-documented-makefiles> --------------

.DEFAULT_GOAL := help
.PHONY: help lint test check build install hash-gen
