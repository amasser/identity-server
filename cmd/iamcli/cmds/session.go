package cmds

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tierklinik-dobersberg/identity-server/cmd/iamcli/session"
)

var tokenStore session.TokenStore

func init() {
	RootCommand.PersistentFlags().StringP("access-token", "t", "", "Access token to use. Can either be a file path or an environment variable prefixed with 'env:'")
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
		tokenStore = session.NewEnvLoader(parts[1])
		return nil
	}

	tokenStore = session.NewFileTokenStore(parts[0])

	return nil
}
