# Vercel Deployment Guide with AWS OIDC

This guide walks you through deploying the Crowd Unlocked API to Vercel with secure AWS OIDC authentication.

## Prerequisites

1. **AWS Infrastructure**: Ensure your dev environment is deployed with the Vercel OIDC module
2. **Vercel Account**: Sign up at [vercel.com](https://vercel.com)
3. **GitHub Repository**: Your code should be pushed to GitHub

## Step 1: Deploy AWS Infrastructure

First, deploy the Vercel OIDC infrastructure:

```bash
cd infra/terraform/dev
tofu plan
tofu apply
```

This will create:
- OIDC provider for Vercel
- IAM role with DynamoDB permissions
- Necessary policies for secure access

## Step 2: Deploy to Vercel

### 2.1 Create New Project

1. Go to [vercel.com/dashboard](https://vercel.com/dashboard)
2. Click "New Project"
3. Import your GitHub repository
4. **Important**: Set the root directory to `apps/web`
5. Vercel will auto-detect it's a Next.js project

### 2.2 Configure Build Settings

Vercel should automatically detect:
- **Framework Preset**: Next.js
- **Root Directory**: `apps/web`
- **Build Command**: `npm run build`
- **Output Directory**: `.next`

### 2.3 Get Vercel Project Information

After deployment, you'll need to get your Vercel project details:

1. Go to your project settings in Vercel
2. Note down your **Team ID** and **Project ID** from the URL or settings
3. The format will be: `team_<team_id>:project_<project_id>:environment_production`

## Step 3: Update Terraform with Vercel Project IDs

Update your `infra/terraform/dev/main.tf` file:

```hcl
module "vercel_oidc" {
  source = "../modules/vercel-oidc"

  project_name    = "crowdunlocked"
  environment     = "dev"
  aws_region      = var.aws_region
  aws_account_id  = local.aws_account_id

  # Update with your actual Vercel project IDs
  vercel_project_ids = [
    "team_YOUR_TEAM_ID:project_YOUR_PROJECT_ID:environment_production",
    "team_YOUR_TEAM_ID:project_YOUR_PROJECT_ID:environment_preview",
  ]

  venues_table_arn   = aws_dynamodb_table.venues.arn
  bookings_table_arn = aws_dynamodb_table.bookings.arn
}
```

Then redeploy:

```bash
cd infra/terraform/dev
tofu apply
```

## Step 4: Configure Vercel Environment Variables

In your Vercel project settings, add these environment variables:

### Required Environment Variables

```bash
# AWS Configuration
AWS_REGION=us-east-1
DYNAMODB_VENUES_TABLE=venues-dev
DYNAMODB_BOOKINGS_TABLE=bookings-dev

# AWS OIDC Role (get this from Terraform output)
AWS_ROLE_ARN=arn:aws:iam::YOUR_ACCOUNT_ID:role/crowdunlocked-dev-vercel-role
```

### Get the Role ARN

Run this command to get your role ARN:

```bash
cd infra/terraform/dev
tofu output vercel_role_arn
```

## Step 5: Configure Vercel for AWS OIDC

### 5.1 Enable OIDC in Vercel

1. Go to your Vercel project settings
2. Navigate to "Environment Variables"
3. Add the `AWS_ROLE_ARN` variable
4. Vercel will automatically use OIDC when this is set

### 5.2 Test the Deployment

After deployment, test your endpoints:

1. **Health Check**: `https://your-app.vercel.app/api/health`
2. **Create Venue**: 
   ```bash
   curl -X POST https://your-app.vercel.app/api/v1/venues \
     -H "Content-Type: application/json" \
     -d '{
       "name": "Test Venue",
       "location": {"latitude": 40.7128, "longitude": -74.0060},
       "address": {
         "street": "123 Main St",
         "city": "New York", 
         "state": "NY",
         "postal_code": "10001",
         "country": "US"
       },
       "venue_types": ["club"]
     }'
   ```

3. **Search Venues**:
   ```bash
   curl "https://your-app.vercel.app/api/v1/venues/search?venue_types=club&limit=5"
   ```

## Step 6: Monitoring and Troubleshooting

### Check Logs

1. **Vercel Logs**: Go to your project dashboard → Functions tab
2. **AWS CloudWatch**: Check `/vercel/crowdunlocked/*` log groups

### Common Issues

1. **403 Forbidden**: Check that your Vercel project IDs are correct in Terraform
2. **Role not found**: Ensure the IAM role was created and the ARN is correct
3. **DynamoDB access denied**: Verify the role has the correct policies attached

### Verify OIDC Setup

Check that the OIDC provider was created:

```bash
aws iam list-open-id-connect-providers
```

## Step 7: Production Deployment

For production:

1. Create a production environment in Terraform
2. Use the same OIDC setup but with production table names
3. Deploy to Vercel production environment
4. Use production domain and SSL certificate

## Security Best Practices

✅ **What we've implemented:**
- No long-lived AWS access keys
- Least privilege IAM policies
- Environment-specific resources
- Secure OIDC token exchange

✅ **Additional recommendations:**
- Enable Vercel's security headers
- Set up monitoring and alerting
- Regular security audits
- Rotate OIDC thumbprints if needed

## Cost Optimization

- **Vercel Pro**: ~$20/month
- **DynamoDB**: Pay-per-use (very cost-effective)
- **AWS OIDC**: Free
- **Total estimated**: $25-35/month

## Next Steps

1. Set up custom domain
2. Configure SSL certificate
3. Add monitoring and alerting
4. Set up staging environment
5. Implement CI/CD pipeline for automatic deployments

## Support

If you encounter issues:
1. Check Vercel function logs
2. Check AWS CloudWatch logs
3. Verify IAM role permissions
4. Test locally with AWS credentials first