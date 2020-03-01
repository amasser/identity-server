package cmds

import (
	"context"
	"fmt"
	"log"

	"github.com/ghodss/yaml"
	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
	"github.com/tierklinik-dobersberg/identity-server/pkg/iam"
)

var groupRootCommand = &cobra.Command{
	Use:     "groups",
	Aliases: []string{"group", "g"},
	Short:   "Manage groups and memberships stored in IAM.",
}

var listGroupsCommand = &cobra.Command{
	Use:   "list",
	Short: "List all groups stored in IAM.",
	Run: func(cmd *cobra.Command, args []string) {
		gc := iamClient.Groups()

		groups, err := gc.Get(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		tw := table.NewWriter()
		tw.AppendHeader(table.Row{"GroupName", "Comment", "URN"})

		for _, g := range groups {
			tw.AppendRow(table.Row{g.Name, g.Comment, g.ID})
		}

		tw.SetStyle(table.StyleLight)
		tw.Style().Options.SeparateColumns = false
		tw.Style().Options.DrawBorder = false

		fmt.Println(tw.Render())
	},
}

var getGroupCommand = &cobra.Command{
	Use:     "get",
	Short:   "Display details information and user memberships from a single group.",
	Aliases: []string{"load", "show"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		gc := iamClient.Groups()

		urn := iam.GroupURN(args[0])
		if urn.GroupName() == "" {
			urn = iam.GroupURN("urn:iam::group/" + urn)
		}

		grp, err := gc.Load(context.Background(), urn)
		if err != nil {
			log.Fatal(err)
		}

		members, err := gc.GetMembers(context.Background(), urn)
		if err != nil {
			log.Fatal(err)
		}

		blob, err := yaml.Marshal(struct {
			iam.Group
			Members []iam.UserURN `json:"members"`
		}{
			Group:   grp,
			Members: members,
		})

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(blob))
	},
}

var deleteGroupCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete a group.",
	Run: func(cmd *cobra.Command, args []string) {
		gc := iamClient.Groups()

		urn := iam.GroupURN(args[0])
		if urn.GroupName() == "" {
			urn = iam.GroupURN("urn:iam::group/" + urn)
		}

		if err := gc.Delete(context.Background(), urn); err != nil {
			log.Fatal(err)
		}
	},
}

var setCommentCommand = &cobra.Command{
	Use:     "set-comment",
	Short:   "Update a groups comment.",
	Aliases: []string{"comment"},
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		gc := iamClient.Groups()

		urn := iam.GroupURN(args[0])
		if urn.GroupName() == "" {
			urn = iam.GroupURN("urn:iam::group/" + urn)
		}

		comment := args[1]

		if err := gc.UpdateComment(context.Background(), urn, comment); err != nil {
			log.Fatal(err)
		}
	},
}

var createGroupCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a new group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		comment, _ := cmd.Flags().GetString("comment")

		gc := iamClient.Groups()

		urn, err := gc.Create(context.Background(), name, comment)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(urn)
	},
}

var addMemberCommand = &cobra.Command{
	Use:     "add-member",
	Short:   "Add a new user to a group",
	Aliases: []string{"add"},
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		user, _ := cmd.Flags().GetString("user")
		group, _ := cmd.Flags().GetString("group")

		userURN := iam.UserURN(user)
		groupURN := iam.GroupURN(group)

		if !userURN.IsValid() {
			userURN = iam.UserURN("urn:iam::user/" + user)
		}

		if !groupURN.IsValid() {
			groupURN = iam.GroupURN("urn:iam::group/" + group)
		}

		gc := iamClient.Groups()

		if err := gc.AddMember(context.Background(), groupURN, userURN); err != nil {
			log.Fatal(err)
		}
	},
}

var deleteMemberCommand = &cobra.Command{
	Use:     "delete-member",
	Short:   "Delete a user from a group",
	Aliases: []string{"remove"},
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		user, _ := cmd.Flags().GetString("user")
		group, _ := cmd.Flags().GetString("group")

		userURN := iam.UserURN(user)
		groupURN := iam.GroupURN(group)

		if !userURN.IsValid() {
			userURN = iam.UserURN("urn:iam::user/" + user)
		}

		if !groupURN.IsValid() {
			groupURN = iam.GroupURN("urn:iam::group/" + group)
		}

		gc := iamClient.Groups()

		if err := gc.DeleteMember(context.Background(), groupURN, userURN); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCommand.AddCommand(groupRootCommand)

	createGroupCommand.Flags().StringP("comment", "c", "", "Comment for the new group")

	addMemberCommand.Flags().StringP("user", "u", "", "Username to add to the group.")
	addMemberCommand.Flags().StringP("group", "g", "", "The target group.")
	addMemberCommand.MarkFlagRequired("user")
	addMemberCommand.MarkFlagRequired("group")

	deleteMemberCommand.Flags().StringP("user", "u", "", "Username to remove from the group.")
	deleteMemberCommand.Flags().StringP("group", "g", "", "The target group.")
	deleteMemberCommand.MarkFlagRequired("user")
	deleteMemberCommand.MarkFlagRequired("group")

	groupRootCommand.AddCommand(
		listGroupsCommand,
		getGroupCommand,
		deleteGroupCommand,
		setCommentCommand,
		createGroupCommand,
		addMemberCommand,
		deleteMemberCommand,
	)
}
