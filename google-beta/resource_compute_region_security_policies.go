package google

import (
	"fmt"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	compute "google.golang.org/api/compute/v0.beta"
	"log"
	"time"
)

// Change
func resourceComputeRegionSecurityPolicies() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeRegionSecurityPoliciesCreate,
		Read:   resourceComputeRegionSecurityPoliciesRead,
		Update: resourceComputeRegionSecurityPoliciesUpdate,
		Delete: resourceComputeRegionSecurityPoliciesDelete,

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
				Required:    true,
				ForceNew:    true,
				Description: `The resource URL for the network edge security service associated with this network edge security service.`,
			},
		},

		UseJSONNumber: true,
	}
}

func resourceComputeRegionSecurityPoliciesCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	sp := d.Get("name").(string)

	regionSecurityPolicies := &compute.RegionSecurityPoliciesService{}
	securityPolicy := &compute.SecurityPolicy{}

	log.Printf("[DEBUG] RegionSecurityPolicies insert request: %#v", regionSecurityPolicies)

	client := config.NewComputeClient(userAgent)

	op, err := client.RegionSecurityPolicies.Insert(project, region, securityPolicy).Do()

	if err != nil {
		return errwrap.Wrapf("Error creating RegionSecurityPolicies: {{err}}", err)
	}

	id, err := replaceVars(d, config, "projects/{{project}}/global/regionSecurityPolicies/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = computeOperationWaitTime(config, op, project, fmt.Sprintf("Creating RegionSecurityPolicies %q", sp), userAgent, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	return resourceComputeRegionSecurityPoliciesRead(d, meta)
}

func resourceComputeRegionSecurityPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	sp := d.Get("name").(string)

	client := config.NewComputeClient(userAgent)

	regionSecurityPolicies, err := client.RegionSecurityPolicies.Get(project, region, sp).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("RegionSecurityPolicies %q", d.Id()))
	}

	if err := d.Set("name", regionSecurityPolicies.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	/*if err := d.Set("description", regionSecurityPolicies.Description); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err := d.Set("fingerprint", regionSecurityPolicies.Fingerprint); err != nil {
		return fmt.Errorf("Error setting fingerprint: %s", err)
	}*/

	return nil
}

func resourceComputeRegionSecurityPoliciesUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	sp := d.Get("name").(string)

	//regionSecurityPolicies := &compute.RegionSecurityPoliciesService{}
	securityPolicy := &compute.SecurityPolicy{}

	/*if d.HasChange("description") {
		networkEdgeSecurityServices.Description = d.Get("description").(string)
		networkEdgeSecurityServices.ForceSendFields = append(networkEdgeSecurityServices.ForceSendFields, "Description")
	}

	if d.HasChange("security_policy") {
		networkEdgeSecurityServices.SecurityPolicy = d.Get("security_policy").(string)
		networkEdgeSecurityServices.ForceSendFields = append(networkEdgeSecurityServices.ForceSendFields, "SecurityPolicy")
	}*/

	/*if len(regionSecurityPolicies.ForceSendFields) > 0 {

	}*/

	client := config.NewComputeClient(userAgent)

	op, err := client.RegionSecurityPolicies.Patch(project, region, sp, securityPolicy).Do()

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error updating RegionSecurityPolicies %q: {{err}}", sp), err)
	}

	err = computeOperationWaitTime(config, op, project, fmt.Sprintf("Updating RegionSecurityPolicies %q", sp), userAgent, d.Timeout(schema.TimeoutUpdate))

	if err != nil {
		return err
	}

	return resourceComputeRegionSecurityPoliciesRead(d, meta)
}

func resourceComputeRegionSecurityPoliciesDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	client := config.NewComputeClient(userAgent)

	// Delete the SecurityPolicy
	op, err := client.RegionSecurityPolicies.Delete(project, region, d.Get("name").(string)).Do()
	if err != nil {
		return errwrap.Wrapf("Error deleting RegionSecurityPolicies: {{err}}", err)
	}

	err = computeOperationWaitTime(config, op, project, "Deleting RegionSecurityPolicies", userAgent, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceComputeRegionSecurityPoliciesImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{"projects/(?P<project>[^/]+)/global/regionSecurityPolicies/(?P<name>[^/]+)", "(?P<project>[^/]+)/(?P<name>[^/]+)", "(?P<name>[^/]+)"}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/global/regionSecurityPolicies/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
