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
	polName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	fmt.Println(spName)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSecurityPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetworkEdgeSecurityServices_basic(spName, polName, "google_compute_security_policy.policy.self_link"),
			},
			{
				ResourceName:      "google_compute_network_edge_security_services.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeNetworkEdgeSecurityServices_basic(spName, polName, polLink string) string {
	return fmt.Sprintf(`
resource "google_compute_network_edge_security_services" "policy" {
  name        = "%s"
  description = "basic network edge security services"
  security_policy = "%s"
}

resource "google_compute_security_policy" "policy" {
	name        = "%s"
	description = "default rule"
	type = "CLOUD_ARMOR_NETWORK"
	
	rule {
	  action   = "deny-502"
	  priority = "2147483647"
  
	  match {
		  versioned_expr = "SRC_IPS_V1"
		  config {
			src_ip_ranges = ["*"]
		  }
		}
	}
  
	rule {
	  action   = "allow"
	  priority = "1000"
	  description = "allow traffic from 198.51.100.0/24"
	  match {
		  versioned_expr = "SRC_IPS_V1"
		  config {
			src_ip_ranges = ["198.51.100.0/24"]
		  }
		}
	}
  
	ddos_protection_config {
	  ddos_protection = "ADVANCED"
	}
  
	adaptive_protection_config {
	  layer_7_ddos_defense_config {
		enable = true
		rule_visibility = "STANDARD"
	  }
	}
  }
`, spName, polLink, polName)
}