package cmd

import (
	"context"
	"fmt"
	"github.com/charmbracelet/huh/spinner"
	"github.com/geowa4/servicelogger/pkg/internalservicelog"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var internalServiceLogCmd = &cobra.Command{
	Use:   "internal",
	Short: "Send an internal service log",
	Long: `Prompt for internal service log 

` + "Example: `servicelogger internal`",
	Run: func(cmd *cobra.Command, args []string) {
		desc, confirmation, err := internalservicelog.Program()
		cobra.CheckErr(err)
		fmt.Println(desc)
		if confirmation {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			go func() {
				defer cancel()
				// TODO send SL
				time.Sleep(3 * time.Second)
			}()
			err = spinner.New().Context(ctx).Run()
			cobra.CheckErr(err)
			_, _ = fmt.Fprint(os.Stderr, "service log sent")
		}
	},
}

func init() {
	rootCmd.AddCommand(internalServiceLogCmd)
}
