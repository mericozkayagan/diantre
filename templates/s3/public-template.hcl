include {
  path = find_in_parent_folders()
}

locals {
  common_vars = yamldecode(file(find_in_parent_folders("common_vars.yaml")))
  name        = "{{ Name }}"
}

terraform {
  source = "{{ Source }}"
}

inputs = {
  bucket        = "${local.common_vars.namespace}-${local.common_vars.environment}-${local.name}"
  description   = "S3 bucket for storing ${local.name} in ${local.common_vars.environment} environment"
  force_destroy = true
  attach_policy = true
  policy        = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "PublicReadGetObject",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "s3:GetObject",
      "Resource": "arn:aws:s3:::${local.common_vars.namespace}-${local.common_vars.environment}-${local.name}/*"
    }
  ]
}
EOF

  cors_rule = [
    {
      allowed_headers = ["Content-Type"]
      allowed_methods = ["PUT"]
      allowed_origins = ["{{ AllowedOrigins }}"]
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
