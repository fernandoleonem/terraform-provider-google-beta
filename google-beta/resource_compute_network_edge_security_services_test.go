package google

import (
	"fmt"
	"testing"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

//Change

func TestAccComputeNetworkEdgeSecurityServices_basic_withBucketImage_realTest(t *testing.T) {
	t.Parallel()

	bucketName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	polName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetworkEdgeSecurityServices_basic_withBucketImage(bucketName, polName, "google_compute_security_policy.policy.self_link"),
			},
			{
				ResourceName:      "google_compute_backend_bucket.image_backend",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeNetworkEdgeSecurityServices_basic_withBucketImage(bucketName, polName, "\"\""),
			},
			{
				ResourceName:      "google_compute_backend_bucket.image_backend",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeNetworkEdgeSecurityServices_basic_withDdos_realTest(t *testing.T) {
	t.Parallel()

	bucketName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	spName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	polName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	polLinkAll := "google_compute_security_policy.policy" 
	polLink := "google_compute_security_policy.policy.self_link"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetworkEdgeSecurityServices_basic_withDdos(bucketName, polLink, spName, polName, polLinkAll),
			},
			{
				ResourceName:      "google_compute_backend_bucket.image_backend",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeNetworkEdgeSecurityServices_basic_withBucketImage(bucketName, polName, polLink string) string {
	return fmt.Sprintf(`
resource "google_compute_backend_bucket" "image_backend" {
  name        = "%s"
  description = "Contains beautiful images"
  bucket_name = google_storage_bucket.image_bucket.name
  enable_cdn  = true
  edge_security_policy = %s
}

resource "google_storage_bucket" "image_bucket" {
  name     = "%s"
  location = "EU"
}


resource "google_compute_security_policy" "policy" {
  name        = "%s"
  description = "basic security policy"
  type = "CLOUD_ARMOR_NETWORK"
}
`, bucketName, polLink, bucketName, polName)
}

func testAccComputeNetworkEdgeSecurityServices_basic_withDdos(bucketName, polLink, spName, polName, polLinkAll string) string {
	return fmt.Sprintf(`
resource "google_compute_backend_bucket" "image_backend" {
	name        = "%s"
	description = "Contains beautiful images"
	bucket_name = google_storage_bucket.image_bucket.name
	enable_cdn  = true
	edge_security_policy = %s
  }
  
  resource "google_storage_bucket" "image_bucket" {
	name     = "%s"
	location = "EU"
  }

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
`, bucketName, polLink, bucketName, spName, polLinkAll, polName) 
}
