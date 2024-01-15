# Service Logger

## Common Use

### Search

First, you'll need to download (or update) your cache of our SOPs and managed notifications.
```shell
servicelogger cache-update 
```

Then, you can run the search program and have it output the template JSON.
```shell
servicelogger search | jq .
```

### List View

```shell
osdctl servicelog list -A $CLUSTER_ID | servicelogger list
```