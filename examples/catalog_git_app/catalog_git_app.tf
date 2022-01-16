provider "nirmata" {
  #  Nirmata API Key. Best configured as the environment variable NIRMATA_TOKEN.
  
  #  token = ""

  #  Nirmata address. Defaults to https://nirmata.io and can be configured as
  #  the environment variable NIRMATA_URL.
  
  #  url = ""
}

resource "nirmata_git_application" "tf-catalog-git" {
  name                = "tf-catapp"
  catalog             = ""
  namespace           = ""
  git_credentials     = ""
  git_repository      = ""
  git_branch          =""
  git_directory_list  = ["*.yaml", "*.yml"]
  git_include_list    = []
  fixed_kustomization = true
  target_based_kustomization = true
  kustomization_file_path = ""
}
output "version" {
  value = nirmata_git_application.tf-catalog-git.version
}