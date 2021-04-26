provider "kubectl" {
  load_config_file       = false
  host                   = "https://${data.google_container_cluster.my_cluster.endpoint}"
  token                  = data.google_client_config.provider.access_token
  cluster_ca_certificate =  base64decode(data.google_container_cluster.my_cluster.master_auth.0.cluster_ca_certificate)
}
// for split file and pass
data "kubectl_filename_list" "manifests" {
    pattern = "${nirmata_cluster_registered.gke-register.yaml_file}/*"
}
resource "kubectl_manifest" "test" {
    count = 7
    yaml_body = file(element(data.kubectl_filename_list.manifests.matches, count.index))
} 
