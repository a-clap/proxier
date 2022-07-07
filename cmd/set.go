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

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Configure files to use proxy",
	Long:  `Append lines to listed files in config.json`,
	Run: func(cmd *cobra.Command, args []string) {
		log := logger.NewStandard()

		p, err := proxier.New(log)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to create proxier %v\n", err)
			os.Exit(1)
		}
		if err = p.Set(true); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to set %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
