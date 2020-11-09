terraform {
  required_version = ">= 0.13.0"
  required_providers {
    k3d = {
      source = "github.com/retr0h/k3d"
      version = "1.0"
    }
  }
}

provider "k3d" {}

resource "k3d_cluster" "local" {
  name = "example-cluster"
}
