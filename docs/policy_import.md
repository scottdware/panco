[panco.dev](https://panco.dev) \| [Documentation Home](https://panco.dev/docs.html) \| [Policy Command](https://panco.dev/policy.html)

# Importing/Modifying A Policy

```
Usage:
  panco policy import [flags]

Flags:
  -p, --delay string    Delay (in milliseconds) to pause between each API call (default "100")
  -d, --device string   Device to connect to
  -f, --file string     Name of the CSV file to export to
  -h, --help            help for import
  -t, --type string     Type of policy to import - <security|nat|pbf|decrypt>
  -u, --user string     User to connect to the device as
  -v, --vsys string     Vsys name when importing to a firewall (default "vsys1")
```

## Overview

Using the `import` command allows you to create new rules, or modify existing rules by adding new values
to them. You can create/modify the following types of policies at this time:

* Security
* NAT
* Policy-Based Forwarding (PBF)
* Decryption

Please use the below link as a guide on how to structure your CSV file when importing rules:

[CSV Structure - Policies](https://panco.dev/csv_policy.html)
