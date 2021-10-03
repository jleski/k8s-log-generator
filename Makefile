BIN_NAME := loggen
BIN_NAME_LN := loggen_linux
SRC := $(wildcard *.go)
CONTAINER_IMAGE=jledev.azurecr.io/loggen:latest
SHELL:="/bin/bash"

.PHONY: build_and_run

build_and_run: $(BIN_NAME) run

container: $(BIN_NAME_LN)
	docker build -t $(CONTAINER_IMAGE) .

$(BIN_NAME): $(SRC)
	go build -ldflags="-s -w" -o $(BIN_NAME) .

run: $(BIN_NAME)
	./$(BIN_NAME) -interval 2

$(BIN_NAME_LN): $(SRC)
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BIN_NAME_LN) .

build:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main main.go


push:
	# for this to work you need to have logged in to the registry using `docker login`
	docker push $(CONTAINER_IMAGE)

kubesecret:
	# first do docker login prior running this target
	kubectl create secret generic jledev-azurecr-cred \
		--from-file=.dockerconfigjson=$(HOME)/.docker/config.json \
		--type=kubernetes.io/dockerconfigjson

deploy:
	# for this to work you need to have authenticated to a Kubernetes cluster
	# and set the desired namespace as default for the cluster context
	kubectl apply -f kubernetes/.
