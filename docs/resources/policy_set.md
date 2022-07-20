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

* `name` - (Required) Enter a unique name for the policy set.
* `is_default` - (Optional) This field indicates that the policy is set as default. The default policy set will be automatically deployed to new clusters.
* `git_credentials` - (Required) Enter the git credential name.
* `git_repository` - (Required)  Enter the repository URL as used in the git clone command.
* `git_branch` - (Required) Enter the Git branch to track. It indicates the name.
* `git_directory_list` - (Optional) This field indicates the directories to track.
* `fixed_kustomization` - (Optional)  This field enables fixed kustomize to select kustomizations for your application.
* `target_based_kustomization` - (Optional) This field enables target based kustomize to select kustomizations for your application.
* `kustomization_file_path` - (Required if fixed_kustomization or target_based_kustomization is set) Enter the kustomization file path. kustomization file path is required if you select fixed_kustomization or target_based_kustomization. 
* `delete_from_cluster` - (Optional) This field indicates the delete from cluster.


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
* `policy_set_name` - (Required) Enter deploy policy set name.
* `cluster` - (Required) Enter the cluster name in which the policy is set to deploy.
* `delete_from_cluster` - (Optional) This field indicates the delete from cluster.
