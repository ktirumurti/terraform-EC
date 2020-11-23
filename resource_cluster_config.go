package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var URL_DEPLOYMENT = "https://api.elastic-cloud.com/api/v1/deployments"
var PRUNE_ORPHANS = "prune_orphans"

func deleteDeployment(apiKey, deploymentId string) error {
	client := &http.Client{}

	req, err := http.NewRequest("POST", URL_DEPLOYMENT+"/"+deploymentId+"/_shutdown", nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "ApiKey "+apiKey)

	_, err = client.Do(req)
	return err
}

func getDeployment(apiKey, deploymentId string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", URL_DEPLOYMENT+"/"+deploymentId, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "ApiKey "+apiKey)

	_, err = client.Do(req)
	return err
}

func createDeployment(apiKey, config string) (string, error) {
	var adjustedConfig map[string]interface{}
	var info map[string]interface{}

	var jsonStr = []byte(config)

	// remove the prune_orphans argument from the config if one exists
	// no-op if key isn't present
	json.Unmarshal(jsonStr, &adjustedConfig)
	delete(adjustedConfig, PRUNE_ORPHANS)
	updatedConfig, err := json.Marshal(adjustedConfig)

	if err != nil {
		return "", err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", URL_DEPLOYMENT, bytes.NewBuffer(updatedConfig))

	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "ApiKey "+apiKey)

	response, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	json.Unmarshal(body, &info)
	deployment_id := info["id"].(string)

	return deployment_id, err
}

func updateDeployment(apiKey, config, deploymentId string) error {
	client := &http.Client{}
	var info map[string]interface{}
	var jsonStr = []byte(config)

	// Check for the existence of a prune_orphans key. If it doesn't exist,
	// add it and set its default value to false
	json.Unmarshal(jsonStr, &info)
	if _, ok := info[PRUNE_ORPHANS]; !ok {
		info[PRUNE_ORPHANS] = false
	}

	newJsonStr, err := json.Marshal(info)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", URL_DEPLOYMENT+"/"+deploymentId, bytes.NewBuffer(newJsonStr))

	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "ApiKey "+apiKey)

	_, err = client.Do(req)
	return err
}

func resourceClusterConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceClusterConfigCreate,
		Read:   resourceClusterConfigRead,
		Update: resourceClusterConfigUpdate,
		Delete: resourceClusterConfigDelete,

		Schema: map[string]*schema.Schema{
			"api_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"config": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

/* Make an http request using the elasticsearch API to create a cluster with the
input config. The deployment id returned in the body of the response will be
maintained internally as part of the terraform state and used for subsequent
deployment access.
*/
func resourceClusterConfigCreate(d *schema.ResourceData, m interface{}) error {
	apiKey := d.Get("api_key").(string)
	config := d.Get("config").(string)
	deployment_id, err := createDeployment(apiKey, config)
	d.SetId(deployment_id)
	return err
}

//Query deployment based on deployment id - used for refreshing tfstate
func resourceClusterConfigRead(d *schema.ResourceData, m interface{}) error {
	apiKey := d.Get("api_key").(string)
	deploymentId := d.Id()
	return getDeployment(apiKey, deploymentId)
}

//Modify an existing deployment
func resourceClusterConfigUpdate(d *schema.ResourceData, m interface{}) error {
	apiKey := d.Get("api_key").(string)
	config := d.Get("config").(string)
	deploymentId := d.Id()
	return updateDeployment(apiKey, config, deploymentId)
}

//Delete an existing deployment
func resourceClusterConfigDelete(d *schema.ResourceData, m interface{}) error {
	apiKey := d.Get("api_key").(string)
	deploymentId := d.Id()
	return deleteDeployment(apiKey, deploymentId)
}
