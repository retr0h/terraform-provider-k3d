package k3d

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCluster() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the resource, also acts as it's unique ID",
				// ValidateFunc: validateName,
			},
		},

		Create: resourceClusterCreate,
		Read:   resourceClusterRead,
		Update: resourceClusterUpdate,
		Delete: resourceClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		// Schema: map[string]*schema.Schema{},
	}
}

func resourceClusterCreate(d *schema.ResourceData, meta interface{}) error {
	name := ""
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}
	out, err := createCluster(name)

	if err != nil {
		return fmt.Errorf("Error creating cluster: '%s'\n%s", name, string(out))
	}

	d.SetId(name)
	return nil
}

func resourceClusterRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	log.Printf("[DEBUG] Read k3d cluster: %s", id)
	cmd := exec.Command("k3d", "cluster", "list", id, "--no-headers")
	out, err := cmd.CombinedOutput()

	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error reading cluster: '%s'\n%s", id, string(out))
	}

	parts := strings.Fields(string(out))

	d.Set("name", parts[0])
	return nil
}

// This may be completely wrong and stupid
func resourceClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	if d.HasChange("name") {
		name := d.Get("name").(string)
		out, err := deleteCluster(id)

		if err != nil {
			return fmt.Errorf("Error deleting cluster: '%s'\n%s", id, string(out))
		}

		out, err = createCluster(name)

		if err != nil {
			return fmt.Errorf("Error creating cluster: '%s'\n%s", name, string(out))
		}

		d.SetId(name)
	}

	return nil
}

func resourceClusterDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	out, err := deleteCluster(id)

	if err != nil {
		return fmt.Errorf("Error deleting cluster: '%s'\n%s", id, string(out))
	}

	return nil
}

func deleteCluster(name string) ([]byte, error) {
	log.Printf("[DEBUG] Deleting k3d cluster: %s", name)
	cmd := exec.Command("k3d", "cluster", "delete", name)

	return cmd.CombinedOutput()
}

func createCluster(name string) ([]byte, error) {
	log.Printf("[DEBUG] Creating k3d cluster: %s", name)
	cmd := exec.Command("k3d", "cluster", "create", name)

	return cmd.CombinedOutput()
}
