provider "nirmata" {
  #  Nirmata API Key. Best configured as the environment variable NIRMATA_TOKEN.
  
    token = ""

  #  Nirmata address. Defaults to https://nirmata.io and can be configured as
  #  the environment variable NIRMATA_URL.
  
  #  url = ""
}

#  An Environment Resource Type is used while run application

resource "nirmata_git_application" "tf-catalog-git-" {
  name                = "tf-catapp"
  catalog             = ""
  git_credentials     = ""
  git_repository      = ""
  git_branch          =""
  git_directory_list  = ["*.yaml", "*.yml"]
  git_include_list    = []
}
