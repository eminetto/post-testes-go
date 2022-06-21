build: mod
	go build -o bin/go-test-workshop

mod:
	go mod download

unit-test: clean generate-mocks
	go test -tags=unit -short -coverprofile=cp.out ./...

integration:
	go test -v ./... -tags=integration -coverprofile=cp.out

e2e:
	go test -v ./... -tags=e2e -coverprofile=cp.out

generate-mocks:
	@mockery --output usecase/person/mocks --dir usecase/person --all
	@mockery --output usecase/weather/mocks --dir usecase/weather --all

clean:
	@rm -rf usecase/person/mocks/*
	@rm -rf usecase/weather/mocks/*
