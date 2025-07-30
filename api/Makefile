.PHONY: test test-unit test-integration test-coverage run build clean deps docker-test

# Configurações
GO_VERSION := 1.21
PROJECT_NAME := besu-api

# Comandos principais
run:
	go run cmd/main.go

build:
	go build -o bin/$(PROJECT_NAME) cmd/main.go

clean:
	rm -rf bin/
	go clean -testcache

# Instalar dependências
deps:
	go mod download
	go mod tidy

# Instalar dependências de teste
deps-test:
	go get github.com/stretchr/testify/assert
	go get github.com/stretchr/testify/mock
	go get github.com/stretchr/testify/suite

# Executar todos os testes
test: deps-test
	go test ./... -v

# Testes unitários apenas
test-unit: deps-test
	go test ./internal/... -v -short

# Testes de integração (requer infraestrutura)
test-integration: deps-test
	go test ./tests/... -v

# Cobertura de testes
test-coverage: deps-test
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Cobertura simplificada
test-cover: deps-test
	go test ./... -cover

# Executar testes com race detection
test-race: deps-test
	go test ./... -race -v

# Benchmark tests
test-bench:
	go test ./... -bench=. -benchmem

# Linting (requer golangci-lint)
lint:
	golangci-lint run

# Formatação
fmt:
	go fmt ./...

# Vet (análise estática)
vet:
	go vet ./...

# Executar todos os checks
check: fmt vet lint test

# Docker para testes (com PostgreSQL)
docker-test:
	docker-compose -f docker-compose.test.yml up -d
	sleep 5
	go test ./... -v
	docker-compose -f docker-compose.test.yml down

# Executar servidor com hot reload (requer air)
dev:
	air

# Gerar mocks (requer mockgen)
mocks:
	mockgen -source=internal/blockchain/client.go -destination=mocks/blockchain_mock.go
	mockgen -source=internal/database/database.go -destination=mocks/database_mock.go

# Ajuda
help:
	@echo "Comandos disponíveis:"
	@echo "  run              - Executar aplicação"
	@echo "  build            - Build da aplicação"
	@echo "  test             - Executar todos os testes"
	@echo "  test-unit        - Executar apenas testes unitários"
	@echo "  test-integration - Executar testes de integração"
	@echo "  test-coverage    - Gerar relatório de cobertura"
	@echo "  test-race        - Executar testes com race detection"
	@echo "  lint             - Executar linting"
	@echo "  fmt              - Formatar código"
	@echo "  check            - Executar todos os checks"
	@echo "  deps             - Instalar dependências"
	@echo "  clean            - Limpar arquivos de build"