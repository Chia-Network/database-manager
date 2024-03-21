package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/chia-network/database-manager/internal/config"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate ensures that the config is valid and all vars exist",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(viper.GetString("config"))
		if err != nil {
			log.Fatalln(err.Error())
		}
		err = cfg.Validate()
		if err != nil {
			log.Fatalln(err.Error())
		}

		log.Println("Done!")
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
