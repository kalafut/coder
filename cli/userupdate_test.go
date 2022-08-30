package cli_test

import (
	"context"
	"testing"

	"github.com/coder/coder/cli/clitest"
	"github.com/coder/coder/coderd/coderdtest"
	"github.com/coder/coder/codersdk"
	"github.com/stretchr/testify/require"
)

func TestUserUpdate(t *testing.T) {
	t.Parallel()
	t.Run("successful updates", func(t *testing.T) {
		client := coderdtest.New(t, nil)
		admin := coderdtest.CreateFirstUser(t, client)
		other := coderdtest.CreateAnotherUser(t, client, admin.OrganizationID)
		targetUser, err := other.User(context.Background(), codersdk.Me)
		require.NoError(t, err, "fetch user")

		t.Run("username", func(t *testing.T) {
			newUsername := "newUsername"
			cmd, root := clitest.New(t, "users", "update", targetUser.Username, "--username", newUsername)
			clitest.SetupConfig(t, client, root)
			err := cmd.Execute()
			require.NoError(t, err)

			updatedUser, err := client.User(context.Background(), newUsername)
			require.NoError(t, err, "fetch updated user")
			require.Equal(t, newUsername, updatedUser.Username, "updated user username")
			require.Equal(t, targetUser.Email, updatedUser.Email, "updated user email")

			// Point target to the updated user
			targetUser = updatedUser
		})

		t.Run("email", func(t *testing.T) {
			newEmail := "bob@example.com"
			cmd, root := clitest.New(t, "users", "update", targetUser.Username, "--email", newEmail)
			clitest.SetupConfig(t, client, root)
			err := cmd.Execute()
			require.NoError(t, err)

			updatedUser, err := client.User(context.Background(), targetUser.Username)
			require.NoError(t, err, "fetch updated user")
			require.Equal(t, targetUser.Username, updatedUser.Username, "updated user username")
			require.Equal(t, newEmail, updatedUser.Email, "updated user email")

			// Point target to the updated user
			targetUser = updatedUser
		})

		t.Run("username and email", func(t *testing.T) {
			newUsername := "anotherNewUsername"
			newEmail := "bob@example.com"
			cmd, root := clitest.New(t, "users", "update", targetUser.Username, "--username", newUsername, "--email", newEmail)
			clitest.SetupConfig(t, client, root)
			err := cmd.Execute()
			require.NoError(t, err)

			updatedUser, err := client.User(context.Background(), newUsername)
			require.NoError(t, err, "fetch updated user")
			require.Equal(t, newUsername, updatedUser.Username, "updated user username")
			require.Equal(t, newEmail, updatedUser.Email, "updated user email")

			// Point target to the updated user
			targetUser = updatedUser
		})
	})

	t.Run("failed updates", func(t *testing.T) {
		client := coderdtest.New(t, nil)
		admin := coderdtest.CreateFirstUser(t, client)
		targetUser, err := coderdtest.CreateAnotherUser(t, client, admin.OrganizationID).User(context.Background(), codersdk.Me)
		require.NoError(t, err, "fetch user")

		conflictUser, err := coderdtest.CreateAnotherUser(t, client, admin.OrganizationID).User(context.Background(), codersdk.Me)
		require.NoError(t, err, "fetch user")

		t.Run("missing user", func(t *testing.T) {
			cmd, root := clitest.New(t, "users", "update", "notauser32408", "--username", "wontmatter")
			clitest.SetupConfig(t, client, root)
			err := cmd.Execute()
			require.ErrorContains(t, err, `"user" must be an existing uuid or username`)
		})

		t.Run("conflicting username", func(t *testing.T) {
			newUsername := conflictUser.Username
			cmd, root := clitest.New(t, "users", "update", targetUser.Username, "--username", newUsername)
			clitest.SetupConfig(t, client, root)
			err := cmd.Execute()
			require.ErrorContains(t, err, `User already exists`)
			require.ErrorContains(t, err, `username`)
		})

		t.Run("conflicting email", func(t *testing.T) {
			newEmail := conflictUser.Email
			cmd, root := clitest.New(t, "users", "update", targetUser.Username, "--email", newEmail)
			clitest.SetupConfig(t, client, root)
			err := cmd.Execute()
			require.ErrorContains(t, err, `User already exists`)
			require.ErrorContains(t, err, `email`)
		})

		t.Run("conflicting username and email", func(t *testing.T) {
			newUsername := conflictUser.Username
			newEmail := conflictUser.Email
			cmd, root := clitest.New(t, "users", "update", targetUser.Username, "--username", newUsername, "--email", newEmail)
			clitest.SetupConfig(t, client, root)
			err := cmd.Execute()
			require.ErrorContains(t, err, `User already exists`)
			require.ErrorContains(t, err, `username`)
			require.ErrorContains(t, err, `email`)
		})

		t.Run("insufficient privilege", func(t *testing.T) {
			t.Skip("not implemented")
		})
	})

	t.Run("prompt-based input", func(t *testing.T) {
		t.Skip("not implemented")
	})
}
