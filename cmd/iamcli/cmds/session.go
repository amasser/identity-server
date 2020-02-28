package cmds

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tierklinik-dobersberg/identity-server/pkg/client"
)

var tokenStore session.TokenStore

func init() {
	RootCommand.PersistentFlags().StringP("access-token", "t", "", `AuthN JWT access token used to authenticate against IAM.
Can either be a file path or an environment variable prefixed
with 'env:'. For example, --access-token env:TOKEN`)
}

func initTokenStore(cmd *cobra.Command) error {
	storeType, err := cmd.Flags().GetString("access-token")
	if err != nil {
		return err
	}

	parts := strings.Split(storeType, ":")
	if len(parts) > 2 {
		return errors.New("invalid value for --access-token")
	}

	if len(parts) == 2 && parts[0] == "env" {
		tokenStore = client.NewEnvLoader(parts[1])
		return nil
	}

	tokenStore = client.NewFileTokenStore(parts[0])

	return nil
}
