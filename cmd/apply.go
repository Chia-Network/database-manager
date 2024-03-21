package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/chia-network/database-manager/internal/config"
	"github.com/chia-network/database-manager/internal/manager"
	"github.com/chia-network/database-manager/internal/mysql"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Applies the configuration to the database",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(viper.GetString("config"))
		if err != nil {
			log.Fatalln(err.Error())
		}
		err = cfg.Validate()
		if err != nil {
			log.Fatalln(err.Error())
		}

		mysqlM, err := mysql.NewMySQLManager(cfg.Connection.Username, cfg.Connection.Password, cfg.Connection.Host, cfg.Connection.Port)
		if err != nil {
			log.Fatalf("Error creating MySQL Manager: %s\n", err.Error())
		}
		mgr, err := manager.NewManager(mysqlM)
		if err != nil {
			log.Fatalf("Error creating manager: %s\n", err.Error())
		}
		err = mgr.Apply()
		if err != nil {
			log.Fatalf("Error during apply: %s\n", err.Error())
		}

		log.Println("Success!")
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// applyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// applyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
