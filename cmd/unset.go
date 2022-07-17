/*
	Copyright © 2022 Adam Cłapiński <clapinskiadam@gmail.com>
*/

package cmd

import (
	"fmt"
	"os"
	"proxier/internal/proxier"
	"proxier/pkg/logger"

	"github.com/spf13/cobra"
)

// unsetCmd represents the unset command
var unsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "Configure files to NOT USE proxy",
	Long:  `Remove lines from listed files in config.json`,
	Run: func(cmd *cobra.Command, args []string) {
		var l logger.Logger = logger.NewDummy()
		if v, err := cmd.Flags().GetBool("verbose"); err != nil {
			panic(err)
		} else if v {
			l = logger.NewStandard()
		}

		p, err := proxier.New(l)
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
