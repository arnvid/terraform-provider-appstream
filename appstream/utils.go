package appstream

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/appstream"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandStringList(configured []interface{}) []*string {
	vs := make([]*string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, aws.String(v.(string)))
		}
	}
	return vs
}

func expandStringSet(configured *schema.Set) []*string {
	return expandStringList(configured.List())
}

func flattenStringList(list []*string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, *v)
	}
	return vs
}

func flattenStringSet(list []*string) *schema.Set {
	return schema.NewSet(schema.HashString, flattenStringList(list))
}

func expandAccessEndpointsConfigs(accessEndpointsConfigs []interface{}) []*appstream.AccessEndpoint {
	accessEndpointsConfig := []*appstream.AccessEndpoint{}

	for _, raw := range accessEndpointsConfigs {
		configAttributes := raw.(map[string]interface{})

		config := &appstream.AccessEndpoint{}
		if v, ok := configAttributes["endpoint_type"]; ok {
			config.EndpointType = aws.String(v.(string))
		}

		if v, ok := configAttributes["vpce_id"]; ok {
			config.VpceId = aws.String(v.(string))
		}

		accessEndpointsConfig = append(accessEndpointsConfig, config)
	}

	return accessEndpointsConfig
}

func expandStorageConnectorConfigs(storageConnectorConfigs []interface{}) []*appstream.StorageConnector {
	storageConnectorConfig := []*appstream.StorageConnector{}

	for _, raw := range storageConnectorConfigs {
		configAttributes := raw.(map[string]interface{})

		config := &appstream.StorageConnector{}
		if v, ok := configAttributes["connector_type"]; ok {
			config.ConnectorType = aws.String(v.(string))
		}

		log.Printf("[***] %s", configAttributes["domains"])
		if v, ok := configAttributes["domains"]; ok && len(v.([]interface{})) > 0 {
			config.Domains = v.([]*string)
		}

		if v, ok := configAttributes["resource_identifier"]; ok && len(v.(string)) > 0 {
			config.ResourceIdentifier = aws.String(v.(string))
		}

		storageConnectorConfig = append(storageConnectorConfig, config)
	}

	return storageConnectorConfig
}

func expandUserSettingConfigs(userSettingConfigs []interface{}) []*appstream.UserSetting {
	userSettingConfig := []*appstream.UserSetting{}

	for _, raw := range userSettingConfigs {
		configAttributes := raw.(map[string]interface{})

		config := &appstream.UserSetting{}
		if v, ok := configAttributes["action"]; ok {
			config.Action = aws.String(v.(string))
		}

		if v, ok := configAttributes["permission"]; ok {
			config.Permission = aws.String(v.(string))
		}

		userSettingConfig = append(userSettingConfig, config)
	}

	return userSettingConfig
}

func expandApplicationSettings(applicationSettings []interface{}) *appstream.ApplicationSettings {
	ApplicationSettings := &appstream.ApplicationSettings{}
	attr := applicationSettings[0].(map[string]interface{})

	if v, ok := attr["enabled"]; ok {
		ApplicationSettings.Enabled = aws.Bool(v.(bool))
	}

	if v, ok := attr["settings_group"]; ok {
		ApplicationSettings.SettingsGroup = aws.String(v.(string))
	}

	return ApplicationSettings
}

func expandVpcConfigs(vpcConfigs []interface{}) *appstream.VpcConfig {
	VpcConfigConfig := &appstream.VpcConfig{}
	attr := vpcConfigs[0].(map[string]interface{})

	if v, ok := attr["security_group_ids"]; ok {
		VpcConfigConfig.SecurityGroupIds = expandStringSet(v.(*schema.Set))
	}

	if v, ok := attr["subnet_ids"]; ok {
		VpcConfigConfig.SubnetIds = expandStringSet(v.(*schema.Set))
	}

	return VpcConfigConfig
}

func expandTags(data_tags map[string]interface{}) map[string]string {
	attr := make(map[string]string)
	for k, v := range data_tags {
		attr[k] = v.(string)
	}

	return attr
}
