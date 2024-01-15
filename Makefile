.PHONY: vendor build run

vendor:
	go mod tidy
	go mod vendor

mockgen:
	go install go.uber.org/mock/mockgen@latest
	mockgen -source usecase/user/interface.go -typed -destination internal/mocks/ucmocks/user_usecase.go -mock_names=Usecase=MockUserUsecase -package ucmocks
	mockgen -source usecase/authtoken/interface.go -typed -destination internal/mocks/ucmocks/authtoken_usecase.go -mock_names=Usecase=MockAuthTokenUsecase -package ucmocks
	mockgen -source usecase/item/interface.go -typed -destination internal/mocks/ucmocks/item_usecase.go -mock_names=Usecase=MockItemTokenUsecase -package ucmocks
	mockgen -source repository/interface.go -typed -destination internal/mocks/repomocks/repository.go -package repomocks

test: mockgen
	rm -f cover.out
	go test -v -coverprofile=cover.out ./handler/... ./usecase/... ./repository/...
	go tool cover -html=cover.out

build: vendor
	docker compose build

run:
	docker compose up