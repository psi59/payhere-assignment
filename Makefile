.PHONY: vendor build deploy

vendor:
	go mod tidy
	go mod vendor

test:
	rm -f cover.out
	go test -v -coverpkg=./... -coverprofile=cover.out ./...
	go tool cover -html=cover.out