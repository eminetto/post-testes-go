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
	@mockery --output person/mocks --dir person --all
	@mockery --output weather/mocks --dir weather --all

clean:
	@rm -rf person/mocks/*
	@rm -rf weather/mocks/*
