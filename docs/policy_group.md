# Group Rules by Tag

```
Usage:
  panco policy group [flags]

Flags:
  -d, --device string        Device to connect to
  -g, --devicegroup string   Device Group name (default "shared")
  -f, --file string          Name of the CSV file
  -h, --help                 help for group
  -p, --pass string          Password for the user account specified
  -t, --type string          <security|nat>
  -u, --user string          User to connect to the device as
  -v, --vsys string          Vsys name (default "vsys1")
```

## Overview

This command allows you to group rules by a tag. This is helpful when looking at the policy and being able to filter
down on a select group of rules for the purpose defined in the tag specified (e.g. "Trust-to-Internet" rules). Currently,
you can **_ONLY_** group security and NAT rules by tags at this time.

Please use the below link as a guide on how to structure your CSV file when grouping rules:

[CSV Structure - Policies (Group Rules By Tag)](https://panco.dev/csv_policy.html#group-rules-by-tags)

Once your CSV file structure is all set, you can apply the changes by running the following command:

```
panco policy group --file <name-of-CSV-file> --type <security|nat>
```