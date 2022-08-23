/*
	Copyright © 2022 Adam Cłapiński <clapinskiadam@gmail.com>
*/

package cmd

import (
	"fmt"
	"github.com/a-clap/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"
	"os"
	"proxier/internal/proxier"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Configure files to use proxy",
	Long:  `Append lines to listed files in config.json`,
	Run: func(cmd *cobra.Command, args []string) {
		lvl := zapcore.ErrorLevel
		if v, err := cmd.Flags().GetBool("verbose"); err != nil {
			panic(err)
		} else if v {
			lvl = zapcore.DebugLevel
		}

		logger.Init(logger.NewDefaultZap(lvl))

		p, err := proxier.New()
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
