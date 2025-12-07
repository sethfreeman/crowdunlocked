# Crowd Unlocked - Quick Start Guide

## 🚀 What You've Got

A production-ready monorepo with:
- ✅ 5 Go microservices (TDD-ready with tests)
- ✅ Next.js web app
- ✅ Flutter mobile app
- ✅ Complete AWS infrastructure (Terraform)
- ✅ GitOps with Flux CD
- ✅ Multi-account AWS setup
- ✅ CI/CD pipeline
- ✅ Observability (X-Ray, CloudWatch)

## 📋 Prerequisites

Install these tools:
```bash
# Required
brew install go node terraform kubectl fluxcd/tap/flux awscli

# Optional (for mobile)
brew install --cask flutter
```

## 🏃 Quick Start (5 minutes)

### 1. Bootstrap the Project
```bash
./scripts/bootstrap.sh
```

### 2. Run Tests
```bash
make test
```

### 3. Start Local Environment
```bash
# Start all services
docker-compose up -d

# Wait for services to be ready (about 10 seconds)
sleep 10

# Create DynamoDB tables (requires AWS CLI with dummy credentials)
AWS_ACCESS_KEY_ID=test AWS_SECRET_ACCESS_KEY=test \
  bash scripts/create-dynamodb-tables.sh

# Verify tables were created
AWS_ACCESS_KEY_ID=test AWS_SECRET_ACCESS_KEY=test \
  aws dynamodb list-tables --endpoint-url http://localhost:8000 --region us-east-1

# Access services:
# - Bookings: http://localhost:8081
# - Releases: http://localhost:8082
# - Web App: http://localhost:3000
# - DynamoDB Local: http://localhost:8000
```

### 4. Test the Bookings API
```bash
# Create a booking
curl -X POST http://localhost:8081/api/v1/bookings \
  -H "Content-Type: application/json" \
  -d '{
    "artist_id": "artist-123",
    "venue_id": "venue-456",
    "event_date": "2024-12-15T20:00:00Z",
    "fee": 5000.0
  }'

# Get booking (use ID from response)
curl http://localhost:8081/api/v1/bookings/{id}

# Check service health
curl http://localhost:8081/health
```

## 🌩️ AWS Deployment

### Step 1: Configure AWS
```bash
aws configure
# Enter your AWS credentials
```

### Step 2: Deploy Mgmt Account
```bash
cd infra/terraform/mgmt

terraform init
terraform apply \
  -var="domain_name=crowdunlocked.com" \
  -var="dev_account_email=dev@example.com" \
  -var="prod_account_email=prod@example.com"
```

### Step 3: Deploy Dev Environment
```bash
cd ../dev
terraform init
terraform apply
```

### Step 4: Deploy Prod Environment
```bash
cd ../prod
terraform init
terraform apply -var="domain_name=crowdunlocked.com"
```

### Step 5: Setup GitOps

**For Dev Cluster:**
```bash
aws eks update-kubeconfig --name crowdunlocked-dev --region us-east-1

flux bootstrap github \
  --owner=YOUR_GITHUB_ORG \
  --repository=crowdunlocked \
  --branch=develop \
  --path=flux/clusters/dev \
  --personal
```

**For Prod Cluster:**
```bash
aws eks update-kubeconfig --name crowdunlocked-prod --region us-east-1

flux bootstrap github \
  --owner=YOUR_GITHUB_ORG \
  --repository=crowdunlocked \
  --branch=main \
  --path=flux/clusters/prod \
  --personal
```

## 🔄 Development Workflow

### Making Changes

1. **Create a feature branch**
   ```bash
   git checkout -b feature/new-booking-feature
   ```

2. **Write tests first (TDD)**
   ```bash
   cd services/bookings
   # Edit internal/domain/booking_test.go
   go test ./...  # Should fail
   ```

3. **Implement the feature**
   ```bash
   # Edit internal/domain/booking.go
   go test ./...  # Should pass
   ```

4. **Test locally**
   ```bash
   docker-compose up --build
   ```

5. **Push and deploy**
   ```bash
   git push origin feature/new-booking-feature
   # Create PR to develop
   # Merge → Auto-deploys to dev cluster
   ```

### Promoting to Production

```bash
# Merge develop to main
git checkout main
git merge develop
git push origin main
# Auto-deploys to prod cluster
```

## 📊 Monitoring

### View Logs
```bash
# Dev cluster
kubectl logs -f deployment/bookings -n default

# CloudWatch
aws logs tail /aws/eks/crowdunlocked-dev/bookings --follow
```

### View Traces
```bash
# Open AWS Console → X-Ray → Service Map
```

### Check Flux Status
```bash
flux get kustomizations
flux get sources git
```

## 🧪 Testing

### Unit Tests
```bash
# All services
make test

# Specific service
cd services/bookings
go test -v -cover ./...
```

### Integration Tests
```bash
# Start local environment
docker-compose up -d

# Create tables with dummy credentials
AWS_ACCESS_KEY_ID=test AWS_SECRET_ACCESS_KEY=test \
  bash scripts/create-dynamodb-tables.sh

# Run your integration tests
make test-integration

# Clean up
docker-compose down
```

## 📱 Mobile Development

```bash
cd apps/mobile

# Get dependencies
flutter pub get

# Run on device
flutter run

# Build for release
flutter build apk        # Android
flutter build ios        # iOS
```

## 🌐 Web Development

```bash
cd apps/web

# Install dependencies
npm install

# Run dev server
npm run dev

# Build for production
npm run build
npm start
```

## 🔧 Troubleshooting

### Services won't start
```bash
# Check Go modules
cd services/bookings
go mod tidy
go mod download

# Rebuild containers
docker-compose build
docker-compose up -d
```

### DynamoDB tables not created
```bash
# The script needs AWS CLI with dummy credentials for local DynamoDB
AWS_ACCESS_KEY_ID=test AWS_SECRET_ACCESS_KEY=test \
  bash scripts/create-dynamodb-tables.sh

# Verify tables exist
AWS_ACCESS_KEY_ID=test AWS_SECRET_ACCESS_KEY=test \
  aws dynamodb list-tables --endpoint-url http://localhost:8000 --region us-east-1
```

### Terraform errors
```bash
# Re-initialize
terraform init -upgrade

# Check state
terraform state list
```

### Flux not syncing
```bash
# Force reconciliation
flux reconcile kustomization flux-system --with-source

# Check logs
flux logs
```

## 📚 Documentation

- **Architecture**: See `docs/ARCHITECTURE.md`
- **Setup Guide**: See `docs/SETUP.md`
- **Project Structure**: See `PROJECT_STRUCTURE.md`
- **Service READMEs**: Each service has its own README

## 🎯 Next Steps

1. ✅ Run local environment
2. ✅ Deploy to AWS
3. ✅ Setup GitOps
4. 🔲 Add your domain to Route 53
5. 🔲 Configure SSO users
6. 🔲 Implement remaining service features
7. 🔲 Add monitoring dashboards
8. 🔲 Setup alerting

## 💡 Tips

- **TDD**: Always write tests first
- **GitOps**: Let Flux handle deployments
- **Observability**: Use X-Ray for debugging
- **Security**: Never commit secrets
- **Scaling**: EKS Auto Mode handles it

## 🆘 Need Help?

- Check service logs: `kubectl logs -f deployment/{service}`
- View X-Ray traces: AWS Console → X-Ray
- Check Flux status: `flux get all`
- Review CloudWatch alarms: AWS Console → CloudWatch

---

**Built with ❤️ for Crowd Unlocked**
