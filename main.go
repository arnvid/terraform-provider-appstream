package main

import (
	//	"context"
	//	"flag"
	//	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/arnvid/terraform-provider-appstream/appstream"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return appstream.Provider()
		},
	})
}
