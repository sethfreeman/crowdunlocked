#!/bin/bash
set -e

# Script to remove manually managed resources from terraform state
# Usage: ./scripts/cleanup-terraform-state.sh <environment>

ENVIRONMENT=$1

if [ -z "$ENVIRONMENT" ]; then
  echo "Usage: $0 <environment>"
  echo "Example: $0 mgmt"
  exit 1
fi

echo "Cleaning up terraform state for $ENVIRONMENT environment..."

cd "infra/terraform/$ENVIRONMENT"

# Remove AWS Organizations resources (managed manually)
echo "Removing AWS Organizations resources from state..."
tofu state rm aws_organizations_organization.main 2>/dev/null || echo "  - aws_organizations_organization.main not in state"
tofu state rm aws_organizations_account.dev 2>/dev/null || echo "  - aws_organizations_account.dev not in state"
tofu state rm aws_organizations_account.prod 2>/dev/null || echo "  - aws_organizations_account.prod not in state"

# Remove Route53 domains (managed manually)
echo "Removing Route53 domain resources from state..."
tofu state rm aws_route53domains_registered_domain.prod 2>/dev/null || echo "  - aws_route53domains_registered_domain.prod not in state"
tofu state rm aws_route53domains_registered_domain.dev 2>/dev/null || echo "  - aws_route53domains_registered_domain.dev not in state"

# Remove ACM certificates (managed manually)
echo "Removing ACM certificate resources from state..."
tofu state rm aws_acm_certificate.prod 2>/dev/null || echo "  - aws_acm_certificate.prod not in state"
tofu state rm aws_acm_certificate.dev 2>/dev/null || echo "  - aws_acm_certificate.dev not in state"
tofu state rm aws_acm_certificate_validation.prod 2>/dev/null || echo "  - aws_acm_certificate_validation.prod not in state"
tofu state rm aws_acm_certificate_validation.dev 2>/dev/null || echo "  - aws_acm_certificate_validation.dev not in state"

# Remove Route53 zones (managed manually)
echo "Removing Route53 zone resources from state..."
tofu state rm aws_route53_zone.prod 2>/dev/null || echo "  - aws_route53_zone.prod not in state"
tofu state rm aws_route53_zone.dev 2>/dev/null || echo "  - aws_route53_zone.dev not in state"

# Remove Route53 records for certificate validation (managed manually)
echo "Removing Route53 validation records from state..."
tofu state rm 'aws_route53_record.prod_cert_validation[0]' 2>/dev/null || echo "  - aws_route53_record.prod_cert_validation[0] not in state"
tofu state rm 'aws_route53_record.dev_cert_validation[0]' 2>/dev/null || echo "  - aws_route53_record.dev_cert_validation[0] not in state"

echo ""
echo "✅ Terraform state cleanup complete for $ENVIRONMENT!"
echo ""
