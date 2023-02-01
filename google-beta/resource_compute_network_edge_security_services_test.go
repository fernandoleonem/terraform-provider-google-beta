package google

import (
	"fmt"
	"testing"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

//Change
func TestAccComputeNetworkEdgeSecurityServices_basic(t *testing.T) {
	t.Parallel()

	spName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	//polName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	fmt.Println(spName)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSecurityPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetworkEdgeSecurityServices_basic(spName),
			},
			{
				ResourceName:      "google_compute_network_edge_security_services.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeNetworkEdgeSecurityServices_basic(spName string) string {
	return fmt.Sprintf(`
resource "google_compute_network_edge_security_services" "primary" {
  name        = "%s"
  description = "basic network edge security services"
}
`, spName)
}