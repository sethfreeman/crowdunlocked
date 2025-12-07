# Management Account Configuration
# Non-sensitive values that can be committed to git

aws_region = "us-east-1"

# Domain names - update these to your actual domains
prod_domain_name = "crowdunlocked.com"
dev_domain_name  = "crowdunlockedbeta.com"

# Contact type for domain registration
domain_contact_type    = "COMPANY"
domain_contact_country = "US"

# Sensitive values (emails, personal info) should be provided via:
# - GitHub Secrets as TF_VAR_* environment variables in CI/CD
# - Local environment variables for manual runs
# - Or local terraform.tfvars (gitignored) for development
#
# Required environment variables:
# - TF_VAR_dev_account_email
# - TF_VAR_prod_account_email
# - TF_VAR_domain_contact_email
# - TF_VAR_domain_contact_first_name
# - TF_VAR_domain_contact_last_name
# - TF_VAR_domain_contact_phone
# - TF_VAR_domain_contact_address
# - TF_VAR_domain_contact_city
# - TF_VAR_domain_contact_state
# - TF_VAR_domain_contact_zip
