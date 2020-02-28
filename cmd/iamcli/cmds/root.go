package cmds

import (
	"log"

	"github.com/spf13/cobra"
)

// RootCommand is the main command of iamcli
var RootCommand = &cobra.Command{
	Use:   "iamcli",
	Short: "Manage users, groups and policies",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initTokenStore(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Fatal("Run with --help for more information")
	},
}
