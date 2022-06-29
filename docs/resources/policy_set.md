---
page_title: "nirmata_policy_set Resource"
---

# nirmata_policy_set Resource

 Policy sets are group of Kyverno policies for compliance and governance.

## Example Usage

```hcl

resource "nirmata_policy_set" "create-policy-set" {
  name                        = "policy-set-name"
  is_default                  = false
  git_credentials             = ""
  git_repository              = "https://github.com/nirmata-add-ons/policies.git"
  git_branch                  = "test"
  git_directory_list          = ["/testpolicies"]
  fixed_kustomization         = true
  target_based_kustomization  = false
  kustomization_file_path     = "/testpolicies/kustomization.yaml"
  delete_from_cluster         = true
}

```

## Argument Reference

* `name` - (Required) a unique name for the policy set.
* `is_default` - (Optional) set this policy set as default. Default policy sets will be automatically deployed to new clusters..
* `git_credentials` - (Required) the git credential name.
* `git_repository` - (Required)  the repository URL as used in the git clone command.
* `git_branch` - (Required) the Git branch to track. name.
* `git_directory_list` - (Optional)  the directories to track.
* `fixed_kustomization` - (Optional)  enable fixed kustomize to select kustomizations for your application.
* `target_based_kustomization` - (Optional) enable target based kustomize to select kustomizations for your application.
* `kustomization_file_path` - (Required if fixed_kustomization or target_based_kustomization is set) the kustomization file path. kustomization file path required if fixed_kustomization or target_based_kustomization selected. 
* `delete_from_cluster` - (Optional) should be delete from cluster.


# nirmata_deploy_policy_set Resource

## Example Usage

```hcl

resource "nirmata_deploy_policy_set" "tf-policy-set-deploy" {
  policy_set_name                = "policy-set-name"
  cluster                         = "cluster-name"
  delete_from_cluster             = true
  depends_on                      = [nirmata_policy_set.create-policy-set]
}

```

## Argument Reference
* `policy_set_name` - (Required) deploy policy set name.
* `cluster` - (Required) cluster name in which policy set to be deploy.
* `delete_from_cluster` - (Optional) should be delete from cluster.
