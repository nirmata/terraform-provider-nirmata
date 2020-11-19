package nirmata

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nirmata/go-client/pkg/client"
	"log"
	"time"
)

var importedClusterSchema = map[string]*schema.Schema{
	"name": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validateName,
	},
	"credentials": {
		Type:     schema.TypeString,
		Required: true,
	},
	"region": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"cluster_type": {
		Type:     schema.TypeString,
		Required: true,
	},
}

func resourceClusterImported() *schema.Resource {
	return &schema.Resource{
		Create: resourceClusterImportedCreate,
		Read:   resourceClusterRead,
		Update: resourceClusterUpdate,
		Delete: resourceClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: importedClusterSchema,
	}
}

func resourceClusterImportedCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	credentials := d.Get("credentials").(string)
	region := d.Get("region").(string)
	clusterType := d.Get("cluster_type").(string)

	cloudCredID, err := apiClient.QueryByName(client.ServiceClusters, "CloudCredentials", credentials)
	if err != nil {
		log.Printf("[ERROR] - %v", err)
		return err
	}

	data := map[string]interface{}{
		"mode":                "providerManaged",
		"clusterTypeSelector": clusterType,
		"credentialsRef":      cloudCredID.UUID(),
		"clusters": map[string]interface{}{
			name: map[string]interface{}{
				"name":   name,
				"region": region,
			},
		},
	}

	log.Printf("[DEBUG] - importing cluster %s with %+v", name, data)

	_, err = apiClient.PostFromJSON(client.ServiceClusters, "ImportClustersAction", data, nil)
	if err != nil {
		log.Printf("[ERROR] - failed to import cluster %s with data %v: %v", name, data, err)
		return err
	}

	//state, waitErr := waitForImportClustersAction(apiClient, d.Timeout(schema.TimeoutCreate), result)
	//if waitErr != nil {
	//	log.Printf("[ERROR] - failed import cluster. Error - %v", waitErr)
	//	return nil
	//}

	return nil
}
