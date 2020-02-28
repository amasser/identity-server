package cmds

import "github.com/spf13/cobra"

// RootCommand is the main command of iamcli
var RootCommand = &cobra.Command{
	Use:   "iamcli",
	Short: "Manage users, groups and policies",
}
