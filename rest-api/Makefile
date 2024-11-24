test:
	go test -cover -v ./...

check-swag:
	@which swag || go install github.com/swaggo/swag/cmd/swag@latest

swag: check-swag
	@swag init -d ./cmd/,./controller --parseDependency --parseInternal --parseDepth 1