# terraform-provider-k3d

This project is a [terraform](http://www.terraform.io/) provider for
[k3d](https://k3d.io/).

## Requirements

* Terraform 0.13.x
* Go 1.15

## Usage

Create Cluster

```hcl

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
  servers = 1
}
```

## Developing

### Dependencies for building from source

If you need to build from source, you should have a working Go environment setup.
If not check out the Go [getting started](http://golang.org/doc/install) guide.

This project uses [Go Modules](https://github.com/golang/go/wiki/Modules) for dependency management.
To fetch all dependencies run `make mod` inside this repository.

### Build

```sh
make build
```

The binary will then be available at `build/$(GOOS)_$(GOARCH)/$(PLUGIN_NAME)_v$(VERSION)`

### Install

```sh
make install
```

This will place the binary under `$(HOME)/.terraform.d/plugins/$(HOSTNAME)/$(USER)/$(NAME)/$(VERSION)/$(GOOS)_$(GOARCH)`.
After installing you will need to run `terraform init` in any project using the plugin.

## License

MIT
