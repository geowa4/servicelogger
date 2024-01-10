package cmd

import (
	"fmt"
	"github.com/geowa4/servicelogger/pkg/list"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Display service logs",
	Long: `Display a filterable list of service logs

` + "Example: `osdctl servicelog list $CLUSTER_ID | servicelogger list`",
	Run: func(cmd *cobra.Command, args []string) {
		slResponseBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "could not read stdin: %v", err)
			os.Exit(1)
		}
		list.Program(slResponseBytes)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
