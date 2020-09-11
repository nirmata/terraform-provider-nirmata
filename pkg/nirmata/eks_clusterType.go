package nirmata

import (
	"fmt"
	"regexp"
	"time"

	guuid "github.com/google/uuid"
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
			"instance_types": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
	instanceTypes := d.Get("instance_types")
	keyName := d.Get("key_name").(string)
	securityGroups := d.Get("security_groups")
	clusterRoleArn := d.Get("cluster_role_arn").(string)
	vpcId := d.Get("vpc_id").(string)
	subnetId := d.Get("subnet_id")
	nodeSecurityGroups := d.Get("node_security_groups")
	nodeIamRole := d.Get("node_iam_role").(string)

	cloudCredID, err := apiClient.QueryByName(client.ServiceClusters, "CloudCredentials", credentials)
	fmt.Printf("Error - %v", cloudCredID)
	if err != nil {
		fmt.Printf("Error - %v", err)
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
					"vpcId":                 vpcId,
					"subnetId":              subnetId,
					"privateEndpointAccess": false,
					"clusterRoleArn":        clusterRoleArn,
					"securityGroups":        securityGroups,
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
				"instanceTypes":  instanceTypes,
				"imageId":        "",
			},
		},
	}

	data, err := apiClient.PostFromJSON(client.ServiceClusters, "clustertypes", clusterType, nil)
	if err != nil {
		return err
	}

	_, nerr := apiClient.PostFromJSON(client.ServiceClusters, "nodepooltypes", nodePoolObj, nil)
	if nerr != nil {
		return err
	}

	clusterTypeID, err := apiClient.QueryByName(client.ServiceClusters, "KubernetesCluster", name)
	if err != nil {
		fmt.Printf("Error ", err)
		return err
	}

	err = apiClient.WaitForState(clusterTypeID,"state","ready",500,"Failed to get cluster status")
	if err != nil {
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
