---
page_title: "nirmata_aws_role_credentials"
---

# nirmata_aws_role_credentials

A reusable configuration set the aws cloud credential.

## Example Usage

```hcl

resource "nirmata_aws_role_credentials" "aws_role" {
  name                      = "aws-credential"
  access_type               = "access_key"  // value are access_key or assume_role
  description               = "AWS Account"
  region                    = "us-west-1"
  access_key_id             = ""            
  secret_key                = ""            
  # aws_role_arn            = ""            
}

```

## Argument Reference

* `name` - (Required) a unique name for the credentials.
* `region` - (Required) use as the default credentials.
* `access_type` - (Required) select type for credentials ( access_key or assume_role).
* `aws_access_key_id` - (Optional) The AWS access key ID.Required if  access_type is access_key 
* `aws_secret_key` - (Optional) The AWS secret access key. Required if access_type is access_key 
* `aws_role_arn` - (Optional) The Amazon Resource Name (ARN) of the role to assume. Required if access_type is assume_role 

