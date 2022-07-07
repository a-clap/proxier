/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"proxier/internal/config"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Allows modification of config.json file via cli",
	Long:  `Currently supports only creating template config`,
	Run: func(cmd *cobra.Command, args []string) {
		template, err := cmd.Flags().GetBool("template")
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error parsing template argument %v", err)
			os.Exit(1)
		}
		if template {
			filename := "template_config.json"
			fmt.Println("Generating", filename, "...")
			err := os.WriteFile(filename, config.Template(), 0755)
			if err != nil {
				log.Fatalln(err)
			} else {
				fmt.Println("Done!")
			}
		} else {
			fmt.Println("Nothing to do")
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	configCmd.Flags().BoolP("template", "t", false, "Generates template config")
}
