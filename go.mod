module github.com/nirmata/terraform-provider-nirmata

go 1.14

require (
	github.com/google/uuid v1.1.1
	github.com/hashicorp/terraform-plugin-sdk v1.9.0
	github.com/nirmata/go-client v1.0.0
)

replace github.com/nirmata/go-client v1.0.0 => github.com/evalsocket/go-client v1.0.2-0.20200922070310-ee3a7873676c
