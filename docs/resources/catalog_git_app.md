---
page_title: "nirmata_git_application Resource"
---

# nirmata_git_application Resource

 Application is a group of workloads, routing and storage configurations.

## Example Usage

```hcl

resource "nirmata_git_application" "tf-catalog-git" {
  name                = "tf-git-app"
  catalog             = "test-catalog"
  namespace           = "test-namespace"
  git_credentials     = "test-creden"
  git_repository      = "repo_path"
  git_branch          = "main"
  git_directory_list  = ["/test"]
  git_include_list    = ["*.yaml", "*.yml"]
  fixed_kustomization = false
  target_based_kustomization = true
  kustomization_file_path = ""
}

output "version" {
  value = nirmata_git_application.tf-catalog-git.version
}

resource "nirmata_run_application" "tf-catalog-run-app" {
  name                = "tf-run-app"
  application         = "test_application"
  catalog             = "test-catalog"
  version             = nirmata_git_application.tf-catalog-git.version
  channel             = "Rapid"
  environments        = ["test-env", "test-envs"]
  depends_on          = [nirmata_git_application.tf-catalog-git]
 }

```

## Argument Reference

* `name` - (Required) a unique name for the application in catalog.
* `namespace` - (Optional) namespace for the git application.
* `git_credentials` - (Required) the git credential name.
* `git_repository` - (Required)  the repository URL as used in the git clone command.
* `git_branch` - (Required) the Git branch to track. name.
* `git_directory_list` - (Optional)  the directories to track.
* `git_include_list` - (Optional)  the file extensions to track.
* `fixed_kustomization` - (Optional)  enable fixed kustomize to select kustomizations for your application.
* `target_based_kustomization` - (Optional) enable target based kustomize to select kustomizations for your application.
* `kustomization_file_path` - (Required if fixed_kustomization or target_based_kustomization is set) the kustomization file path. kustomization file path required if fixed_kustomization or target_based_kustomization selected. 


* `name` - (Required) A unique name to identify your application.
* `catalog` - (Required) the name of catalog.
* `application` - (Required) the application name.
* `channel` - (Required) The channel from which the application should be deployed.
* `environments` - (Required) the list of environments to deploy an application .
* `version` - (Required)  the version for the application.