package cmds

import (
	"github.com/spf13/cobra"
	"github.com/tierklinik-dobersberg/identity-server/pkg/client"
)

var iamClient *client.IdentityClient

func initClient(cmd *cobra.Command) error {
	iamClient = client.NewIdentityClient(iamServerURL,
		client.WithTokenLoader(tokenStore),
	)

	return nil
}
