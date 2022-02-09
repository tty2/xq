##@ Test
lint: ## Run linters only.
	@echo -e "\033[2m→ Running linters...\033[0m"
	golangci-lint run --config .golangci.yml

test: ## Run go tests for files with tests.
	@echo -e "\033[2m→ Run tests for all files...\033[0m"
	go test -v ./...
	@if [ $$(cat fixtures/hashes.txt | sed -n 3p) = $$(cat fixtures/film.xml | go run main.go processor.go query.go | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 4p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 5p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go tags .PurchaseOrder | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 6p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go tags .PurchaseOrder.Items | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 7p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go tags .PurchaseOrder.Items.Item | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 8p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go tags .PurchaseOrder.Address | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 9p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go tags .PurchaseOrder.DeliveryNotes | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 9p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go tags .PurchaseOrder.Address.Name | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 9p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go tags .PurchaseOrder.Address.Street | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 9p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go tags .PurchaseOrder.Address.City | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 9p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go tags .PurchaseOrder.Address.State | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 9p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go tags .PurchaseOrder.Address.Zip | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 9p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go tags .PurchaseOrder.Address.Country | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 9p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go tags .PurchaseOrder.Items.Item.ProductName | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 9p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go tags .PurchaseOrder.Items.Item.Quantity | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 9p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go tags .PurchaseOrder.Items.Item.USPrice | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 9p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go tags .PurchaseOrder.Items.Item.Comment | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 9p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go tags .PurchaseOrder.Items.Item.image | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 9p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go tags .PurchaseOrder.Items.Item.ShipDate | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 10p) = $$(cat fixtures/film.xml | go run main.go processor.go query.go tags .objects.object | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 11p) = $$(cat fixtures/film.xml | go run main.go processor.go query.go attr .objects.object | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 12p) = $$(cat fixtures/film.xml | go run main.go processor.go query.go attr .objects.object.poster | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 13p) = $$(cat fixtures/film.xml | go run main.go processor.go query.go attr .objects.object.poster#url | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 14p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go attr .PurchaseOrder.Address#Type | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 15p) = $$(cat fixtures/flat.xml | go run main.go processor.go query.go .PurchaseOrder.Address.Name | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 16p) = $$(cat fixtures/film.xml | go run main.go processor.go query.go .objects.object.actors.actor[0] | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 17p) = $$(cat fixtures/film.xml | go run main.go processor.go query.go .objects.object.actors.actor[6] | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 18p) = $$(cat fixtures/film.xml | go run main.go processor.go query.go .objects[0].object[0].actors.actor[1] | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 19p) = $$(cat fixtures/film.xml | go run main.go processor.go query.go .objects[1].object.actors.actor | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 20p) = $$(cat fixtures/film.xml | go run main.go processor.go query.go .objects.object.poster[1] | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 21p) = $$(cat fixtures/film.xml | go run main.go processor.go query.go attr .objects.object.poster[0] | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 22p) = $$(cat fixtures/film.xml | go run main.go processor.go query.go attr .objects.object.poster[0]#url | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi
	@if [ $$(cat fixtures/hashes.txt | sed -n 23p) = $$(cat fixtures/film.xml | go run main.go processor.go query.go tags .objects.object[0].actors | md5sum | awk '{print $$1}') ]; then echo "PASSED"; else exit 125; fi

	cat fixtures/film.xml | go run -race main.go processor.go query.go > /dev/null

beautify:
	gofumpt -l -w ./$$(go list -f {{.Dir}} ./... | grep -v /vendor/)


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
	echo "These hashes are not related with any security information but only xq output hashes for fixtures generated by 'make hash-gen' command" > fixtures/hashes.txt
	echo "" >> fixtures/hashes.txt
	cat fixtures/film.xml | go run main.go processor.go query.go | md5sum | awk '{print $$1}' >>  fixtures/hashes.txt
	cat fixtures/flat.xml | go run main.go processor.go query.go | md5sum | awk '{print $$1}' >> fixtures/hashes.txt
	cat fixtures/flat.xml | go run main.go processor.go query.go tags .PurchaseOrder | md5sum | awk '{print $$1}' >> fixtures/hashes.txt
	cat fixtures/flat.xml | go run main.go processor.go query.go tags .PurchaseOrder.Items | md5sum | awk '{print $$1}' >> fixtures/hashes.txt
	cat fixtures/flat.xml | go run main.go processor.go query.go tags .PurchaseOrder.Items.Item | md5sum | awk '{print $$1}' >> fixtures/hashes.txt
	cat fixtures/flat.xml | go run main.go processor.go query.go tags .PurchaseOrder.Address | md5sum | awk '{print $$1}' >> fixtures/hashes.txt
	cat fixtures/flat.xml | go run main.go processor.go query.go tags .PurchaseOrder.DeliveryNotes | md5sum | awk '{print $$1}' >> fixtures/hashes.txt
	cat fixtures/film.xml | go run main.go processor.go query.go tags .objects.object | md5sum | awk '{print $$1}' >>  fixtures/hashes.txt
	cat fixtures/film.xml | go run main.go processor.go query.go attr .objects.object | md5sum | awk '{print $$1}' >>  fixtures/hashes.txt
	cat fixtures/film.xml | go run main.go processor.go query.go attr .objects.object.poster | md5sum | awk '{print $$1}' >>  fixtures/hashes.txt
	cat fixtures/film.xml | go run main.go processor.go query.go attr .objects.object.poster#url | md5sum | awk '{print $$1}' >>  fixtures/hashes.txt
	cat fixtures/flat.xml | go run main.go processor.go query.go attr .PurchaseOrder.Address#Type | md5sum | awk '{print $$1}' >>  fixtures/hashes.txt
	cat fixtures/flat.xml | go run main.go processor.go query.go .PurchaseOrder.Address.Name | md5sum | awk '{print $$1}' >>  fixtures/hashes.txt
	cat fixtures/film.xml | go run main.go processor.go query.go .objects.object.actors.actor[0] | md5sum | awk '{print $$1}' >>  fixtures/hashes.txt
	cat fixtures/film.xml | go run main.go processor.go query.go .objects.object.actors.actor[6] | md5sum | awk '{print $$1}' >>  fixtures/hashes.txt
	cat fixtures/film.xml | go run main.go processor.go query.go .objects[0].object[0].actors.actor[1] | md5sum | awk '{print $$1}' >>  fixtures/hashes.txt
	cat fixtures/film.xml | go run main.go processor.go query.go .objects[1].object.actors.actor | md5sum | awk '{print $$1}' >>  fixtures/hashes.txt
	cat fixtures/film.xml | go run main.go processor.go query.go .objects.object.poster[1] | md5sum | awk '{print $$1}' >>  fixtures/hashes.txt
	cat fixtures/film.xml | go run main.go processor.go query.go attr .objects.object.poster[0] | md5sum | awk '{print $$1}' >>  fixtures/hashes.txt
	cat fixtures/film.xml | go run main.go processor.go query.go attr .objects.object.poster[0]#url | md5sum | awk '{print $$1}' >>  fixtures/hashes.txt
	cat fixtures/film.xml | go run main.go processor.go query.go tags .objects.object[0].actors | md5sum | awk '{print $$1}' >>  fixtures/hashes.txt

gen: ## Generate result to `generate` folder
	@echo -e "\033[2m→ Generating test files...\033[0m"
	cat fixtures/film.xml | go run main.go processor.go query.go > generate/film.xml
	cat fixtures/flat.xml | go run main.go processor.go query.go > generate/flat.xml

##@ Other
#------------------------------------------------------------------------------
help:  ## Display help
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
#------------- <https://suva.sh/posts/well-documented-makefiles> --------------

.DEFAULT_GOAL := help
.PHONY: help lint test check build install hash-gen
