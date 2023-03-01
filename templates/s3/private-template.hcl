include {
  path = find_in_parent_folders()
}

locals {
  common_vars     = yamldecode(file(find_in_parent_folders("common_vars.yaml")))
  name            = "{{ Name }}"
  allowed_origins = ["{{ AllowedOrigins }}"]
}

terraform {
  source = "{{ Source }}"
}

inputs = {
  bucket               = "${local.common_vars.namespace}-${local.common_vars.environment}-${local.name}"
  description          = "S3 bucket for storing ${local.name} in ${local.common_vars.environment} Environment"
  force_destroy        = true
  attach_public_policy = false
  block_public_acls    = true
  block_public_policy  = true

  cors_rule = [
    {
      allowed_headers = ["Content-Type"]
      allowed_methods = ["GET", "PUT"]
      allowed_origins = "${local.allowed_origins}"
    }
  ]

  server_side_encryption_configuration = {
    rule = {
      apply_server_side_encryption_by_default = {
        sse_algorithm = "AES256"
      }
    }
  }

  versioning = {
    enabled = true
  }

  tags = local.common_vars.tags
}
