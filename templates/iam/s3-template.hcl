include {
  path = find_in_parent_folders()
}

locals {
  common_vars = yamldecode(file(find_in_parent_folders("common_vars.yaml")))
  name        = "{{ Name }} "
}

terraform {
  source = "{{ Source }}"
}

inputs = {
  name        = "${local.common_vars.namespace}-${local.common_vars.environment}-${local.name}-access"
  path        = "/"
  description = "IAM policy to allow read and write to {{ Name }} Bucket in ${local.common_vars.environment} environment."

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
      {
          "Effect": "Allow",
          "Action": [
              "s3:PutObject",
              "s3:GetObject",
              "s3:DeleteObject"
          ],
          "Resource": "${dependency.{{ Name }}-bucket.outputs.s3_bucket_arn}/*"
      }
  ]
}
EOF

  tags = local.common_vars.tags
}

dependency "{{ Name }}-bucket" {
  config_path = "../../../s3/{{ Name }}"
}

