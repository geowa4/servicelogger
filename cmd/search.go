package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/geowa4/servicelogger/pkg/editor"
	"github.com/geowa4/servicelogger/pkg/search"
	"github.com/spf13/cobra"
	"os"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for a service log and fill in its variables",
	Long: `Run an interactive TUI to search and discover service log templates. The output of this command will be the service log with all its variables filled in based on user input. For example, take this template.

{
  "severity": "Warning",
  "service_name": "SREManualAction",
  "log_type": "cluster-configuration",
  "summary": "Action required: Pod(s) preventing Node Drain",
  "description": "Your cluster is attempting to drain a node but there are pod(s) preventing the drain. The SRE team has identified the pod(s) as '${POD}' running in namespace(s) '${NAMESPACE}'. Please re-schedule the impacted pod(s) so that the node can drain.",
  "internal_only": false,
  "doc_references": [
    "https://access.redhat.com/documentation/en-us/red_hat_openshift_service_on_aws/4/html/nodes/working-with-pods#nodes-pods-pod-distruption-about_nodes-pods-configuring"
  ],
  "_tags": [
    "sop_KubeNodeUnschedulableSRE",
    "sop_MCDDrainError"
  ]
}

The TUI will prompt you to fill in those variables, and the output might look like this.

{
  "severity": "Warning",
  "service_name": "SREManualAction",
  "log_type": "cluster-configuration",
  "summary": "Action required: Pod(s) preventing Node Drain",
  "description": "Your cluster is attempting to drain a node but there are pod(s) preventing the drain. The SRE team has identified the pod(s) as 'loki' running in namespace(s) 'logging-ns'. Please re-schedule the impacted pod(s) so that the node can drain.",
  "internal_only": false,
  "doc_references": [
    "https://access.redhat.com/documentation/en-us/red_hat_openshift_service_on_aws/4/html/nodes/working-with-pods#nodes-pods-pod-distruption-about_nodes-pods-configuring"
  ],
  "_tags": [
    "sop_KubeNodeUnschedulableSRE",
    "sop_MCDDrainError"
  ]
}
`,
	Run: func(cmd *cobra.Command, args []string) {
		template, err := search.Program()
		cobra.CheckErr(err)
		if template != nil {
			template = editor.Program(template)
		}
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
