package appstream

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/appstream"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppstreamStack() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppstreamStackCreate,
		Read:   resourceAppstreamStackRead,
		Update: resourceAppstreamStackUpdate,
		Delete: resourceAppstreamStackDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"access_endpoints": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"vpce_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"application_settings": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"settings_group": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"embed_host_domains": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"feedback_url": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"redirect_url": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"storage_connectors": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,							
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connector_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"domains": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed: true,							
						},
						"resource_identifier": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,							
						},
					},
				},
			},

			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			"user_settings": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeString,
							Required: true,
						},
						"permission": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceAppstreamStackCreate(d *schema.ResourceData, meta interface{}) error {

	svc := meta.(*AWSClient).appstreamconn

	CreateStackInputOpts := &appstream.CreateStackInput{}
	if v, ok := d.GetOk("access_endpoints"); ok {
		CreateStackInputOpts.AccessEndpoints = expandAccessEndpointsConfigs(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("application_settings"); ok {
		CreateStackInputOpts.ApplicationSettings = expandApplicationSettings(v.([]interface{}))
	}

	if v, ok := d.GetOk("name"); ok {
		CreateStackInputOpts.Name = aws.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		CreateStackInputOpts.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("display_name"); ok {
		CreateStackInputOpts.DisplayName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("embed_host_domains"); ok {
		CreateStackInputOpts.EmbedHostDomains = expandStringSet(v.(*schema.Set))
	}

	if v, ok := d.GetOk("feedback_url"); ok {
		CreateStackInputOpts.FeedbackURL = aws.String(v.(string))
	}

	if v, ok := d.GetOk("redirect_url"); ok {
		CreateStackInputOpts.RedirectURL = aws.String(v.(string))
	}

	if v, ok := d.GetOk("storage_connectors"); ok {
		CreateStackInputOpts.StorageConnectors = expandStorageConnectorConfigs(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("user_settings"); ok {
		CreateStackInputOpts.UserSettings = expandUserSettingConfigs(v.(*schema.Set).List())
	}

	log.Printf("[DEBUG] Run configuration: %s", CreateStackInputOpts)

	resp, err := svc.CreateStack(CreateStackInputOpts)

	if err != nil {
		log.Printf("[ERROR] Error creating Appstream stack: %s", err)
		return err
	}

	log.Printf("[DEBUG] Appstream stack created %s ", resp)

	if v, ok := d.GetOk("tags"); ok {
		time.Sleep(2 * time.Second)

		stack_name := aws.StringValue(CreateStackInputOpts.Name)
		get, err := svc.DescribeStacks(&appstream.DescribeStacksInput{
			Names: aws.StringSlice([]string{stack_name}),
		})

		if err != nil {
			log.Printf("[ERROR] Error describing Appstream stack: %s", err)
			return err
		}

		if get.Stacks == nil {
			log.Printf("[DEBUG] Appstream stack (%s) not found", d.Id())
		}

		tag, err := svc.TagResource(&appstream.TagResourceInput{
			ResourceArn: get.Stacks[0].Arn,
			Tags:        aws.StringMap(expandTags(v.(map[string]interface{}))),
		})

		if err != nil {
			log.Printf("[ERROR] Error tagging Appstream stack: %s", err)
			return err
		}

		log.Printf("[DEBUG] %s", tag)
	}

	d.SetId(*CreateStackInputOpts.Name)

	return resourceAppstreamStackRead(d, meta)
}

func resourceAppstreamStackRead(d *schema.ResourceData, meta interface{}) error {

	svc := meta.(*AWSClient).appstreamconn

	var names []*string
	stackName := d.Id()
	names = append(names, &stackName)
	describeStacksInput := appstream.DescribeStacksInput{Names: names}

	resp, err := svc.DescribeStacks(&describeStacksInput)
	if err != nil {
		log.Printf("[ERROR] Error describing stacks: %s", err)
		return err
	}

	for _, v := range resp.Stacks {

		if aws.StringValue(v.Name) == d.Get("name") {
			d.Set("name", v.Name)
			d.Set("description", v.Description)
			d.Set("display_name", v.DisplayName)
			d.Set("embed_host_domains", v.EmbedHostDomains)
			d.Set("feedback_url", v.FeedbackURL)
			d.Set("redirect_url", v.RedirectURL)

			ae_res := make([]map[string]interface{}, 0)
			for _, raw := range v.AccessEndpoints {
				ae_attr := map[string]interface{}{}
				ae_attr["endpoint_type"] = aws.StringValue(raw.EndpointType)
				ae_attr["vpce_id"] = aws.StringValue(raw.VpceId)
				ae_res = append(ae_res, ae_attr)
			}

			if len(ae_res) > 0 {
				if err := d.Set("access_endpoints", ae_res); err != nil {
					log.Printf("[ERROR] Error setting access endpoints: %s", err)
				}
			}

			sc_res := make([]map[string]interface{}, 0)
			for _, raw := range v.StorageConnectors {
				sc_attr := map[string]interface{}{}
				sc_attr["connector_type"] = aws.StringValue(raw.ConnectorType)
				sc_attr["domains"] = aws.StringValueSlice(raw.Domains)
				sc_attr["resource_identifier"] = aws.StringValue(raw.ResourceIdentifier)
				sc_res = append(sc_res, sc_attr)
			}

			if len(sc_res) > 0 {
				if err := d.Set("storage_connectors", sc_res); err != nil {
					log.Printf("[ERROR] Error setting storage connector: %s", err)
				}
			}

			tg, err := svc.ListTagsForResource(&appstream.ListTagsForResourceInput{
				ResourceArn: v.Arn,
			})

			if err != nil {
				log.Printf("[ERROR] Error listing stack tags: %s", err)
				return err
			}

			if tg.Tags == nil {
				log.Printf("[DEBUG] Appstream stack tags (%s) not found", d.Id())
				return nil
			}

			if len(tg.Tags) > 0 {
				tags_attr := make(map[string]string)
				tags := tg.Tags
				for k, v := range tags {
					tags_attr[k] = aws.StringValue(v)
				}
				d.Set("tags", tags_attr)
			}

			us_res := make([]map[string]interface{}, 0)
			for _, raw := range v.UserSettings {
				us_attr := map[string]interface{}{}
				us_attr["action"] = aws.StringValue(raw.Action)
				us_attr["permission"] = aws.StringValue(raw.Permission)
				us_res = append(us_res, us_attr)
			}
			if len(us_res) > 0 {
				if err := d.Set("user_settings", us_res); err != nil {
					log.Printf("[ERROR] Error setting user settings: %s", err)
				}
			}
			return nil
		}
	}

	d.SetId("")

	return nil
}

func resourceAppstreamStackUpdate(d *schema.ResourceData, meta interface{}) error {

	svc := meta.(*AWSClient).appstreamconn

	UpdateStackInputOpts := &appstream.UpdateStackInput{}

	if v, ok := d.GetOk("name"); ok {
		UpdateStackInputOpts.Name = aws.String(v.(string))
	}

	if d.HasChange("access_endpoints") {
		log.Printf("[DEBUG] Modify appstream stack")
		access_endpoints := d.Get("access_endpoints").(*schema.Set).List()
		UpdateStackInputOpts.AccessEndpoints = expandAccessEndpointsConfigs(access_endpoints)
	}

	if d.HasChange("application_settings") {
		log.Printf("[DEBUG] Modify appstream stack")
		application_settings := d.Get("application_settings").([]interface{})
		UpdateStackInputOpts.ApplicationSettings = expandApplicationSettings(application_settings)
	}

	if d.HasChange("description") {
		log.Printf("[DEBUG] Modify appstream stack")
		description := d.Get("description").(string)
		UpdateStackInputOpts.Description = aws.String(description)
	}

	if d.HasChange("display_name") {
		log.Printf("[DEBUG] Modify appstream stack")
		displayname := d.Get("display_name").(string)
		UpdateStackInputOpts.DisplayName = aws.String(displayname)
	}

	if d.HasChange("embed_host_domains") {
		log.Printf("[DEBUG] Modify appstream stack")
		embed_host_domains := d.Get("embed_host_domains").(*schema.Set)
		UpdateStackInputOpts.EmbedHostDomains = expandStringSet(embed_host_domains)
	}

	if d.HasChange("feedback_url") {
		log.Printf("[DEBUG] Modify appstream stack")
		feedbackurl := d.Get("feedback_url").(string)
		UpdateStackInputOpts.FeedbackURL = aws.String(feedbackurl)
	}

	if d.HasChange("redirect_url") {
		log.Printf("[DEBUG] Modify appstream stack")
		redirecturl := d.Get("redirect_url").(string)
		UpdateStackInputOpts.RedirectURL = aws.String(redirecturl)
	}

	if d.HasChange("storage_connectors") {
		log.Printf("[DEBUG] Modify appstream stack")
		storage_connectors := d.Get("storage_connectors").(*schema.Set).List()
		UpdateStackInputOpts.StorageConnectors = expandStorageConnectorConfigs(storage_connectors)
	}

	if d.HasChange("user_settings") {
		log.Printf("[DEBUG] Modify appstream stack")
		user_settings := d.Get("user_settings").(*schema.Set).List()
		UpdateStackInputOpts.UserSettings = expandUserSettingConfigs(user_settings)
	}

	resp, err := svc.UpdateStack(UpdateStackInputOpts)

	if err != nil {
		log.Printf("[ERROR] Error updating Appstream Stack: %s", err)
		return err
	}

	log.Printf("[DEBUG] Appstream Stack updated %s ", resp)
	if v, ok := d.GetOk("tags"); ok && d.HasChange("tags") {
		time.Sleep(2 * time.Second)

		stack_name := aws.StringValue(UpdateStackInputOpts.Name)
		get, err := svc.DescribeStacks(&appstream.DescribeStacksInput{
			Names: aws.StringSlice([]string{stack_name}),
		})

		if err != nil {
			log.Printf("[ERROR] Error describing Appstream Stack: %s", err)
			return err
		}

		if get.Stacks == nil {
			log.Printf("[DEBUG] Appstream Stack (%s) not found", d.Id())
		}
		tag, err := svc.TagResource(&appstream.TagResourceInput{
			ResourceArn: get.Stacks[0].Arn,
			Tags:        aws.StringMap(expandTags(v.(map[string]interface{}))),
		})
		if err != nil {
			log.Printf("[ERROR] Error tagging Appstream Stack: %s", err)
			return err
		}
		log.Printf("[DEBUG] %s", tag)
	}
	return nil
}

func resourceAppstreamStackDelete(d *schema.ResourceData, meta interface{}) error {
	svc := meta.(*AWSClient).appstreamconn

	resp, err := svc.DeleteStack(&appstream.DeleteStackInput{
		Name: aws.String(d.Id()),
	})

	if err != nil {
		log.Printf("[ERROR] Error deleting Appstream stack: %s", err)
		return err
	}
	log.Printf("[DEBUG] %s", resp)
	return nil

}

func expandAccessEndpointConfigs(accessEndpointConfigs []interface{}) []*appstream.AccessEndpoint {
	accessEndpointConfig := []*appstream.AccessEndpoint{}

	for _, raw := range accessEndpointConfigs {
		configAttributes := raw.(map[string]interface{})
		configEndpointType := configAttributes["endpoint_type"].(string)
		configVpceId := configAttributes["vpce_id"].(string)
		config := &appstream.AccessEndpoint{
			EndpointType: aws.String(configEndpointType),
			VpceId:       aws.String(configVpceId),
		}
		accessEndpointConfig = append(accessEndpointConfig, config)
	}
	return accessEndpointConfig
}
