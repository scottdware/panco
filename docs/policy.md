# Policy Command

```
Usage:
  panco policy [flags]
  panco policy [command]

Available Commands:
  export      Export a security, NAT or Policy-Based Forwarding policy
  group       Group security or NAT rules by tags
  hitcount    Get the hit count data on a security, NAT or Policy-Based Forwarding policy - FIREWALL ONLY
  import      Import (create, modify) a security, NAT or Policy-Based Forwarding policy
  move        Move multiple rules within a security, NAT or Policy-Based Forwarding policy

Flags:
  -h, --help   help for policy

Use "panco policy [command] --help" for more information about a command.
```

## Overview

The `policy` command allows you to work with security, NAT or Policy-Based Forwarding policies. You will
be able to export a policy, import (create/modify) a policy, move multiple rules within different
rulebases, group security or NAT rules by tags, and more.

**_Important_**: Please refer to the [CSV Structure - Policies](https://panco.dev/csvPolicy.html) page
on how to structure your CSV files when importing, grouping or moving rules.

Click on any one of the available commands to view the full documentation and usage:

* [export](https://panco.dev/policyExport.html)
* [group](https://panco.dev/policyGroup.html)
* [hitcount](https://panco.dev/policyHitcount.html)
* [import](https://panco.dev/policyImport.html)
* [move](https://panco.dev/policyMove.html)