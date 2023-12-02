package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/geowa4/servicelogger/pkg/search"
	"github.com/spf13/cobra"
	"os"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for a service log",
	Long:  `Run an interactive TUI to search and discover service log templates`,
	Run: func(cmd *cobra.Command, args []string) {
		template := search.Program()
		//TODO: send template to parameter fill-in program
		templateJson, err := json.Marshal(template)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error printing selected template: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(templateJson))
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
