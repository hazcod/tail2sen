# tail2sentinel

A Go program that exports Tailscale network logs and events to Microsoft Sentinel SIEM.

## Running

First create a yaml file, such as `config.yml`:
```yaml
log:
  level: INFO

microsoft:
  app_id: ""
  secret_key: ""
  tenant_id: ""
  subscription_id: ""
  resource_group: ""
  workspace_name: ""

  dcr:
    endpoint: ""
    rule_id: ""
    stream_name: ""

  expires_months: 6
  update_table: false

tailscale:
  api_token: ""
```

And now run the program from source code:
```shell
% make
go run ./cmd/... -config=dev.yml
INFO[0000] shipping logs                                 module=sentinel_logs table_name=TailscaleLogs total=82
INFO[0002] shipped logs                                  module=sentinel_logs table_name=TailscaleLogs
INFO[0002] successfully sent logs to sentinel            total=82
```

Or binary:
```shell
% tail2sen -config=config.yml
```

## Building

```shell
% make build
```
