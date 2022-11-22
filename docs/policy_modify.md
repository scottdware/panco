# Modifying/Editing A Policy

```
Modify (edit) a security, NAT or Policy-Based Forwarding policy

Usage:
  panco policy modify [flags]

Flags:
  -d, --device string        Device to connect to
  -g, --devicegroup string   Device Group name when importing to Panorama (default "shared")
  -f, --file string          Name of the CSV file to export to
  -h, --help                 help for modify
  -l, --location string      Location of the rulebase - <pre|post> (default "pre")
  -t, --type string          Type of policy to import - <security|nat|pbf>
  -u, --user string          User to connect to the device as
  -v, --vsys string          Vsys name when importing to a firewall (default "vsys1")
```

## Overview

Using the `modify` command allows you to modify (or edit) existing rules, by adding or removing entries from
each of the rule fields. You can modify/edit the following types of policies at this time:

* Security
* NAT
* Policy-Based Forwarding (PBF)

Please use the below link as a guide on how to structure your CSV file when modifying rules:

[CSV Structure - Policies](https://panco.dev/csv_policy.html)
