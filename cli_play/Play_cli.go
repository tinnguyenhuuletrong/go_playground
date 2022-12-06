package cli_play

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

type user struct {
	Id    string `json:"id"`
	Email string `json:"name"`
}

func newUser(id, email string) *user {
	return &user{
		Id:    id,
		Email: email,
	}
}

var allUsers []user = []user{
	*newUser("1", "a@email.com"),
	*newUser("2", "b@email.com"),
	*newUser("3", "c@email.com"),
	*newUser("4", "d@email.com"),
}

func Play_cli_cobra() {
	var rootCmd = &cobra.Command{
		Use:   "app",
		Short: "My root command",
	}

	usersCmd := buildCmdUsers()

	rootCmd.AddCommand(usersCmd)
	rootCmd.Execute()
}

func buildCmdUsers() *cobra.Command {
	var usersCmd = &cobra.Command{
		Use:                   "users",
		Short:                 "User resources",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	var userList = &cobra.Command{
		Use:   "list",
		Short: "List resources",
		Run: func(cmd *cobra.Command, args []string) {
			fFormat := cmd.Flags().Lookup("format").Value.String()
			fmt.Printf("Inside userList Run with args: %v, flags: %v\n", args, fFormat)

			switch fFormat {
			case "json":
				fmt.Println(dump2Json(allUsers))
			default:
				fmt.Println(allUsers)
			}

		},
	}
	userList.Flags().String("format", "go", "output format json|go")

	var userAdd = &cobra.Command{
		Use:   "add",
		Short: "Add resources",
		Run: func(cmd *cobra.Command, args []string) {
			fEmail := cmd.Flags().Lookup("email").Value.String()
			fmt.Printf("Inside userAdd Run with args: %v, flags: %v\n", args, fEmail)

			id := fmt.Sprintf("%d", len(allUsers))
			newIns := *newUser(id, fEmail)
			allUsers = append(allUsers, newIns)

			fmt.Printf("added: %+v\n", newIns)
		},
	}
	userAdd.Flags().String("email", "", "email address")
	userAdd.MarkFlagRequired("email")

	usersCmd.AddCommand(userList)
	usersCmd.AddCommand(userAdd)
	return usersCmd
}

func dump2Json(v any) string {
	bytes, err := json.MarshalIndent(v, " ", "  ")
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
