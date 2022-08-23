/*
	Copyright © 2022 Adam Cłapiński <clapinskiadam@gmail.com>
*/

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"proxier/internal/proxier"
)

// unsetCmd represents the unset command
var unsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "Configure files to NOT USE proxy",
	Long:  `Remove lines from listed files in config.json`,
	Run: func(cmd *cobra.Command, args []string) {
		p, err := proxier.New()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to create proxier %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Removing lines...")
		if err = p.Unset(true); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to unset %v\n", err)
			os.Exit(1)
		}
		fmt.Println("All good!")
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(unsetCmd)
}
