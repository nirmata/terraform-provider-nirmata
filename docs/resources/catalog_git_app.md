---
page_title: "nirmata_git_application Resource"
---

# nirmata_git_application Resource

 Application is a group of workloads, routing, and storage configurations.

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

resource "nirmata_promote_version" "tf-catalog-promote-version" {
  rollout_name        = "tf-version"
  catalog             = "test-catalog"
  application         = "test-application"
  version             = nirmata_git_application.tf-catalog-git.version
  channel             = "Rapid"
  depends_on          = [nirmata_git_application.tf-catalog-git]
 }
```

## Argument Reference

* `name` - (Required) A unique name for the application in the catalog.
* `namespace` - (Optional) This field indicates the namespace for the git application.
* `git_credentials` - (Required) This field indicates the git credentials name.
* `git_repository` - (Required)  This is the repository URL as used in the git clone command.
* `git_branch` - (Required) Enter the git branch name to track.
* `git_directory_list` - (Optional)  This field indicates the git directories to track.
* `git_include_list` - (Optional)  the file extensions to track.
* `fixed_kustomization` - (Optional)  enable fixed kustomize to select kustomizations for your application.
* `target_based_kustomization` - (Optional) enable target based kustomize to select kustomizations for your application.
* `kustomization_file_path` - (Required if fixed_kustomization or target_based_kustomization is set) the kustomization file path. kustomization file path required if fixed_kustomization or target_based_kustomization selected. 


* `name` - (Required) Enter a unique name to identify your application.
* `catalog` - (Required) Enter the name of the catalog.
* `application` - (Required) Enter the application name.
* `channel` - (Required) Enter the channel from which the application should be deployed.
* `environments` - (Required) Enter the list of environments to deploy an application.
* `version` - (Required) Enter the version for the application.

* `rollout_name` - (Required) Enter a unique name for rollout.
* `catalog` - (Required) Enter the name of the catalog.
* `application` - (Required) Enter the application name.
* `channel` - (Required) Enter the channel from which the application should be deployed.
* `version` - (Required) Enter the version of the application.