.PHONY: all

CURDIR := $(shell pwd)

help: ## Print this help
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build Docker Container
	gcloud builds submit -t gcr.io/${GOOGLE_CLOUD_PROJECT}/serverless_example:1.0.0 .

apply: ## Apply to k8s
	@cd ./k8s; \
	kustomize build . | \
	sed 's/\$$APP_IMAGE'"/gcr.io\/${GOOGLE_CLOUD_PROJECT}\/serverless_example:1.0.0/g" | \
	kubectl apply -f -

delete: ## Delete resources from k8s
	kubectl delete -f ./k8s/deployment.yaml -f ./k8s/service.yaml
