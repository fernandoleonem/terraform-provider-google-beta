package google

import (
	"fmt"
	"testing"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

//Change

func TestAccComputeNetworkEdgeSecurityServices_basic_withDdos_realTest(t *testing.T) {
	t.Parallel()

	spName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	polLink := "google_compute_security_policy.policy.ddos_protection_config.ddos_protection"
	polName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetworkEdgeSecurityServices_basic(spName, polLink, polName),
			},
			{
				ResourceName:      "google_compute_backend_bucket.image_backend",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}


func testAccComputeNetworkEdgeSecurityServices_basic(spName, polLink, polName string) string {
	return fmt.Sprintf(`
  resource "google_compute_network_edge_security_services" "services" {
	name        = "%s"
	description = "basic network edge security services"
	security_policy = "%s"
}
	  
resource "google_compute_security_policy" "policy" {
  name        = "%s"
  description = "basic security policy"
  type = "CLOUD_ARMOR_NETWORK"
}
`, spName, polLink, polName) 
}
