package cmds

import (
	"context"
	"fmt"
	"log"

	"github.com/ghodss/yaml"
	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
	"github.com/tierklinik-dobersberg/identity-server/pkg/iam"
	"golang.org/x/crypto/ssh/terminal"
)

var userRootCommand = &cobra.Command{
	Use:   "users",
	Short: "Manage users stored in IAM.",
}

var listUsersCommand = &cobra.Command{
	Use:   "list",
	Short: "List all users stored in IAM.",
	Run: func(cmd *cobra.Command, args []string) {
		uc := iamClient.Users()

		users, err := uc.Users(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		tw := table.NewWriter()
		tw.AppendHeader(table.Row{"", "Username", "URN"})

		for _, u := range users {
			tw.AppendRow(table.Row{u.AccountID, u.Username, u.ID})
		}

		tw.SetStyle(table.StyleLight)
		tw.Style().Options.SeparateColumns = false
		tw.Style().Options.DrawBorder = false

		fmt.Println(tw.Render())
	},
}

var loadUserCommand = &cobra.Command{
	Use:     "get",
	Aliases: []string{"load", "show"},
	Short:   "Load all data available for a given user.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		uc := iamClient.Users()

		urn := iam.UserURN(args[0])

		if urn.AccountID() == "" {
			urn = iam.UserURN("urn:iam::user/" + urn)
		}

		user, err := uc.LoadUser(context.Background(), urn)
		if err != nil {
			log.Fatal(err)
		}

		blob, err := yaml.Marshal(user)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(blob))
	},
}

var deleteUserCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete an existing user.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		uc := iamClient.Users()

		urn := iam.UserURN(args[0])
		if urn.AccountID() == "" {
			urn = iam.UserURN("urn:iam::user/" + urn)
		}

		if err := uc.DeleteUser(context.Background(), urn); err != nil {
			log.Fatal(err)
		}
	},
}

var createUserCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a new user account in IAM.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		password, _ := cmd.Flags().GetString("password")

		if password == "" {
			fmt.Print("Password: ")
			pwd, err := terminal.ReadPassword(0)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("")

			password = string(pwd)
		}

		uc := iamClient.Users()

		urn, err := uc.CreateUser(context.Background(), username, password, nil)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(urn)
	},
}

func init() {
	RootCommand.AddCommand(userRootCommand)

	createUserCommand.Flags().StringP("password", "p", "", "Password for the new user.")

	userRootCommand.AddCommand(
		listUsersCommand,
		loadUserCommand,
		deleteUserCommand,
		createUserCommand,
	)
}
