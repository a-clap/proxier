/*
	Copyright © 2022 Adam Cłapiński <clapinskiadam@gmail.com>
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// unsetCmd represents the unset command
var unsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "Configure files to NOT USE proxy",
	Long:  `Remove lines from listed files in config.json`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("unset called")
	},
}

func init() {
	rootCmd.AddCommand(unsetCmd)
}
