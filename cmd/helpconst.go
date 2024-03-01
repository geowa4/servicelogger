package cmd

const commonOcmFlagLongHelpStanza = `This command requires two mostly static flgs. The first is "--ocm-url", which can be set as an environment variable "OCM_URL" or set in the global config file.

The second flag is "--ocm-token", which can be set as an environment variable "OCM_TOKEN" or  set in the global config file. It is recommended to use the output of the "ocm token" command.
`
const commonSendArgLongHelpStanza = commonOcmFlagLongHelpStanza + `This command also requires one or multiple cluster IDs to send. The most common case is to send to a single cluster. To do this, specify the "--cluster-id" flag or set the "CLUSTER_ID" environment variable. Another common case is to send the same service log to multiple clusters in parallel. For this case, use the "--cluster-ids" flag where the cluster IDs are comma-separated. The environment variable "CLUSTER_IDS" (note the pluralization) may be used but instead of being comma-separate, the cluster IDs must be separated by whitespace.
`
