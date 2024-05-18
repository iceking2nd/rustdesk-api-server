/*
Copyright © 2023 Daniel Wu <wxc@wxccs.org>
*/
package cmd

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/iceking2nd/rustdesk-api-server/app/Models"
	"github.com/iceking2nd/rustdesk-api-server/utils/Hash"
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
			&Models.ActivateToken{},
			&Models.Tag{},
			&Models.Team{},
			&Models.User{},
			&Models.Group{},
			&Models.Token{},
			&Models.Address{},
			&Models.Settings{},
		); err != nil {
			log.Fatalln(err.Error())
		}

		var count int64 = -1
		err := db.Model(&Models.Team{}).Count(&count).Error
		if count <= 0 {
			db.Create(&Models.Team{
				GUID:  uuid.New().String(),
				Name:  "Default",
				EMail: "Default",
				Info:  "{}",
			})
			count = -1
		}
		if err != nil {
			log.Println(err.Error())
		} else {
			fmt.Println("Create default team successfully.")
		}

		err = db.Model(&Models.Group{}).Count(&count).Error
		if count <= 0 {
			var team Models.Team
			err = db.First(&team).Error
			if err != nil {
				log.Println(err.Error())
			} else {
				db.Create(&Models.Group{
					GUID:   uuid.New().String(),
					TeamID: team.ID,
					Name:   "Default",
					Note:   "Default Group",
					Info:   "{}",
				})
			}
			count = -1
		}
		if err != nil {
			log.Println(err.Error())
		} else {
			fmt.Println("Create default group successfully.")
		}

		err = db.Model(&Models.User{}).Count(&count).Error
		if count <= 0 {
			var group Models.Group
			err = db.First(&group).Error
			if err != nil {
				log.Println(err.Error())
			} else {
				pwd, _ := Hash.StringToSHA512("test1234")
				db.Create(&Models.User{
					GUID:     uuid.New().String(),
					Username: "admin@example.com",
					Password: pwd,
					Name:     "Administrator",
					GroupID:  group.ID,
					IsAdmin:  true,
					Status:   1,
					Note:     "Default Administrator account",
					Info:     "{}",
				})
			}
			count = -1
		}
		if err != nil {
			log.Println(err.Error())
		} else {
			fmt.Println("Create default administrator account successfully.")
		}

		err = db.Model(&Models.Settings{}).Count(&count).Error
		if count <= 0 {
			db.Create(&Models.Settings{
				Key:   "LICENSE",
				Value: "community license",
			})
			count = -1
		}

		if err != nil {
			log.Println(err.Error())
		} else {
			fmt.Println("Create default settings successfully.")
		}

		fmt.Println("Database schemas migrated.")
		os.Exit(0)

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
