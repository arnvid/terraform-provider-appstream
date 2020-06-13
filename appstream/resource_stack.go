package appstream

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/appstream"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
        "log"
        "time"
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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connector_type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"user_settings": {
				Type: schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_download":{
							Type: schema.TypeString,
							Optional: true,
						},
						"file_upload":{
							Type: schema.TypeString,
							Optional: true,
						},
						"copy_from_local": {
							Type: schema.TypeString,
							Optional: true,
						},
						"copy_to_local": {
							Type: schema.TypeString,
							Optional: true,
						},
						"allow_local_device_printing": {
							Type: schema.TypeString,
							Optional: true,
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

	if v, ok := d.GetOk("name"); ok {
		CreateStackInputOpts.Name = aws.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		CreateStackInputOpts.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("display_name"); ok {
		CreateStackInputOpts.DisplayName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("feedback_url"); ok {
		CreateStackInputOpts.FeedbackURL = aws.String(v.(string))
	}

	if v, ok := d.GetOk("redirect_url"); ok {
		CreateStackInputOpts.RedirectURL = aws.String(v.(string))
	}

	if v, ok := d.GetOk("storage_connectors"); ok {
		storageConnectorConfigs := v.(*schema.Set).List()
		CreateStackInputOpts.StorageConnectors = expandStorageConnectorConfigs(storageConnectorConfigs)
	}

	if v, ok := d.GetOk("user_settings"); ok {
		userSettingConfigs := v.(*schema.Set).List()
		CreateStackInputOpts.UserSettings = expandUserSettingConfigs(userSettingConfigs)
	}

	log.Printf("[DEBUG] Run configuration: %s", CreateStackInputOpts)

	resp, err := svc.CreateStack(CreateStackInputOpts)

	if err != nil {
		log.Printf("[ERROR] Error creating Appstream Cluster: %s", err)
		return err
	}
	log.Printf("[DEBUG] Appstream stack created %s ", resp)
	time.Sleep(2 * time.Second)
	if v, ok := d.GetOk("tags"); ok {

		data_tags := v.(map[string]interface{})
		attr := make(map[string]string)
		for k, v := range data_tags {
			attr[k] = v.(string)
		}

		tags := aws.StringMap(attr)

		stack_name := aws.StringValue(CreateStackInputOpts.Name)
		get, err := svc.DescribeStacks(&appstream.DescribeStacksInput{
			Names: aws.StringSlice([]string{stack_name}),
		})
		if err != nil {
			log.Printf("[ERROR] Error describing Appstream Stack: %s", err)
			return err
		}
		if get.Stacks == nil {
			log.Printf("[DEBUG] Apsstream Stack (%s) not found", d.Id())
		}

		stackArn := get.Stacks[0].Arn

		tag, err := svc.TagResource(&appstream.TagResourceInput{
			ResourceArn: stackArn,
			Tags:        tags,
		})
		if err != nil {
			log.Printf("[ERROR] Error tagging Appstream Stack: %s", err)
			return err
		}
		log.Printf("[DEBUG] %s", tag)
	}

	log.Printf("[DEBUG] %s", resp)
	d.SetId(*CreateStackInputOpts.Name)

	return resourceAppstreamStackRead(d, meta)
}

func resourceAppstreamStackRead(d *schema.ResourceData, meta interface{}) error {

	svc := meta.(*AWSClient).appstreamconn

	resp, err := svc.DescribeStacks(&appstream.DescribeStacksInput{})
	if err != nil {
		log.Printf("[ERROR] Error describing stacks: %s", err)
		return err
	}

	for _, v := range resp.Stacks {

		if aws.StringValue(v.Name) == d.Get("name") {
			d.Set("name", v.Name)
			d.Set("description", v.Description)
			d.Set("display_name", v.DisplayName)
			d.Set("feedback_url", v.FeedbackURL)
			d.Set("redirect_url", v.RedirectURL)

			attr := map[string]interface{}{}
			res := make([]map[string]interface{}, 0)

			sc := v.StorageConnectors
			if len(sc) > 0 {
				attr["connector_type"] = aws.StringValue(sc[0].ConnectorType)
				res = append(res, attr)
			}

			if len(res) > 0 {
				if err := d.Set("storage_connectors", res); err != nil {
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
				log.Printf("[DEBUG] Apsstream Stack tags (%s) not found", d.Id())
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
			return nil
		}
	}

	d.SetId("")
	return nil

}

func resourceAppstreamStackUpdate(d *schema.ResourceData, meta interface{}) error {

	svc := meta.(*AWSClient).appstreamconn

	UpdateStackInputOpts := &appstream.UpdateStackInput{}

	d.Partial(true)

	if v, ok := d.GetOk("name"); ok {
		UpdateStackInputOpts.Name = aws.String(v.(string))
	}

	if d.HasChange("description") {
		d.SetPartial("description")
		log.Printf("[DEBUG] Modify appstream stack")
		description := d.Get("description").(string)
		UpdateStackInputOpts.Description = aws.String(description)
	}

	if d.HasChange("display_name") {
		d.SetPartial("display_name")
		log.Printf("[DEBUG] Modify appstream stack")
		displayname := d.Get("display_name").(string)
		UpdateStackInputOpts.DisplayName = aws.String(displayname)
	}

	if d.HasChange("feedback_url") {
		d.SetPartial("feedback_url")
		log.Printf("[DEBUG] Modify appstream stack")
		feedbackurl := d.Get("feedback_url").(string)
		UpdateStackInputOpts.FeedbackURL = aws.String(feedbackurl)
	}

	if d.HasChange("redirect_url") {
		d.SetPartial("redirect_url")
		log.Printf("[DEBUG] Modify appstream stack")
		redirecturl := d.Get("redirect_url").(string)
		UpdateStackInputOpts.RedirectURL = aws.String(redirecturl)
	}

	resp, err := svc.UpdateStack(UpdateStackInputOpts)

	if err != nil {
		log.Printf("[ERROR] Error updating Appstream Stack: %s", err)
		return err
	}
	log.Printf("[DEBUG] %s", resp)
	d.Partial(false)
	return nil

}

func resourceAppstreamStackDelete(d *schema.ResourceData, meta interface{}) error {

	svc := meta.(*AWSClient).appstreamconn

	resp, err := svc.DeleteStack(&appstream.DeleteStackInput{
		Name: aws.String(d.Id()),
	})
	if err != nil {
		log.Printf("[ERROR] Error deleting Appstream Stack: %s", err)
		return err
	}
	log.Printf("[DEBUG] %s", resp)
	return nil

}

func expandStorageConnectorConfigs(storageConnectorConfigs []interface{}) []*appstream.StorageConnector {
	storageConnectorConfig := []*appstream.StorageConnector{}

	for _, raw := range storageConnectorConfigs {
		configAttributes := raw.(map[string]interface{})
		configConnectorType := configAttributes["connector_type"].(string)
		config := &appstream.StorageConnector{
			ConnectorType: aws.String(configConnectorType),
		}
		storageConnectorConfig = append(storageConnectorConfig, config)
	}
	return storageConnectorConfig
}

func expandUserSettingConfigs(userSettingConfigs []interface{}) []*appstream.UserSetting {
	userSettingConfig := []*appstream.UserSetting{}

	for _, raw := range userSettingConfigs {
		configAttributes := raw.(map[string]interface{})
		configFileDownload := configAttributes["file_download"].(string)
		configFileUpload := configAttributes["file_upload"].(string)
		configCopyFromLocal := configAttributes["copy_from_local"].(string)
		configCopytoLocal := configAttributes["copy_to_local"].(string)
		configAllowLocalPrint := configAttributes["allow_local_device_printing"].(string)
		if configAttributes["file_download"] != nil {
			config := &appstream.UserSetting{
				Action: aws.String("FILE_DOWNLOAD"),
				Permission: aws.String(configFileDownload),
			}
			userSettingConfig = append(userSettingConfig, config)
		}
		if configAttributes["file_upload"] != nil {
			config := &appstream.UserSetting{
				Action: aws.String("FILE_UPLOAD"),
				Permission: aws.String(configFileUpload),
			}
			userSettingConfig = append(userSettingConfig, config)
		}
		if configAttributes["copy_from_local"] != nil {
			config := &appstream.UserSetting{
				Action: aws.String("CLIPBOARD_COPY_FROM_LOCAL_DEVICE"),
				Permission: aws.String(configCopyFromLocal),
			}
			userSettingConfig = append(userSettingConfig, config)
		}
		if configAttributes["copy_to_local"] != nil {
			config := &appstream.UserSetting{
				Action: aws.String("CLIPBOARD_COPY_TO_LOCAL_DEVICE"),
				Permission: aws.String(configCopytoLocal),
			}
			userSettingConfig = append(userSettingConfig, config)
		}
		if configAttributes["allow_local_device_printing"] != nil {
			config := &appstream.UserSetting{
				Action: aws.String("PRINTING_TO_LOCAL_DEVICE"),
				Permission: aws.String(configAllowLocalPrint),
			}
			userSettingConfig = append(userSettingConfig, config)
		}

	}
	return userSettingConfig
}
