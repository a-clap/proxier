/*
	Copyright © 2022 Adam Cłapiński <clapinskiadam@gmail.com>
*/

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"proxier/internal/proxier"
	"proxier/pkg/logger"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Configure files to use proxy",
	Long:  `Append lines to listed files in config.json`,
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
		fmt.Println("Appending lines...")
		if err = p.Set(true); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to set %v\n", err)
			os.Exit(1)
		}
		fmt.Println("All good!")

		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
