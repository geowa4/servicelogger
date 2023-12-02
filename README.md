# Service Logger

## Common Use

First, you'll need to download (or update) your cache of our SOPs and managed notifications.
```shell
go run main.go cache-update 
```

Then, you can run the search program and have it output the template JSON.
```shell
go run main.go search | jq .
```
