# Terraform Provider for Nirmata

## Building

````bash
go build 
````

## Executing the samples

The samples are avaiilable in the [samples](samples) folder. 

To run the samples, initialize the Terraform provider and then run the `plan` and `apply` commands. Here is an example of how to run the GKE cluster provisioning:

1. Build the provider 

````bash
go build
````

2. Set your `NIRMATA_API` environment variable to contain your API key. You can optionally set `NIRMATA_URL` to point to the Nirmata address (defualts to https://nirmata.io.) 

3. Initialize the Terraform provider with the correct directory

```bash
terraform init samples/gke
````

4. Run `plan` to build the execution plan:

````bash
terraform init samples/gke
````

5. Run `apply` to execute the plan:

````bash
terraform init samples/gke
````

6. Run `show` to see the created resources:

````bash
terraform show samples/gke
````

7. Run `destroy` to delete the created resources:

````bash
terraform destroy samples/gke
````
