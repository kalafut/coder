package cli

import (
	"fmt"

	"github.com/coder/coder/cli/cliui"
	"github.com/coder/coder/codersdk"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

func userUpdate() *cobra.Command {
	var (
		email    string
		username string
	)
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a user's username and/or email",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := CreateClient(cmd)
			if err != nil {
				return err
			}

			identifier := args[0]
			if identifier == "" {
				return xerrors.Errorf("user identifier cannot be an empty string")
			}

			user, err := client.User(cmd.Context(), identifier)
			if err != nil {
				return xerrors.Errorf("fetch user: %w", err)
			}

			// Only prompt if no flags were provided
			if username == "" && email == "" {
				username, err = cliui.Prompt(cmd, cliui.PromptOptions{
					Text:    "Username:",
					Default: user.Username,
				})
				if err != nil {
					return err
				}

				email, err = cliui.Prompt(cmd, cliui.PromptOptions{
					Text:    "Email:",
					Default: user.Email,
					Validate: func(s string) error {
						err := validator.New().Var(s, "email")
						if err != nil {
							return xerrors.New("That's not a valid email address!")
						}
						return err
					},
				})
				if err != nil {
					return err
				}
			}

			if username == "" {
				username = user.Username
			}

			if email == "" {
				email = user.Email
			}

			updatedUser, err := client.UpdateUser(cmd.Context(), user.ID.String(), codersdk.UpdateUserRequest{
				Email:    email,
				Username: username,
			})
			if err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), `
User has been updated!

The new user details are:
    Username: %s
    Email:    %s

`, cliui.Styles.Keyword.Render(updatedUser.Username), cliui.Styles.Keyword.Render(updatedUser.Email))
			return nil
		},
	}
	cmd.Flags().StringVarP(&email, "email", "e", "", "Specifies the new email address for the user.")
	cmd.Flags().StringVarP(&username, "username", "u", "", "Specifies the new username for the user.")
	return cmd
}
