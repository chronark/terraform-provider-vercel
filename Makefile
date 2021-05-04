default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	@go get github.com/mfridman/tparse
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m  -race -covermode=atomic  -json | tparse -all -dump

test:
	go test ./... -v

build: fmt
	go build -o terraform-provider-vercel
	mkdir -p ~/.terraform.d/plugins/hashicorp.com/chronark/vercel/9000.1/linux_amd64
	mv terraform-provider-vercel ~/.terraform.d/plugins/hashicorp.com/chronark/vercel/9000.1/linux_amd64


fmt:

	go generate -v ./...
	golangci-lint run -v
	go fmt ./...
	terraform fmt -recursive .

rm-state:
	rm -rf examples/e2e/terraform.tfstate*


init: build
	rm -rf examples/e2e/.terraform*
	terraform -chdir=examples/e2e init -upgrade
e2e: init
	terraform -chdir=examples/e2e apply


release:
	@go get github.com/caarlos0/svu
	@echo "Releasing $$(svu next)..."
	
	@git tag $$(svu next) && git push --tags
	@echo "Done"