.PHONY: help install test lint fmt type-check build terraform-validate terraform-plan

help:
	@echo "Crowd Unlocked - Available targets:"
	@echo ""
	@echo "Web app (apps/web):"
	@echo "  install            - Install web app dependencies"
	@echo "  test               - Run web app tests"
	@echo "  lint               - Lint the web app"
	@echo "  type-check         - Type-check the web app"
	@echo "  build              - Build the web app"
	@echo "  fmt                - Format code (web + terraform)"
	@echo ""
	@echo "Mobile app (apps/mobile):"
	@echo "  test-mobile        - Run Flutter tests"
	@echo ""
	@echo "Infrastructure (infra/terraform/dev):"
	@echo "  terraform-validate - Validate OpenTofu configs"
	@echo "  terraform-plan     - Plan OpenTofu changes (needs AWS creds)"

install:
	@cd apps/web && npm install

# Web app tests (Jest)
test:
	@echo "Running web app tests..."
	@cd apps/web && npm test

lint:
	@echo "Linting web app..."
	@cd apps/web && npm run lint

type-check:
	@echo "Type-checking web app..."
	@cd apps/web && npm run type-check

build:
	@echo "Building web app..."
	@cd apps/web && npm run build

# Mobile app tests
test-mobile:
	@echo "Running mobile app tests..."
	@cd apps/mobile && flutter test

# Format code (web app + terraform)
fmt:
	@echo "Formatting code..."
	@cd apps/web && npm run format
	@tofu fmt -recursive infra/

# OpenTofu validation
terraform-validate:
	@echo "Validating OpenTofu configurations..."
	@cd infra/terraform/mgmt && tofu init -backend=false && tofu validate
	@cd infra/terraform/dev && tofu init -backend=false && tofu validate

# OpenTofu plan (requires AWS credentials and tfvars)
terraform-plan:
	@echo "Planning OpenTofu changes (dev)..."
	@cd infra/terraform/dev && tofu plan -var-file=dev.tfvars
