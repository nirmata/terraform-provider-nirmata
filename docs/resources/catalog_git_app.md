---
page_title: "nirmata_git_application Resource"
---

# nirmata_git_application Resource

 Application is a group of workloads, routing and storage configurations.

## Example Usage

```hcl

resource "nirmata_git_application" "tf-catalog-git-" {
  name                = "tf-git-app"
  catalog             = ""
  git_credentials     = ""
  git_repository      = ""
  git_branch          =""
  git_directory_list  = ["*.yaml", "*.yml"]
  git_include_list    = []
}

```

## Argument Reference

* `name` - (Required) a unique name for the application in catalog.
* `git_credentials` - (Required) the git credential name.
* `git_repository` - (Required)  the repository URL as used in the git clone command.
* `git_branch` - (Required) the Git branch to track. name.
* `git_directory_list` - (Optional)  the directories to track.
* `git_include_list` - (Optional)  the file extensions to track.