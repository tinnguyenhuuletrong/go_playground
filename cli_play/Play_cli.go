package cli_play

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/syndtr/goleveldb/leveldb"

	"ttin.com/play2022/utils"
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

var seedUsers []user = []user{
	*newUser("1", "a@email.com"),
	*newUser("2", "b@email.com"),
	*newUser("3", "c@email.com"),
	*newUser("4", "d@email.com"),
}

type datastore struct {
	db       *leveldb.DB
	allUsers []user
}

func newDatastore() *datastore {
	ins := &datastore{
		allUsers: seedUsers,
	}
	ins.open()
	return ins
}

func (ins *datastore) open() {
	db, err := leveldb.OpenFile("./tmp/cli-db", nil)
	if err != nil {
		panic(err)
	}
	ins.db = db
}
func (ins *datastore) close() {
	ins.db.Close()
}

func (ins *datastore) save() {
	userKey := []byte("users")
	ins.db.Put(userKey, utils.Dump2Gob(ins.allUsers), nil)
}

func (ins *datastore) loadData() {
	userKey := []byte("users")
	val, err := ins.db.Get(userKey, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(val) > 0 {
		ins.allUsers = utils.Gob2Obj[[]user](val)
	}
}

func Play_cli_cobra() {
	datastore := newDatastore()
	datastore.loadData()
	defer func() {
		datastore.save()
		datastore.close()
	}()

	var rootCmd = &cobra.Command{
		Use:   "app",
		Short: "My root command",
	}

	usersCmd := buildCmdUsers(datastore)

	rootCmd.AddCommand(usersCmd)
	rootCmd.Execute()
}

func buildCmdUsers(store *datastore) *cobra.Command {
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
				fmt.Println(utils.Dump2Json(store.allUsers))
			default:
				fmt.Println(store.allUsers)
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

			id := fmt.Sprintf("%d", len(store.allUsers)+1)
			newIns := *newUser(id, fEmail)
			store.allUsers = append(store.allUsers, newIns)

			fmt.Printf("added: %+v\n", newIns)
		},
	}
	userAdd.Flags().String("email", "", "email address")
	userAdd.MarkFlagRequired("email")

	usersCmd.AddCommand(userList)
	usersCmd.AddCommand(userAdd)
	return usersCmd
}
