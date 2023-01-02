package google

import (
	"context"
	"fmt"
	"log"

	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	compute "google.golang.org/api/compute/v0.beta"
)

func resourceComputeNetworkEdgeSecurityServices() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeNetworkEdgeSecurityServicesCreate,
		Read:   resourceComputeSecurityPolicyRead,
		Update: resourceComputeSecurityPolicyUpdate,
		Delete: resourceComputeNetworkEdgeSecurityServicesDelete,

		Importer: &schema.ResourceImporter{
			State: resourceNetworkEdgeSecurityServicesImporter,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Update: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateGCEName,
				Description:  `Name of the resource. Provided by the client when the resource is created.`,
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `An optional description of this resource. Provide this property when you create the resource.`,
			},

			"fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Fingerprint of this resource. A hash of the contents stored in this object. This field is used in optimistic locking. This field will be ignored when inserting a NetworkEdgeSecurityService.`,
			},

			"security_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The resource URL for the network edge security service associated with this network edge security service.`,
			},


		},

		UseJSONNumber: true,
	}	
}

func resourceComputeNetworkEdgeSecurityServicesCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	sp := d.Get("name").(string)
	networkEdgeSecurityServices := &compute.NetworkEdgeSecurityService{
		Name:        sp,
		Description: d.Get("description").(string),
	}

	if v, ok := d.GetOk("fingerprint"); ok {
		networkEdgeSecurityServices.Fingerprint = v.(string)
	}

	if v, ok := d.GetOk("security_policy"); ok {
		networkEdgeSecurityServices.SecurityPolicy = v.(string)
	}

	log.Printf("[DEBUG] NetworkEdgeSecurityService insert request: %#v", networkEdgeSecurityServices)

	client := config.NewComputeClient(userAgent)

	op, err := client.NetworkEdgeSecurityServices.Insert(project, d.Get("region").(string), networkEdgeSecurityServices).Do()

	if err != nil {
		return errwrap.Wrapf("Error creating NetworkEdgeSecurityService: {{err}}", err)
	}

	id, err := replaceVars(d, config, "projects/{{project}}/global/networkEdgeSecurityServices/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = computeOperationWaitTime(config, op, project, fmt.Sprintf("Creating NetworkEdgeSecurityService %q", sp), userAgent, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	return resourceComputeSecurityPolicyRead(d, meta)


}

func resourceComputeNetworkEdgeSecurityServicesDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	client := config.NewComputeClient(userAgent)

	// Delete the SecurityPolicy
	op, err := client.NetworkEdgeSecurityServices.Delete(project, d.Get("region").(string), d.Get("name").(string)).Do()
	if err != nil {
		return errwrap.Wrapf("Error deleting NetworkEdgeSecurityServices: {{err}}", err)
	}

	err = computeOperationWaitTime(config, op, project, "Deleting NetworkEdgeSecurityServices", userAgent, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceNetworkEdgeSecurityServicesImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{"projects/(?P<project>[^/]+)/global/networkEdgeSecurityServices/(?P<name>[^/]+)", "(?P<project>[^/]+)/(?P<name>[^/]+)", "(?P<name>[^/]+)"}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/global/networkEdgeSecurityServices/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}