package k3d

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rancher/k3d/v3/pkg/cluster"
)

func resourceCluster() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the resource, also acts as it's unique ID",
				ValidateFunc: validateName,
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

func validateName(val interface{}, key string) (warns []string, errs []error) {
	name := val.(string)

	if err := cluster.CheckName(name); err != nil {
		errs = append(errs, fmt.Errorf("%s", err))
	}

	return
}

func resourceClusterCreate(d *schema.ResourceData, meta interface{}) error {

	name := ""
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	if err := createCluster(name); err != nil {
		return err
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
		return fmt.Errorf("Reading cluster: '%s'\n\n%s", id, string(out))
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

		if err := deleteCluster(id); err != nil {
			return err
		}

		if err := createCluster(name); err != nil {
			return err
		}

		d.SetId(name)
	}

	return nil
}

func resourceClusterDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	if err := deleteCluster(id); err != nil {
		return err
	}

	return nil
}

func deleteCluster(name string) error {
	log.Printf("[DEBUG] Deleting k3d cluster: %s", name)
	cmd := exec.Command("k3d", "cluster", "delete", name)
	out, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("Deleting cluster: '%s'\n\n%s", name, string(out))
	}

	return nil
}

func createCluster(name string) error {
	log.Printf("[DEBUG] Creating k3d cluster: %s", name)
	cmd := exec.Command("k3d", "cluster", "create", name)
	out, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("Creating cluster: '%s'\n\n%s", name, string(out))
	}

	return nil
}
