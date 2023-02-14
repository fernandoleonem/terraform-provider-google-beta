package google

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

// Change
func TestAccComputeNetworkEdgeSecurityServices_withSecurityPolicies(t *testing.T) {
	t.Parallel()

	spName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	polLink := "google_compute_security_policy.policy.self_link"
	polNameStandard := fmt.Sprintf("tf-test-%s", randString(t, 10))
	polNameAdvanced := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSecurityPolicy_withDdosProtectionConfig(polNameStandard),
			},
			{
				ResourceName:      "google_compute_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSecurityPolicy_withDdosProtectionConfig_update(polNameAdvanced),
			},
			{
				ResourceName:      "google_compute_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeNetworkEdgeSecurityServices_withSecurityPolicies(spName, polLink),
			},
			{
				ResourceName:      "google_compute_network_edge_security_services.services",
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
    ddos_protection = "STANDARD"
  }

  adaptive_protection_config {
    layer_7_ddos_defense_config {
      enable = true
      rule_visibility = "STANDARD"
	}
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
`, spName)
}

func testAccComputeNetworkEdgeSecurityServices_withSecurityPolicies(spName, polLink string) string {
	return fmt.Sprintf(`
resource "google_compute_network_edge_security_services" "services" {
	name        = "%s"
	description = "basic network edge security services"
	security_policy = %s
}
`, spName, polLink)
}