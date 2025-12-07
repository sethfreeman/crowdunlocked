# Crowd Unlocked Setup Guide

## Prerequisites

- AWS CLI configured with appropriate credentials
- Terraform >= 1.6
- kubectl
- Flux CLI
- Docker
- Go 1.22+
- Node.js 20+
- Flutter SDK

## Infrastructure Setup

### 1. Mgmt Account

```bash
cd infra/terraform/mgmt
terraform init
terraform plan -var="domain_name=crowdunlocked.com" \
  -var="dev_account_email=dev@crowdunlocked.com" \
  -var="prod_account_email=prod@crowdunlocked.com"
terraform apply
```

### 2. Dev Environment

```bash
cd infra/terraform/dev
terraform init
terraform apply
```

### 3. Production Environment

```bash
cd infra/terraform/prod
terraform init
terraform plan -var="domain_name=crowdunlocked.com"
terraform apply
```

## Kubernetes Setup

### Install Flux on Dev Cluster

```bash
aws eks update-kubeconfig --name crowdunlocked-dev --region us-east-1

flux bootstrap github \
  --owner=crowdunlocked \
  --repository=crowdunlocked \
  --branch=develop \
  --path=flux/clusters/dev \
  --personal
```

### Install Flux on Prod Cluster

```bash
aws eks update-kubeconfig --name crowdunlocked-prod --region us-east-1

flux bootstrap github \
  --owner=crowdunlocked \
  --repository=crowdunlocked \
  --branch=main \
  --path=flux/clusters/prod \
  --personal
```

## Local Development

### Run Services Locally

```bash
# Terminal 1 - Bookings
cd services/bookings
go run cmd/server/main.go

# Terminal 2 - Releases
cd services/releases
go run cmd/server/main.go
```

### Run Web App

```bash
cd apps/web
npm install
npm run dev
```

### Run Mobile App

```bash
cd apps/mobile
flutter pub get
flutter run
```

## Testing

```bash
# Run all tests
make test

# Test specific service
cd services/bookings
go test -v ./...
```

## Deployment

Deployments are automated via GitOps:

- Push to `develop` branch → Auto-deploy to dev cluster
- Push to `main` branch → Auto-deploy to prod cluster

## Monitoring

- CloudWatch Logs: `/aws/eks/crowdunlocked-{env}/{service}`
- X-Ray Traces: AWS Console → X-Ray
- CloudWatch Alarms: Configured for error rates
