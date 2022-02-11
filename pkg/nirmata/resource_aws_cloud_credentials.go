package nirmata

import (
	"fmt"
	"log"

	guuid "github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nirmata/go-client/pkg/client"
)

func resourceAWSCloudCredentials() *schema.Resource {
	return &schema.Resource{

		Create: resourceAWSCredentialsCreate,
		Read:   resourceAWSCredentialsRead,
		Update: resourceAWSCredentialsUpdate,
		Delete: resourceAWSCredentialsDelete,
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
			"access_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aws_role_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_key_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"secret_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAWSCredentialsCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	region := d.Get("region").(string)
	access_type := d.Get("access_type").(string)
	aws_role_arn := d.Get("aws_role_arn").(string)
	access_key_id := d.Get("access_key_id").(string)
	secret_key := d.Get("secret_key").(string)

	data := map[string]interface{}{
		"name":        name,
		"description": description,
		"region":      region,
		"type":        "AWS",
	}

	if access_type == "assume_role" {
		if aws_role_arn == "" {
			return fmt.Errorf("\n [ERROR] - aws role arn  is required")
		}
		data["awsRoleArn"] = aws_role_arn
		data["awsExternalId"] = guuid.New()
	} else if access_type == "access_key" {
		if access_key_id == "" || secret_key == "" {
			return fmt.Errorf("\n [ERROR] - access key id or secret key is required")
		}
		data["awsAccessKeyId"] = access_key_id
		data["awsSecretKey"] = secret_key
	} else {
		return fmt.Errorf(" [ERROR] - invalid access type. Select assume_role or access_key")
	}

	log.Printf("[DEBUG] - creating cloud credentials %s with %+v", name, data)
	credentialsData, err := apiClient.PostFromJSON(client.ServiceClusters, "CloudCredentials", data, nil)
	if err != nil {
		log.Printf("[ERROR] - failed to create cloud credential %s with data %v: %v", name, data, err)
		return err
	}
	credentialsDataUUID := credentialsData["id"].(string)
	d.SetId(credentialsDataUUID)
	log.Printf("[INFO] - created cloud credential %s %s", name, credentialsDataUUID)

	return nil
}

func resourceAWSCredentialsRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	id := client.NewID(client.ServiceClusters, "CloudCredentials", d.Id())

	_, err := apiClient.Get(id, &client.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] failed to retrieve cloud credential details %s (%s): %v", name, id, err)
		return err
	}

	log.Printf("[INFO] - retrieved cloud credential %s %s", name, id.UUID())
	return nil
}

var awsCredentialsMap = map[string]string{
	"access_key_id": "awsAccessKeyId",
	"secret_key":    "awsSecretKey",
	"aws_role_arn":  "awsRoleArn",
}

func resourceAWSCredentialsUpdate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	id := client.NewID(client.ServiceClusters, "Cloudcredentials", d.Id())
	credentialsChanges := buildChanges(d, awsCredentialsMap, "access_key_id", "secret_key", "aws_role_arn")
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

func resourceAWSCredentialsDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	name := d.Get("name").(string)
	id := client.NewID(client.ServiceClusters, "CloudCredentials", d.Id())

	if err := apiClient.Delete(id, nil); err != nil {
		return err
	}

	log.Printf("[INFO] - deleted cloud credential %s %s", name, id.UUID())
	return nil
}
