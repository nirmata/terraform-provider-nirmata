package nirmata

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"regexp"
	"strings"
	"time"

	guuid "github.com/google/uuid"
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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
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
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subnet_id": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"cluster_role_arn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"security_groups": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"key_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_security_groups": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"node_iam_role": {
				Type:     schema.TypeString,
				Required: true,
			},
			"disk_size": {
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
			"log_types": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
		},
	}
}

func resourceEksClusterTypeCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)

	clouduuid := guuid.New()
	nodepooluuid := guuid.New()

	name := d.Get("name").(string)
	version := d.Get("version").(string)
	credentials := d.Get("credentials").(string)
	region := d.Get("region").(string)
	diskSize := d.Get("disk_size").(int)
	instanceType := d.Get("instance_type")
	keyName := d.Get("key_name").(string)
	securityGroups := d.Get("security_groups")
	clusterRoleArn := d.Get("cluster_role_arn").(string)
	vpcID := d.Get("vpc_id").(string)
	subnetID := d.Get("subnet_id")
	nodeSecurityGroups := d.Get("node_security_groups")
	nodeIamRole := d.Get("node_iam_role").(string)
	logTypes := d.Get("log_types")

	cloudCredID, err := apiClient.QueryByName(client.ServiceClusters, "CloudCredentials", credentials)
	if err != nil {
		log.Printf("Error - %v", err)
		return err
	}

	clusterType := map[string]interface{}{
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
					"vpcId":                 vpcID,
					"subnetId":              subnetID,
					"privateEndpointAccess": false,
					"clusterRoleArn":        clusterRoleArn,
					"securityGroups":        securityGroups,
					"logTypes":              logTypes,
				},
			},
		},
	}

	nodePoolObj := map[string]interface{}{
		"id":              nodepooluuid,
		"modelIndex":      "NodePoolType",
		"name":            name + "-default-node-pool-type",
		"cloudConfigSpec": clouduuid,
		"spec": map[string]interface{}{
			"modelIndex": "NodePoolSpec",
			"eksConfig": map[string]interface{}{
				"securityGroups": nodeSecurityGroups,
				"nodeIamRole":    nodeIamRole,
				"keyName":        keyName,
				"diskSize":       diskSize,
				"instanceType":   instanceType,
				"imageId":        "",
			},
		},
	}
	txn := make(map[string]interface{})
	var objArr = make([]interface{}, 0)
	objArr = append(objArr, clusterType, nodePoolObj)
	txn["create"] = objArr
	data, err := apiClient.PostFromJSON(client.ServiceClusters, "txn", txn, nil)
	if err != nil {
		log.Printf("[ERROR] - failed to create cluster type  with data : %v", err)
		return err
	}

	changeID := data["changeId"].(string)
	d.SetId(changeID)

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
		log.Printf("[ERROR] - %v", err)
		return err
	}

	params := map[string]string{
		"action": "delete",
	}

	if err := apiClient.Delete(id, params); err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return nil
		}

		log.Printf("[ERROR] - %v", err)
		return err
	}

	log.Printf("Deleted cluster type %s", name)
	return nil
}
