---
layout: "nirmata"
page_title: "Provider: Nirmata"
description: |-
  
---

# Nirmata Provider


## Example Usage

```hcl
# Configure the Nirmata Provider
provider "nirmata" {
  version = "~> 0.1"
  token   = ""
  url     = ""
}

# Create a nirmata host group 
resource "nirmata_host_group_direct_connect" "dc-host-group" {
  name = "dc-hg-1"
}
```


### Authentication

!> **Warning:** Hard-coding credentials into any Terraform configuration is not
recommended, and risks secret leakage should this file ever be committed to a
public version control system.

Static credentials can be provided by adding an `token` and `url`
in-line in the Nirmata provider block:

Usage:

```hcl
provider "nirmata" {
  token   = ""
  url     = ""
}
```

### Environment Variables

You can provide your credentials via the `NIRMATA_TOKEN` and
`NIRMATA_URL`, environment variables, representing your Nirmata
Access Key and URL, respectively. 

```hcl
provider "nirmata" {}
```

Usage:

```sh
$ export NIRMATA_TOKEN=""
$ export NIRMATA_URL="https://nirmata.io"
$ terraform plan
```


Please note that the [Nirmata Go SDK](https://github.com/nirmata/client-go), the underlying authentication handler used by the Terraform Nirmata Provider

## Argument Reference

In addition to [generic `provider` arguments](https://www.terraform.io/docs/configuration/providers.html)
(e.g. `alias` and `version`), the following arguments are supported in the nirmata
 `provider` block:

* `token` - (Required) 

* `url` - (Required) 


