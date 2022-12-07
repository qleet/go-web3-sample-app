.DEFAULT_GOAL := help
CURRENTTAG:=$(shell git describe --tags --abbrev=0)
NEWTAG ?= $(shell bash -c 'read -p "Please provide a new tag (currnet tag - ${CURRENTTAG}): " newtag; echo $$newtag')
GOFLAGS=-mod=mod
KINDCLUSTER=qleet-os

#help: @ List available tasks
help:
	@clear
	@echo "Usage: make COMMAND"
	@echo "Commands :"
	@grep -E '[a-zA-Z\.\-]+:.*?@ .*$$' $(MAKEFILE_LIST)| tr -d '#' | awk 'BEGIN {FS = ":.*?@ "}; {printf "\033[32m%-13s\033[0m - %s\n", $$1, $$2}'

#clean: @ Cleanup
clean:
	@rm -rf ./dist
	@rm -rf ./completions

#test: @ Run tests
test:
	@go generate
	@export GOFLAGS=$(GOFLAGS); go test $(go list ./...)

#build: @ Build binary
build:
	@export GOFLAGS=$(GOFLAGS); export CGO_ENABLED=0; go build -a -o go-web3-sample-app main.go

#run: @ Run binary
run:
	@export RPCENDPOINT=https://rpc.ankr.com/eth; export GOFLAGS=$(GOFLAGS); go run main.go

#get: @ Download and install dependency packages
get:
	@export GOFLAGS=$(GOFLAGS); go get . ; go mod tidy

#release: @ Create and push a new tag
release: build
	$(eval NT=$(NEWTAG))
	@echo -n "Are you sure to create and push ${NT} tag? [y/N] " && read ans && [ $${ans:-N} = y ]
	@echo ${NT} > ./version.txt
	@git add -A
	@git commit -a -s -m "Cut ${NT} release"
	@git tag -a -m "Cut ${NT} release" ${NT}
	@git push origin ${NT}
	@git push
	@echo "Done."

#update: @ Update dependencies to latest versions
update:
	@export GOFLAGS=$(GOFLAGS); go get -u; go mod tidy

#version: @ Print current version(tag)
version:
	@echo $(shell git describe --tags --abbrev=0)

#image-build: @ Build a Docker image
image-build:
	docker build -t go-web3-sample-app:$(CURRENTTAG) .

#image-run: @ Run a Docker image
image-run: image-stop
	@export RPCENDPOINT=https://rpc.ankr.com/eth && docker-compose -f "docker-compose.yml" up --build
#	@docker run --rm -p 8080:8080 --name web3 go-web3-sample-app:$(CURRENTTAG) -e RPCENDPOINT='https://rpc.ankr.com/eth'
#	@docker run --rm --env RPCENDPOINT=https://rpc.ankr.com/eth alpine:3 env

#image-stop: @ Stop a Docker image
image-stop:
	@docker-compose -f "docker-compose.yml" down --volumes
#	@docker stop web3 || true

#kind-deploy: @ Deploy to a local KinD cluster
kind-deploy: image-build
	@kind load docker-image go-web3-sample-app:$(CURRENTTAG) -n $(KINDCLUSTER) && \
	cat ./k8s/ns.yaml | kubectl apply -f - && \
	cat ./k8s-kind-qleet/cm.yaml | kubectl apply --namespace=web3 -f - && \
	yq eval '.spec.template.spec.containers[0].image = "go-web3-sample-app:$(CURRENTTAG)"' ./k8s/deployment.yaml | kubectl apply --namespace=web3 -f - && \
	cat ./k8s/service.yaml | kubectl apply --namespace=web3 -f -

#kind-undeploy: @ Undeploy from a local KinD cluster
kind-undeploy:
	@kubectl delete -f ./k8s/deployment.yaml --namespace=web3 --ignore-not-found=true && \
	kubectl delete -f ./k8s/cm.yaml --namespace=web3 --ignore-not-found=true && \
	kubectl delete -f ./k8s/ns.yaml --ignore-not-found=true

#kind-redeploy: @ Redeploy to a local KinD cluster
kind-redeploy:
	@kubectl delete -f ./k8s/deployment.yaml --namespace=web3 --ignore-not-found=true && \
	kubectl apply -f ./k8s/cm.yaml --namespace=web3 && \
	yq eval '.spec.template.spec.containers[0].image = "go-web3-sample-app:$(CURRENTTAG)"' ./k8s/deployment.yaml | kubectl apply --namespace=web3 -f -

# ssh into pod
# kubectl exec --stdin --tty -n web3 web3-sample-app-569598dd94-qvg4m -- /bin/sh

# pod logs
# kubectl logs -n web3 web3-sample-app-569598dd94-qvg4m
