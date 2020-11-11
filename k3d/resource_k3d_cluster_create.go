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
		Create: resourceClusterCreate,
		Read:   resourceClusterRead,
		// Update: resourceClusterUpdate,
		Delete: resourceClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(defaultCreateTimeout),
			// Update: schema.DefaultTimeout(defaultUpdateTimeout),
			Delete: schema.DefaultTimeout(defaultDeleteTimeout),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "The name of the resource, also acts as it's unique ID",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateName,
			},
			"servers": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "1",
				Description: "Specify how many servers you want to create (default 1)",
			},
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

	return resourceClusterRead(d, meta)
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

	// Computed values

	return nil
}

func resourceClusterDelete(d *schema.ResourceData, meta interface{}) error {
	if err := deleteCluster(d); err != nil {
		return err
	}

	d.SetId("")

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
