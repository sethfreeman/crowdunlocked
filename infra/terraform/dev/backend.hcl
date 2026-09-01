# Backend config for `tofu init -backend-config=backend.hcl`.
#
# The terraform state bucket lives in the AWS account behind the `default`
# CLI profile (021645491430), NOT the dev account. The dev *resources* live in
# the dev account (179151668767) via the `crowdunlocked-dev` profile.
#
# So the backend (state) and the provider (resources) use different profiles.
profile = "default"
