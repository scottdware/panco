# Importing/Modifying A Policy

```
Usage:
  panco policy import [flags]

Flags:
  -d, --device string        Device to connect to
  -g, --devicegroup string   Device Group name when exporting from Panorama (default "shared")
  -f, --file string          Name of the CSV file to export to
  -h, --help                 help for import
  -p, --pass string          Password for the user account specified
  -t, --type string          Type of policy to import
  -u, --user string          User to connect to the device as
  -v, --vsys string          Vsys name when exporting from a firewall (default "vsys1")
```

## Overview

Using the `import` command allows you to create new rules, or modify existing rules by adding new values
to them. You can create/modify the following types of policies at this time:

* Security
* NAT
* Policy-Based Forwarding (PBF)

Please use the below link as a guide on how to structure your CSV file when importing rules:

[CSV Structure - Policies](https://panco.dev/csvPolicy.html)