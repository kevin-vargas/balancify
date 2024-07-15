dev:
	@docker build --progress=plain -t test -f ./ias/Dockerfile .
dev-run:
	@docker run -p 80:8080 test
start:
	@cd app && npm run build -- --mode staging
	@go run cmd/main.go
