package nirmata

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nirmata/go-client/pkg/client"
)

func resoureAwsRoleCredentials() *schema.Resource {
	return &schema.Resource{

		Create: resourceAwsRoleCreate,
		Read:   resourceAwsRoleRead,
		Update: resourceAwsRoleUpdate,
		Delete: resourceAwsRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateName,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aws_access_key_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"aws_secret_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"aws_role_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"aws_external_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAwsRoleCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	region := d.Get("region").(string)
	awsAccessKeyId := d.Get("aws_access_key_id").(string)
	awsSecretKey := d.Get("aws_secret_key").(string)
	awsRoleArn := d.Get("aws_role_arn").(string)
	awsExternalId := d.Get("aws_external_id").(string)

	data := map[string]interface{}{
		"name":           name,
		"description":    description,
		"type":           "AWS",
		"region":         region,
		"awsAccessKeyId": awsAccessKeyId,
		"awsSecretKey":   awsSecretKey,
		"awsRoleArn":     awsRoleArn,
		"awsExternalId":  awsExternalId,
	}

	if awsExternalId == "" && awsAccessKeyId == "" {
		return fmt.Errorf("\nError - access type is required")
	}

	log.Printf("[DEBUG] - creating aws cloud credentials %s with %+v", name, data)
	credData, err := apiClient.PostFromJSON(client.ServiceClusters, "Cloudcredentials", data, nil)
	if err != nil {
		log.Printf("[ERROR] - failed to creating aws cloud credentials %s with data %v: %v", name, data, err)
		return err
	}

	credUUID := credData["id"].(string)
	d.SetId(credUUID)
	log.Printf("[INFO] - created aws cloud credentials%s %s", name, credUUID)

	return nil
}

func resourceAwsRoleRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

var credentialsMap = map[string]string{
	"aws_access_key_id": "awsAccessKeyId",
	"aws_secret_key":    "awsSecretKey",
	"aws_role_arn":      "awsRoleArn",
	"aws_external_id":   "awsExternalId",
}

func resourceAwsRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	id := client.NewID(client.ServiceClusters, "Cloudcredentials", d.Id())
	credentialsChanges := buildChanges(d, credentialsMap, "aws_access_key_id", "aws_secret_key", "aws_role_arn", "aws_external_id")
	if len(credentialsChanges) > 0 {
		credentialsData, err := apiClient.Get(id, &client.GetOptions{})
		if err != nil {
			log.Printf("[ERROR] - failed to retrieve %s from %v: %v", "Cloudcredentials", id.Map(), err)
			return err
		}
		d, plainErr := client.NewObject(credentialsData)
		if plainErr != nil {
			log.Printf("[ERROR] - failed to decode %s %v: %v", "Cloudcredentials", d, err)
			return err
		}
		_, plainErr = apiClient.PutWithIDFromJSON(d.ID(), credentialsChanges)
		if plainErr != nil {
			log.Printf("[ERROR] - failed to update %s %v: %v", "Cloudcredentials", d.ID().Map(), err)
			return err
		}
		log.Printf("[DEBUG] updated %v %v", d.ID().Map(), credentialsChanges)
	}
	return nil
}

func resourceAwsRoleDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	id := client.NewID(client.ServiceClusters, "Cloudcredentials", d.Id())

	if err := apiClient.Delete(id, nil); err != nil {
		return err
	}

	log.Printf("[INFO] - deleted aws cloud credentials %s %s", name, id.UUID())
	return nil
}
