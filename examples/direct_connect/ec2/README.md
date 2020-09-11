## Requirements

No requirements.

## Providers

| Name | Version |
|------|---------|
| aws | n/a |
| nirmata | n/a |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| awsprops | n/a | `map` | <pre>{<br>  "ami": "ami-03ba3948f6c37a4b0",<br>  "instance_count": 3,<br>  "itype": "t3a.medium",<br>  "keyname": "terraform-test-west-1",<br>  "publicip": true,<br>  "region": "us-west-1",<br>  "secgroupname": "terraform-test",<br>  "subnet": "subnet-12345678909876543",<br>  "vpc": "vpc-00012345678909876"<br>}</pre> | no |

## Outputs

| Name | Description |
|------|-------------|
| agent\_script | Nirmata agent install command |
| id | List of IDs of instances |
| instance\_count | Number of instances to launch specified as argument to this module |
| public\_ip | List of public IP addresses assigned to the instances, if applicable |

