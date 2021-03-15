package cmd

import (
	"fmt"

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
		fmt.Println("scrape called")
	},
}

func init() {
	rootCmd.AddCommand(scrapeCmd)

	scrapeCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.yaml", "config file (default is $HOME/config.yaml")
}
