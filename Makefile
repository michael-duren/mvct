test:
	@echo "Testing"
	
	@go test -v ./...
	@echo "Tests completed"
	@echo ""

scaffold:
	@echo "Scaffolding"
	@go run ./cmd/cli/main.go

build:
	@echo "Building"
	@go build -o bin/app ./cmd/cli/main.go
	@echo "Build completed"
	@echo ""
