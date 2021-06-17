package appstream

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/appstream"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppstreamUsageReportSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppstreamUsageReportSubscriptionCreate,
		Read:   resourceAppstreamUsageReportSubscriptionRead,
		Delete: resourceAppstreamUsageReportSubscriptionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"s3_bucket_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schedule": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAppstreamUsageReportSubscriptionCreate(d *schema.ResourceData, meta interface{}) error {

	svc := meta.(*AWSClient).appstreamconn

	CreateUsageReportSubscriptionInputOpts := &appstream.CreateUsageReportSubscriptionInput{}

	log.Printf("[DEBUG] Run configuration: %s", CreateUsageReportSubscriptionInputOpts)

	resp, err := svc.CreateUsageReportSubscription(CreateUsageReportSubscriptionInputOpts)

	if err != nil {
		log.Printf("[ERROR] Error creating Appstream Usage Report: %s", err)
		return err
	}

	log.Printf("[DEBUG] Usage Report created %s", resp)

	d.SetId(*resp.S3BucketName)
	d.Set("s3_bucket_name", resp.S3BucketName)
	d.Set("schedule", resp.Schedule)

	return nil
}

func resourceAppstreamUsageReportSubscriptionRead(d *schema.ResourceData, meta interface{}) error {

	svc := meta.(*AWSClient).appstreamconn

	resp, err := svc.DescribeUsageReportSubscriptions(&appstream.DescribeUsageReportSubscriptionsInput{})
	if err != nil {
		log.Printf("[ERROR] Error describing Appstream Usage Reports Subscription: %s", err)
		return err
	}
	for _, v := range resp.UsageReportSubscriptions {

		if aws.StringValue(v.S3BucketName) == d.Get("s3_bucket_name") {
			d.Set("s3_bucket_name", v.S3BucketName)
			d.Set("schedule", v.Schedule)

			return nil
		}
	}

	d.SetId("")
	return nil

}

func resourceAppstreamUsageReportSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	svc := meta.(*AWSClient).appstreamconn

	del, err := svc.DeleteUsageReportSubscription(&appstream.DeleteUsageReportSubscriptionInput{})

	if err != nil {
		log.Printf("[ERROR] Error deleting Appstream Usage Report Subscription: %s", err)
		return err
	}
	log.Printf("[DEBUG] %s", del)

	return nil
}
