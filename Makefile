.PHONY: vendor build deploy

vendor:
	go mod tidy
	go mod vendor
