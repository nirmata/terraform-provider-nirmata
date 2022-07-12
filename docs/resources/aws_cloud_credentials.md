---
page_title: "nirmata_aws_role_credentials"
---

# nirmata_aws_role_credentials

A reusable configuration set the aws cloud credential.

## Example Usage

```hcl

resource "nirmata_aws_cloud_credentials" "aws_cloud_credential" {
  name                      = "aws-credential"
  access_type               = "access_key"  // The value should be either access_key or assume_role.
  description               = "AWS Account"
  region                    = "us-west-1"
  access_key_id             = ""            
  secret_key                = ""            
  # aws_role_arn            = ""            
}

```

## Argument Reference

* `name` - (Required) A unique name for the AWS cloud credentials.
* `region` - (Required) Use "default credentials" as the value for this field.
* `access_type` - (Required) Select the access type for the AWS credentials (it is either access_key or assume_role).
* `aws_access_key_id` - (Optional) The AWS access key ID. This value is required if the access_type is access_key. 
* `aws_secret_key` - (Optional) Enter the AWS secret access key. This value is required if the access_type is access_key. 
* `aws_role_arn` - (Optional) The Amazon Resource Name (ARN) is the unique identifier of AWS resources. It is the role that is assumed for access type. This value is required if the access_type is assume_role. 
