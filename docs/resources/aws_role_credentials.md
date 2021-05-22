---
page_title: "nirmata_aws_role_credentials"
---

# nirmata_aws_role_credentials

A reusable configuration set the aws credential.

## Example Usage

```hcl

resource "nirmata_aws_role_credentials" "aws_role" {
  name                      = ""
  region                    = ""
  description               = ""
  aws_access_key_id         = ""
  aws_secret_key            = ""
  # aws_role_arn            = ""
  # aws_external_id         = ""
}

```

## Argument Reference

* `name` - (Required) a unique name for the credentials.
* `region` - (Required) use as the default credentials.
* `aws_access_key_id` - (Required) The AWS access key ID.
* `aws_secret_key` - (Required) The AWS secret access key.
* `aws_role_arn` - (Optional) The Amazon Resource Name (ARN) of the role to assume.
* `aws_external_id` - (Optional) A unique identifier that might be required when you assume a role in another account.
