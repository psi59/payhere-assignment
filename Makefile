.PHONY: vendor build deploy

version=$(shell git rev-parse --short HEAD)
image=payhere

vendor:
	go mod tidy
	go mod vendor

test:
	rm -f cover.out
	go test -v -coverprofile=cover.out ./...
	go tool cover -html=cover.out

build:
	DOCKER_BUILDKIT=1 docker build --build-arg=GOARCH=arm64 -t $(image):$(version) .

run:
	docker run \
	--name payhere \
	-p 127.0.0.1:1202:1202 \
	-v /Users/psi59/Workspace/payhere-assignment/config/serve.docker.yaml:/www/config.yaml \
	-d \
	$(image):$(version) \
	serve \
	--config-path=/www/config.yaml