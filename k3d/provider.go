package k3d

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"k3d_cluster": resourceCluster(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
	}
}
