package cmd

import (
	"log"

	"github.com/acuas/YouOffScraper/cmd/instance"
	"github.com/spf13/cobra"
)

var cfgFile string

// scrapeCmd represents the scrape command
var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "TODO",
	Long:  `TODO`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Fatal(err)
		}
		err = instance.Run(cfg, args[0])
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(scrapeCmd)

	scrapeCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.yaml", "config file (default is $HOME/config.yaml")
}
