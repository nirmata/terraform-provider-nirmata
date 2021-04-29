# Terraform Provider for Nirmata

<h3> See documrntation and examples at the [Terraform Registry](https://registry.terraform.io/providers/nirmata/nirmata/latest) </h3>


## Releasing

To release a new version create a tag (using semantic versioning) and push upstream. The release process is completed via a [GitHub Action](.github/workflows/release.yml)

```bash
git tag -a v1.x.x -m "...."
git push --tags
```

## Building

```bash
make
```

## Testing locally

1. Configure `dev_overrides` in your `.terraform.rc` or `terraform.rc` file. See: https://www.terraform.io/docs/cli/config/config-file.html#development-overrides-for-provider-developers.

```hcl
provider_installation {
    dev_overrides {
        "registry.terraform.io/nirmata/nirmata" = "<repo path>/dist/<platform>_<architecture>"
    }

    # For all other providers, install them directly from their origin provider
    # registries as normal. If you omit this, Terraform will _only_ use
    # the dev_overrides block, and so no other providers will be available.
    direct {}
}
```

For example on Windows the `nirmata/nirmata` provider would be set to `"C:\\go\\src\\github.com\\nirmata\\terraform-provider-nirmata\\dist\\windows_amd64"`

2. Build the plugin using `make`.

3. Set your `NIRMATA_TOKEN` environment variable to contain your Nirmata API key. You can optionally set `NIRMATA_URL` to point to the Nirmata address (defaults to https://nirmata.io.)

4. Navigate to the examples and initialize the Terraform provider:

```bash
terraform init 
```

If you see the error below, delete the `.terraform.lock.hcl` file and re-run the `init` command:

```bash
Error while installing local/nirmata/nirmata v99.0.0: the local package for
local/nirmata/nirmata 99.0.0 doesn't match any of the checksums previously
recorded in the dependency lock file (this might be because the available
checksums are for packages targeting different platforms)
```

6. Run `plan` to build the execution plan:

```bash
terraform plan
```

7. Run `apply` to execute the plan:

```bash
terraform apply
```

8. Run `show` to see the created resources:

```bash
terraform show
```

9. Run `destroy` to delete the created resources:

````bash
terraform destroy samples/cloud_provider/gke
````

## Troubleshooting

Set the TF_LOG environment variable to `DEBUG` or `TRACE`.

```bash
export TF_LOG=DEBUG
```
