/*
Copyright © 2023 Chandler <chandler@chand1012.dev>
*/
package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"github.com/chand1012/arail/pkg/config"
)

var (
	AIModel string
	APIKey  string
	FromEnv bool
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:     "config",
	Short:   "Set configuration for Arail",
	Long:    `Set configuration for Arail. This includes the API key for OpenAI API and the model to use for ChatGPT.`,
	Aliases: []string{"c", "conf"},
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.Load()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		if FromEnv {
			log.Info("Getting API and Model from environment variables...")
			if AIModel == "" {
				AIModel = os.Getenv("OPENAI_MODEL")
			}
			if APIKey == "" {
				APIKey = os.Getenv("OPENAI_API_KEY")
			}
		}

		if AIModel != "" {
			log.Info("Setting model to '" + AIModel + "'...")
			conf.Model = AIModel
		}

		if os.Getenv("OPENAI_API_KEY") != "" && !FromEnv {
			log.Warn("OPENAI_API_KEY environment variable is already set.")
		}

		if APIKey != "" {
			log.Info("Setting API key...")
			conf.APIKey = APIKey
		}

		err = config.Save(conf)

		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().StringVarP(&AIModel, "model", "m", "", "The model to use for ChatGPT")
	configCmd.Flags().StringVarP(&APIKey, "apikey", "k", "", "The API key to use for OpenAI API")
	configCmd.Flags().BoolVarP(&FromEnv, "from-env", "e", false, "Use the OPENAI_API_KEY and OPENAI_MODEL environment variables")
}
