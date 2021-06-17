package appstream

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/appstream"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppstreamStackAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppstreamStackAttachmentCreate,
		Read:   resourceAppstreamStackAttachmentRead,
		Update: resourceAppstreamStackAttachmentUpdate,
		Delete: resourceAppstreamStackAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"appstream_stack_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"appstream_fleet_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAppstreamStackAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	svc := meta.(*AWSClient).appstreamconn
	AssociateFleetInputOpts := &appstream.AssociateFleetInput{}

	if stack, ok := d.GetOk("appstream_stack_id"); ok {
		AssociateFleetInputOpts.StackName = aws.String(stack.(string))
	}

	if fleet, ok := d.GetOk("appstream_fleet_id"); ok {
		AssociateFleetInputOpts.FleetName = aws.String(fleet.(string))
	}

	log.Printf("[DEBUG] Run configuration: %s", AssociateFleetInputOpts)

	resp, err := svc.AssociateFleet(AssociateFleetInputOpts)
	if err != nil {
		log.Printf("[ERROR] Error associating Appstream Fleet to Stack: %s", err)
		return err
	}

	log.Printf("[DEBUG] Appstream Fleet associated to Stack %s ", resp)

	d.SetId(fmt.Sprintf("%s_%s", *AssociateFleetInputOpts.StackName, *AssociateFleetInputOpts.FleetName))

	return resourceAppstreamStackAttachmentRead(d, meta)
}

func resourceAppstreamStackAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	svc := meta.(*AWSClient).appstreamconn

	AssociationId := strings.Split(d.Id(), "_")

	resp, err := svc.ListAssociatedFleets(&appstream.ListAssociatedFleetsInput{
		StackName: &AssociationId[0],
	})

	if err != nil {
		log.Printf("[ERROR] Error describing associations: %s", err)
		return err
	}

	stack := AssociationId[0]
	for _, fleet := range resp.Names {
		d.Set("appstream_stack_id", stack)
		d.Set("appstream_fleet_id", fleet)

		return nil
	}

	d.SetId("")

	return nil
}

func resourceAppstreamStackAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	svc := meta.(*AWSClient).appstreamconn
	DisassociateFleetInputOpts := &appstream.DisassociateFleetInput{}
	AssociateFleetInputOpts := &appstream.AssociateFleetInput{}

	AssociationId := strings.Split(d.Id(), "_")
	DisassociateFleetInputOpts.StackName = &AssociationId[0]
	DisassociateFleetInputOpts.FleetName = &AssociationId[1]

	AssociateFleetInputOpts.StackName = DisassociateFleetInputOpts.StackName
	AssociateFleetInputOpts.FleetName = DisassociateFleetInputOpts.FleetName

	// d.Partial(true)

	if d.HasChange("appstream_stack_id") {
		//		d.SetPartial("appstream_stack_id")
		log.Printf("[DEBUG] Modify appstream association")
		appstream_stack_id := d.Get("appstream_stack_id").(string)
		AssociateFleetInputOpts.StackName = aws.String(appstream_stack_id)
	}

	if d.HasChange("appstream_fleet_id") {
		//		d.SetPartial("appstream_fleet_id")
		log.Printf("[DEBUG] Modify appstream association")
		appstream_fleet_id := d.Get("appstream_fleet_id").(string)
		AssociateFleetInputOpts.FleetName = aws.String(appstream_fleet_id)
	}

	if d.HasChanges("appstream_stack_id", "appstream_fleet_id") {
		dis_resp, dis_err := svc.DisassociateFleet(&appstream.DisassociateFleetInput{
			StackName: aws.String(*DisassociateFleetInputOpts.StackName),
			FleetName: aws.String(*DisassociateFleetInputOpts.FleetName),
		})

		if dis_err != nil {
			log.Printf("[ERROR] Error disassociating Appstream Fleet from Stack: %s", dis_err)
			return dis_err
		}

		log.Printf("[DEBUG] %s", dis_resp)

		ass_resp, ass_err := svc.AssociateFleet(AssociateFleetInputOpts)
		if ass_err != nil {
			log.Printf("[ERROR] Error associating Appstream Fleet to Stack: %s", ass_err)
			return ass_err
		}

		d.SetId(fmt.Sprintf("%s_%s", *AssociateFleetInputOpts.StackName, *AssociateFleetInputOpts.FleetName))

		log.Printf("[DEBUG] Appstream Fleet associated to Stack %s ", ass_resp)
	}

	//	d.Partial(false)

	return nil
}

func resourceAppstreamStackAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	svc := meta.(*AWSClient).appstreamconn

	AssociationId := strings.Split(d.Id(), "_")

	resp, err := svc.DisassociateFleet(&appstream.DisassociateFleetInput{
		StackName: aws.String(AssociationId[0]),
		FleetName: aws.String(AssociationId[1]),
	})

	if err != nil {
		log.Printf("[ERROR] Error disassociating Appstream Fleet from Stack: %s", err)
		return err
	}

	log.Printf("[DEBUG] %s", resp)

	return nil
}
