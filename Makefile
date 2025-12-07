.PHONY: help test test-unit test-integration test-web test-mobile test-watch lint fmt terraform-validate terraform-plan

help:
	@echo "Crowd Unlocked - Available targets:"
	@echo ""
	@echo "Testing:"
	@echo "  test              - Run all tests (unit + integration)"
	@echo "  test-unit         - Run unit tests only"
	@echo "  test-integration  - Run integration tests"
	@echo "  test-web          - Run web app tests"
	@echo "  test-mobile       - Run mobile app tests"
	@echo "  test-watch        - Run tests in watch mode"
	@echo ""
	@echo "Code Quality:"
	@echo "  lint              - Run linters"
	@echo "  fmt               - Format code"
	@echo ""
	@echo "Infrastructure:"
	@echo "  terraform-validate - Validate OpenTofu configs"
	@echo "  terraform-plan     - Plan OpenTofu changes"
	@echo ""
	@echo "Build & Deploy:"
	@echo "  build             - Build all services"
	@echo "  docker-build      - Build all Docker images"
	@echo "  deploy-dev        - Deploy to dev cluster"
	@echo "  deploy-prod       - Deploy to prod cluster"

# Run all tests
test: test-unit test-web
	@echo "✅ All tests passed!"

# Unit tests for Go services
test-unit:
	@echo "Running unit tests for Go services..."
	@cd services/bookings && go test -v -cover ./...
	@cd services/releases && go test -v -cover ./...
	@cd services/publicity && go test -v -cover ./...
	@cd services/social && go test -v -cover ./...
	@cd services/money && go test -v -cover ./...

# Integration tests (requires Docker)
test-integration:
	@echo "Starting test dependencies..."
	@docker-compose up -d dynamodb-local
	@echo "Running integration tests..."
	@cd services/bookings && go test -tags=integration -v ./tests/integration/...
	@cd services/releases && go test -tags=integration -v ./tests/integration/...
	@cd services/publicity && go test -tags=integration -v ./tests/integration/...
	@cd services/social && go test -tags=integration -v ./tests/integration/...
	@cd services/money && go test -tags=integration -v ./tests/integration/...
	@docker-compose down

# Web app tests
test-web:
	@echo "Running web app tests..."
	@cd apps/web && npm test

# Mobile app tests
test-mobile:
	@echo "Running mobile app tests..."
	@cd apps/mobile && flutter test

# Watch mode for rapid feedback
test-watch:
	@echo "Running tests in watch mode..."
	@echo "Press Ctrl+C to stop"
	@find services -name "*.go" | entr -c make test-unit

# Linting
lint:
	@echo "Running linters..."
	@cd services/bookings && golangci-lint run
	@cd services/releases && golangci-lint run
	@cd services/publicity && golangci-lint run
	@cd services/social && golangci-lint run
	@cd services/money && golangci-lint run
	@cd apps/web && npm run lint

# Format code
fmt:
	@echo "Formatting code..."
	@cd services/bookings && go fmt ./...
	@cd services/releases && go fmt ./...
	@cd services/publicity && go fmt ./...
	@cd services/social && go fmt ./...
	@cd services/money && go fmt ./...
	@cd apps/web && npm run format
	@tofu fmt -recursive infra/

# OpenTofu validation
terraform-validate:
	@echo "Validating OpenTofu configurations..."
	@cd infra/terraform/management && tofu init -backend=false && tofu validate
	@cd infra/terraform/dev && tofu init -backend=false && tofu validate
	@cd infra/terraform/prod && tofu init -backend=false && tofu validate

# OpenTofu plan (requires AWS credentials and tfvars)
terraform-plan:
	@echo "Planning OpenTofu changes..."
	@cd infra/terraform/management && tofu plan -var-file=mgmt.tfvars
	@cd infra/terraform/dev && tofu plan -var-file=dev.tfvars
	@cd infra/terraform/prod && tofu plan -var-file=prod.tfvars

build:
	@echo "Building all services..."
	@cd services/bookings && go build -o ../../bin/bookings ./cmd/server
	@cd services/releases && go build -o ../../bin/releases ./cmd/server
	@cd services/publicity && go build -o ../../bin/publicity ./cmd/server
	@cd services/social && go build -o ../../bin/social ./cmd/server
	@cd services/money && go build -o ../../bin/money ./cmd/server

docker-build:
	@echo "Building Docker images..."
	@docker build -t crowdunlocked/bookings:latest -f services/bookings/Dockerfile .
	@docker build -t crowdunlocked/releases:latest -f services/releases/Dockerfile .
	@docker build -t crowdunlocked/publicity:latest -f services/publicity/Dockerfile .
	@docker build -t crowdunlocked/social:latest -f services/social/Dockerfile .
	@docker build -t crowdunlocked/money:latest -f services/money/Dockerfile .
