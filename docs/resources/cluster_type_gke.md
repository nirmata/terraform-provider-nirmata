---
page_title: "nirmata_cluster_type_gke Resource"
---

# nirmata_cluster_type_gke Resource

A Google Kubernetes Engine (GKE) cluster type. A cluster type can be used to create multiple clusters.

## Example Usage

Create a GKE cluster type with a Vault Agent Injector add-on:

```hcl
resource "nirmata_cluster_type_gke" "gke-us-west" {

  name = "gke-us-west"
  version = ""1.17.13-gke.2001"
  credentials = "gcp"
  location_type =  "Zonal"
  region = "us-central1"
  zone = "us-central1-a"
  network = "default"
  subnetwork = "default"

  nodepools { 
    machine_type = "c2-standard-16"
    disk_size= 120
    enable_preemptible_nodes  =  true
    node_annotations = {
       node = "annotate"
    }
  }

  addons {
    name            = "vault-agent-injector"
    addon_selector  = "vault-agent-injector"
    catalog         = "default-catalog"
    channel         = "Stable"
    sequence_number = 15
  }

  vault_auth {
    name             = "gke-vault-auth"
    path             = "nirmata/$(cluster.name)"
    addon_name       = "vault-agent-injector"
    credentials_name = "vault_access"
     delete_auth_path = true

    roles {
      name                 = "sample-role"
      service_account_name = "application-sample-sa"
      namespace            = "application-sample-ns"
      policies             = "application-sample-policy"
    }
  }
}
```

## Argument Reference

* `name` - (Required) a unique name for the cluster.
* `version` - (Required) the GKE version (e.g. 1.16.12-gke.3)
* `credentials` - (Required) the cloud credentials to use.
* `location_type` - (Required) Regional or Zonal. A regional cluster has multiple replicas of the control plane, running in multiple zones within a given region. A zonal cluster has a single replica of the control plane running in a single zone.
* `region` - (Optional) the GCP region into which the cluster should be deployed (e.g. "us-central1"). Required when location_type is `Regional`.
* `zone` - (Optional) the GCP zone into which the cluster should be deployed (e.g. "us-central1-a"). Required if location_type is `Zonal`.
* `node_locations` - (Optional) the nodes should be deployed. Selecting more than one zone increases availability. (e.g. ["asia-east1-a"]). Required if location_type is `Regional`.
* `enable_secrets_encryption` - (Optional) enables envelope encryption for Kubernetes Secrets.
* `enable_workload_identity` - (Optional) Workload Identity is the recommended way to access Google Cloud services from applications running within GKE due to its improved security properties and manageability.
* `workload_pool` - (Optional) Workload Identity relies on a Workload Pool to aggregate identity across multiple clusters. Required if enable_secrets_encryption is true
* `secrets_encryption_key` - (Optional) the Resource ID of the key you want to use (e.g. projects/project-name/locations/global/keyRings/my-keyring/cryptoKeys/my-key). Required if enable_workload_identity is true.
* `network` - (Required) the cluster network (e.g. "default")
* `subnetwork` - (Required) the node subnetwork (e.g. "default")
* `cluster_ipv4_cidr` - (Optional) All pods in the cluster are assigned an IP address from this range. Enter a range (in CIDR notation) within a network range, a mask, or leave this field blank to use a default range. This setting is permanent.
* `services_ipv4_cidr` - (Optional) Cluster services will be assigned an IP address from this IP address range. Enter a range (in CIDR notation) within a network range, a mask, or leave this field blank to use a default range. This setting is permanent.
* `cloud_run` - (Optional) Cloud Run for Anthos enables you to easily deploy stateless apps and functions to this cluster using the Cloud Run experience. Cloud Run for Anthos automatically manages underlying resources and scales your app based on requests.
* `enable_network_policy` - (Optional) The Kubernetes Network Policy API allows the cluster administrator to specify what pods are allowed to communicate with each other. Google Kubernetes Engine has partnered with Tigera to provide Project Calico  to enforce network policies within your cluster.
* `enable_http_load_balancing` - (Optional) The HTTP Load Balancing add-on is required to use the Google Cloud Load Balancer with Kubernetes Ingress. If enabled, a controller will be installed to coordinate applying load balancing configuration changes to your GCP project
* `enable_vertical_pod_autoscaling` - (Optional) Vertical Pod Autoscaling automatically analyzes and adjusts your containers' CPU requests and memory requests.
* `enable_horizontal_pod_autoscaling` - 
* `enable_maintenance_policy` - (Optional) To specify regular times for maintenance, enable maintenance windows. Normally, routine Kubernetes Engine maintenance may run at any time on your cluster.
* `maintenance_start_time` - (Optional) Start time for the maintenance window.
* `maintenance_duration` - (Optional) Duration for the maintenance window in hours.
* `maintenance_recurrence` -  (Optional) Recurrence rule specification (RRULE) for the maintenance window. Example RRule to run maintenance during weekends: 'FREQ=WEEKLY;BYDAY=SA,SU'.
* `maintenance_exclusion_timewindow` - (Optional) To specify times when routine, non-emergency maintenance won't happen, set up to 3 maintenance exclusions. Normally, routine Kubernetes Engine maintenance may run at any time on your cluster.
* `system_metadata` - (Optional) key-value pairs that will be provisioned as a ConfigMap called system-metadata-map in the cluster.
* `allow_override_credentials` - (Optional) Allow passing in cloud credentials when a cluster is created using this cluster type.
* `cluster_field_override` - (Optional) Allow override of cluster settings ('network' and 'subnetwork') when a cluster is created using this cluster type.
* `nodepool_field_override` - (Optional)  Allow override of node fields ('disk_size' and 'machine_type') when a cluster is created using this cluster type.
* `nodepools` - (Optional) A list of [nodepool](#nodepool) types.
* `addons` - (Optional) a list of add-on services.
* `vault_auth` - (Optional) vault authentication configuration.

## Nested Blocks

### nodepool

* `machine_type` - (Required) the GCP machine type (e.g. "e2-standard-2")
* `disk_size` - (Required) the worker node disk size in GB (e.g. 60)
* `service_account` - (Optional) Applications running on the VM use the service account to call Google Cloud APIs. Use Permissions on the console menu to create a service account or use the default service account if available. Service account is permanent
* `enable_preemptible_nodes` - (Optional) Preemptible nodes are Compute Engine instances that last up to 24 hours and provide no availability guarantees, but are priced lower. This setting is permanent.
* `node_annotations` -  (Optional) Annotations to set on each node in this pool. This setting is permanent.
* `node_labels` - (Optional) Labels to set on each Kubernetes node in this node pool. This setting is permanent.

### addons

* `name` - (Required) a unique name for the add-on service
* `addon_selector` - (Required) the catalog application name
* `catalog` - (Required) the catalog name
* `channel` - (Required) the release channel
* `sequence_number` - (Optional) a sequence number to control installation order

### vault_auth

* `name` - (Required) a unique name
* `path` - (Required) the vault authentication path. The variable $(cluster.name) is allowed in the path for uniquenes.
* `addon_name` - (Required) the associated Vault Agent Injector add-on
* `credentials_name` - (Required) the Vault credentials to use 
* `roles` - (Required) a list of application roles to configure for add-on services
* `delete_auth_path` - (Optional) delete auth path on cluster delete

#### roles

* `name` - (Required) a unique name
* `service_account_name` - (Required) the allowed service account name
* `namespace` - (Required) the allowed namespace
* `policies` - (Required) the applied policies
