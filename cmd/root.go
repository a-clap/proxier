/*
	Copyright © 2022 Adam Cłapiński <clapinskiadam@gmail.com>
*/

package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "proxier",
	Short: "Application which helps make necessary changes in certain files to use/don't use proxy server",
	Long: `If you are lazy bastard, as I am, this application will enable/disable proxy in certain files in Linux.
 
What you need to do:
1. Create config.json (you can create template one with command config --template)
2. Call application with cmd 'set' (or 'unset').
By defaults application will create backup files in subdirectory backup/.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "make an application full of logs")
	rootCmd.PersistentFlags().BoolP("backup", "b", true, "make backup of files, which will be overridden")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
