# Management Account Configuration
# NOTE: Organization, accounts, domains, and certificates must be created manually
# See docs/AWS_ORGANIZATION_SETUP.md for instructions

aws_region = "us-east-1"

# Domain names (must match manually registered domains)
prod_domain_name = "crowdunlocked.com"
dev_domain_name  = "crowdunlockedbeta.com"

# GitHub repository for OIDC
github_org  = "sethfreeman"
github_repo = "crowdunlocked"

# AWS Account IDs (from manual account creation)
dev_account_id  = "987470856210"
prod_account_id = "379211248770"


# Trigger workflow run
