package cmds

import (
	"log"
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
		return initTokenStore(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Fatal("Run with --help for more information")
	},
}

func init() {
	iamServerEnv := os.Getenv("IAM_SERVER_URL")
	RootCommand.PersistentFlags().StringVarP(&iamServerURL, "server", "s", iamServerEnv, `Address of you IAM server url. If left empty,
it defaults to the value of the IAM_SERVER_URL
environment variable.`)
}
