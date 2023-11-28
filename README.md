# tail2sentinel

A Go program that exports Tailscale network logs and events to Microsoft Sentinel SIEM.
Two tables are used; `TailscaleAudit` for audit logs and `TailscaleNetwork` for network logs.

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
  
  audit_output:
      resource_group: ""
      workspace_name: ""
    
      dcr:
        endpoint: ""
        rule_id: ""
        stream_name: ""
    
      expires_months: 6
      update_table: false
      
    network_output:
      resource_group: ""
      workspace_name: ""

      dcr:
        endpoint: ""
        rule_id: ""
        stream_name: ""

      expires_months: 6
      update_table: false

tailscale:
  tailnet: ""
  client_id: ""
  client_secret: ""
  lookback_days: 30
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
