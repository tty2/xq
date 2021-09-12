##@ Test
lint: ## Run linters only.
	@echo -e "\033[2m→ Running linters...\033[0m"
	golangci-lint run --config .golangci.yml

test: ## Run go tests for files with tests.
	@echo -e "\033[2m→ Run tests for all files...\033[0m"
	@if [ $$(cat fixtures/hashes.txt | sed -n 3p) = $$(cat fixtures/film.xml | go run main.go processor.go query.go fullparse.go | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 4p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go fullparse.go | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 5p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go fullparse.go tags .PurchaseOrder | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 6p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go fullparse.go tags .PurchaseOrder.Items | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 7p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go fullparse.go tags .PurchaseOrder.Items.Item | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 8p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go fullparse.go tags .PurchaseOrder.Address | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 9p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go fullparse.go tags .PurchaseOrder.DeliveryNotes | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 9p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go fullparse.go tags .PurchaseOrder.Address.Name | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 9p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go fullparse.go tags .PurchaseOrder.Address.Street | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 9p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go fullparse.go tags .PurchaseOrder.Address.City | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 9p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go fullparse.go tags .PurchaseOrder.Address.State | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 9p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go fullparse.go tags .PurchaseOrder.Address.Zip | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 9p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go fullparse.go tags .PurchaseOrder.Address.Country | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 9p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go fullparse.go tags .PurchaseOrder.Items.Item.ProductName | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 9p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go fullparse.go tags .PurchaseOrder.Items.Item.Quantity | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 9p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go fullparse.go tags .PurchaseOrder.Items.Item.USPrice | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 9p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go fullparse.go tags .PurchaseOrder.Items.Item.Comment | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 9p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go fullparse.go tags .PurchaseOrder.Items.Item.image | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 9p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go fullparse.go tags .PurchaseOrder.Items.Item.ShipDate | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi

beautify:
	gofumpt -l -w $$(go list -f {{.Dir}} ./... | grep -v /vendor/)


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
	echo "These hashes don't related with any security information but only xq output hashes for fixtures generated by 'make hash-gen' command" > fixtures/hashes.txt
	echo "" >> fixtures/hashes.txt
	cat fixtures/film.xml | go run main.go processor.go query.go fullparse.go | md5sum | awk '{print $$1}' >>  fixtures/hashes.txt
	cat fixtures/flat.xml | go run main.go processor.go query.go fullparse.go | md5sum | awk '{print $$1}' >> fixtures/hashes.txt
	cat fixtures/flat.xml | go run main.go processor.go query.go fullparse.go tags .PurchaseOrder | md5sum | awk '{print $$1}' >> fixtures/hashes.txt
	cat fixtures/flat.xml | go run main.go processor.go query.go fullparse.go tags .PurchaseOrder.Items | md5sum | awk '{print $$1}' >> fixtures/hashes.txt
	cat fixtures/flat.xml | go run main.go processor.go query.go fullparse.go tags .PurchaseOrder.Items.Item | md5sum | awk '{print $$1}' >> fixtures/hashes.txt
	cat fixtures/flat.xml | go run main.go processor.go query.go fullparse.go tags .PurchaseOrder.Address | md5sum | awk '{print $$1}' >> fixtures/hashes.txt
	cat fixtures/flat.xml | go run main.go processor.go query.go fullparse.go tags .PurchaseOrder.DeliveryNotes | md5sum | awk '{print $$1}' >> fixtures/hashes.txt

gen: ## Generate result to `generate` folder
	@echo -e "\033[2m→ Generating test files...\033[0m"
	cat fixtures/film.xml | go run main.go processor.go query.go fullparse.go > generate/film.xml
	cat fixtures/flat.xml | go run main.go processor.go query.go fullparse.go > generate/flat.xml

##@ Other
#------------------------------------------------------------------------------
help:  ## Display help
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
#------------- <https://suva.sh/posts/well-documented-makefiles> --------------

.DEFAULT_GOAL := help
.PHONY: help lint test check build install hash-gen
