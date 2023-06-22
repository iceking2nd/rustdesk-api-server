/*
Copyright © 2023 Daniel Wu <wxc@wxccs.org>
*/
package cmd

import (
	"fmt"
	"github.com/iceking2nd/rustdesk-api-server/app/Models"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// dbMigrateCmd represents the dbMigrate command
var dbMigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate db schemas",
	/*Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,*/
	Run: func(cmd *cobra.Command, args []string) {
		db := Models.Init()
		if err := db.AutoMigrate(
			&Models.User{},
			&Models.Token{},
			&Models.Tag{},
			&Models.Address{},
		); err != nil {
			log.Fatalln(err.Error())
		} else {
			fmt.Println("Database schemas migrated.")
			os.Exit(0)
		}
	},
}

func init() {
	dbCmd.AddCommand(dbMigrateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dbMigrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dbMigrateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
