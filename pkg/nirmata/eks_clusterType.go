package nirmata

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	client "github.com/nirmata/go-client/pkg/client"
)

func resourceEksClusterType() *schema.Resource {
	return &schema.Resource{
		Create: resourceEksClusterTypeCreate,
		Read:   resourceEksClusterTypeRead,
		Update: resourceEksClusterTypeUpdate,
		Delete: resourceEksClusterTypeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if len(value) > 64 {
						errors = append(errors, fmt.Errorf(
							"%q cannot be longer than 64 characters", k))
					}
					if !regexp.MustCompile(`^[\w+=,.@-]*$`).MatchString(value) {
						errors = append(errors, fmt.Errorf(
							"%q must match [\\w+=,.@-]", k))
					}
					return
				},
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"credentials": {
				Type:     schema.TypeString,
				Required: true,
			},

			"region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpcid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subneid": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"clusterrolearn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"securitygroups": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"keyname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instancetypes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"disksize": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					if v.(int) < 9 {
						errors = append(errors, fmt.Errorf(
							"%q The disk size must be grater than 9", k))
					}
					return
				},
			},
		},
	}
}

func resourceEksClusterTypeCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)

	clouduuid, _ := newUUID()
	nodepooluuid, _ := newUUID()

	name := d.Get("name").(string)
	version := d.Get("version").(string)
	credentials := d.Get("credentials").(string)
	region := d.Get("region").(string)
	diskSize := d.Get("disksize").(int)
	instancetype := d.Get("instancetypes")
	keyname := d.Get("keyname").(string)
	securitygroups := d.Get("securitygroups")
	clusterrolearn := d.Get("clusterrolearn").(string)
	vpcid := d.Get("vpcid").(string)
	subneid := d.Get("subneid")

	cloudCredID, err := apiClient.QueryByName(client.ServiceClusters, "CloudCredentials", credentials)
	fmt.Printf("Error - %v", cloudCredID)
	if err != nil {
		fmt.Printf("Error - %v", err)
		return err
	}

	clustertype := map[string]interface{}{
		"name":        name,
		"description": "",
		"modelIndex":  "ClusterType",
		"spec": map[string]interface{}{
			"clusterMode": "providerManaged",
			"modelIndex":  "ClusterSpec",
			"version":     version,
			"cloud":       "aws",
			"addons": map[string]interface{}{
				"dns":        false,
				"modelIndex": "AddOns",
				"addons": map[string]interface{}{
					"name":          "kyverno",
					"addOnSelector": "kyverno",
					"catalog":       "default-addon-catalog",
				},
			},
			"cloudConfigSpec": map[string]interface{}{
				"credentials":   cloudCredID.UUID(),
				"id":            clouduuid,
				"modelIndex":    "CloudConfigSpec",
				"nodePoolTypes": nodepooluuid,
				"eksConfig": map[string]interface{}{
					"region":                region,
					"vpcId":                 vpcid,
					"subnetId":              subneid,
					"privateEndpointAccess": false,
					"clusterRoleArn":        clusterrolearn,
				},
			},
		},
	}

	nodepoolobj := map[string]interface{}{
		"id":              nodepooluuid,
		"modelIndex":      "NodePoolType",
		"name":            name + "-default-node-pool-type",
		"cloudConfigSpec": clouduuid,
		"spec": map[string]interface{}{
			"modelIndex": "NodePoolSpec",
			"eksConfig": map[string]interface{}{
				"securityGroups": securitygroups,
				"keyName":        keyname,
				"diskSize":       diskSize,
				"instanceTypes":  instancetype,
				"imageId":        "",
			},
		},
	}

	data, err := apiClient.PostFromJSON(client.ServiceClusters, "clustertypes", clustertype, nil)
	if err != nil {
		return err
	}

	_, nerr := apiClient.PostFromJSON(client.ServiceClusters, "nodepooltypes", nodepoolobj, nil)
	if nerr != nil {
		return err
	}

	pmcID := data["id"].(string)
	d.SetId(pmcID)
	return nil
}

func resourceEksClusterTypeRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceEksClusterTypeUpdate(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceEksClusterTypeDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)

	name := d.Get("name").(string)

	id, err := apiClient.QueryByName(client.ServiceClusters, "clustertypes", name)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	params := map[string]string{
		"action": "delete",
	}

	if err := apiClient.Delete(id, params); err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Printf("Deleted cluster type %s", name)

	return nil
}
