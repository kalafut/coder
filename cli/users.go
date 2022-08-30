package cli

import (
	"github.com/spf13/cobra"

	"github.com/coder/coder/codersdk"
)

func users() *cobra.Command {
	cmd := &cobra.Command{
		Short:   "Create, update, remove, and list users",
		Use:     "users",
		Aliases: []string{"user"},
	}
	cmd.AddCommand(
		userCreate(),
		userList(),
		userSingle(),
		userUpdate(),
		createUserStatusCommand(codersdk.UserStatusActive),
		createUserStatusCommand(codersdk.UserStatusSuspended),
	)
	return cmd
}
