package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-providers/terraform-provider-vault/deprecated"
)

func main() {
	// TODO how to handle the new vs. deprecated providers?
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: deprecated.Provider})
}
