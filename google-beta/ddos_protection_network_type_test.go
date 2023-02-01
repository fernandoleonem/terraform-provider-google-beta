package google

import (
	"fmt"
	"testing"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeSecurityPolicy_withDdosProtectionConfig(t *testing.T) {
	t.Parallel()

	spName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSecurityPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSecurityPolicy_withDdosProtectionConfig(spName),
			},
			{
				ResourceName:      "google_compute_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSecurityPolicy_withDdosProtectionConfig_update(spName),
			},
			{
				ResourceName:      "google_compute_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})	
}

func testAccComputeSecurityPolicy_withDdosProtectionConfig(spName string) string {
	return fmt.Sprintf(`
resource "google_compute_security_policy" "policy" {
  name        = "%s"
  description = "default rule"
  type = "CLOUD_ARMOR_NETWORK"
  region = "us-central1"
  
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
	description = "allow traffic from 192.0.2.0/24"
	match {
		versioned_expr = "SRC_IPS_V1"
		config {
		  src_ip_ranges = ["192.0.2.0/24"]
		}
		expr {
			expression = "us-central1"
		}
	  }
  }

  ddos_protection_config {
    ddos_protection = "STANDARD"
  }
}
`, spName)
}

func testAccComputeSecurityPolicy_withDdosProtectionConfig_update(spName string) string {
	return fmt.Sprintf(`
resource "google_compute_security_policy" "policy" {
  name        = "%s"
  description = "default rule"
  type = "CLOUD_ARMOR_NETWORK"
  region = "us-central1"
  
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
	description = "allow traffic from 192.0.2.0/24"
	match {
		versioned_expr = "SRC_IPS_V1"
		config {
		  src_ip_ranges = ["192.0.2.0/24"]
		}
		expr {
			expression = "us-central1"
		}
	  }
  }

  ddos_protection_config {
    ddos_protection = "ADVANCED"
  }
}
`, spName)
}