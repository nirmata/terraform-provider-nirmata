# Terraform Provider for Nirmata

## Building

````bash
go build
````

## Executing the samples

The samples are available in the [samples](samples) folder.

To run the samples, initialize the Terraform provider and then run the `plan` and `apply` commands. Here is an example of how to run the GKE cluster provisioning:

1. Build the provider

````bash
go build
````

2. Set your `NIRMATA_TOKEN` environment variable to contain your API key. You can optionally set `NIRMATA_URL` to point to the Nirmata address (defaults to https://nirmata.io.)

3. Edit the terraform config file samples/clustertypes/gke/gke.tf and include your credentials, and desired region, machinetype, and disksize.  (You can replace with the desired eks, or aks example.)

4. Initialize the Terraform provider with the correct directory

```bash
terraform init samples/clustertypes/gke
````

5. Run `plan` to build the execution plan:

````bash
terraform plan samples/clustertypes/gke
````

6. Run `apply` to execute the plan:

````bash
terraform apply samples/clustertypes/gke
````

7. Run `show` to see the created resources:

````bash
terraform show
````

8. Run `destroy` to delete the created resources:

````bash
terraform destroy samples/clustertypes/gke
````
