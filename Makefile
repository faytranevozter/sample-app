run: 
	@go run app/main.go

test:
	@go test -v -cover -covermode=atomic ./...

testcoverage:
	@go test -v -cover -covermode=atomic -coverprofile=coverage.out ./... && \
	go tool cover -func=coverage.out

generate-mocks:
	@rm -rf mocks && \
	mockery --all

build: 
	@go build -ldflags="-s -w" -o build-app app/main.go
