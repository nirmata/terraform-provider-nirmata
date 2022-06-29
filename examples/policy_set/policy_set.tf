provider "nirmata" {
  #  Nirmata API Key. Best configured as the environment variable NIRMATA_TOKEN.
  
  #  token = ""

  #  Nirmata address. Defaults to https://nirmata.io and can be configured as
  #  the environment variable NIRMATA_URL.
  
  #  url = ""
}

resource "nirmata_policy_set" "create-policy-set" {
  name                        = "policy-set-name"
  is_default                  = false
  git_credentials             = "Dolis"
  git_repository              = "https://github.com/nirmata-add-ons/policies.git"
  git_branch                  = "test"
  git_directory_list          = ["/testpolicies"]
  fixed_kustomization         = true
  target_based_kustomization  = false
  kustomization_file_path     = "/testpolicies/kustomization.yaml"
  delete_from_cluster         = true
}

resource "nirmata_deploy_policy_set" "tf-policy-set-deploy" {
  policy_set_name                = "policy-set-name"
  cluster                         = "test-csi-3"
  delete_from_cluster             = true
  depends_on                      = [nirmata_policy_set.create-policy-set]
}
