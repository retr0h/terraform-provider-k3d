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
			"servers": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "1",
				Description: "Specify how many servers you want to create (default 1)",
			},
		},

		Create: resourceClusterCreate,
		Read:   resourceClusterRead,
		Update: resourceClusterUpdate,
		Delete: resourceClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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
	if err := createCluster(d); err != nil {
		return err
	}

	name := d.Get("name").(string)
	d.SetId(name)

	return nil
}

func resourceClusterRead(d *schema.ResourceData, meta interface{}) error {
	out, err := listCluster(d)
	if err != nil {
		d.SetId("")

		return err
	}

	parts := strings.Fields(string(out))
	name := parts[0]
	servers := strings.Split(parts[1], "/")[0]
	d.Set("name", name)
	d.Set("servers", servers)

	return nil
}

// This may be completely wrong and stupid
func resourceClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	changed := false
	if d.HasChange("name") {
		changed = true
	}

	if d.HasChange("servers") {
		changed = true
	}

	if changed {
		if err := deleteCluster(d); err != nil {
			return err
		}

		if err := createCluster(d); err != nil {
			return err
		}

		name := d.Get("name").(string)
		d.SetId(name)
	}

	return nil
}

func resourceClusterDelete(d *schema.ResourceData, meta interface{}) error {
	if err := deleteCluster(d); err != nil {
		return err
	}

	return nil
}

func deleteCluster(d *schema.ResourceData) error {
	id := d.Id()

	log.Printf("[DEBUG] Deleting k3d cluster: %s", id)
	args := []string{"cluster", "delete", id}
	cmd := exec.Command("k3d", args...)
	out, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("Deleting cluster: '%s'\n\n%s", id, string(out))
	}

	return nil
}

func createCluster(d *schema.ResourceData) error {
	name := d.Get("name").(string)
	servers := d.Get("servers").(string)

	log.Printf("[DEBUG] Creating k3d cluster: %s", name)
	args := []string{"cluster", "create", name, "--servers", servers}
	cmd := exec.Command("k3d", args...)
	out, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("Creating cluster: '%s'\n\n%s", name, string(out))
	}

	return nil
}

func listCluster(d *schema.ResourceData) ([]byte, error) {
	id := d.Id()

	log.Printf("[DEBUG] Read k3d cluster: %s", id)
	args := []string{"cluster", "list", id, "--no-headers"}
	cmd := exec.Command("k3d", args...)
	out, err := cmd.CombinedOutput()

	if err != nil {
		return out, fmt.Errorf("Reading cluster: '%s'\n\n%s", id, string(out))
	}

	return out, nil
}
