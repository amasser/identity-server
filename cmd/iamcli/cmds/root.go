package cmds

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	iamServerURL string
)

// RootCommand is the main command of iamcli
var RootCommand = &cobra.Command{
	Use:   "iamcli",
	Short: "Manage users, groups and policies of your IAM server instance.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := initTokenStore(cmd); err != nil {
			return err
		}

		if err := initClient(cmd); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	iamServerEnv := os.Getenv("IAM_SERVER_URL")
	RootCommand.PersistentFlags().StringVarP(&iamServerURL, "server", "s", iamServerEnv, `Address of you IAM server url. If left empty,
it defaults to the value of the IAM_SERVER_URL
environment variable.`)
}
