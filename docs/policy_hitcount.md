# Exporting Policy Hit Count Data From Firewalls

```
Usage:
  panco policy hitcount [flags]

Flags:
  -d, --device string   Device to connect to
  -f, --file string     Name of the CSV file you'd like to export to (default "PaloAltoPolicy")
  -h, --help            help for hitcount
  -p, --pass string     Password for the user account specified
  -t, --type string     Type of policy to gather hit count on - <security|nat|pbf|all>
  -u, --user string     User to connect to the device as
  -v, --vsys string     Vsys name when exporting from a firewall (default "vsys1")
```

## Overview

The `hitcount` command allows you to export rule hit count on a firewall for the following policy types:

* Security
* NAT
* Policy-Based Forwarding (PBF)

Currently, this **_ONLY_** functions when ran against an individual firewall, and not Panorama.

You can choose to specify each individual rulebase, or all of them at once. Each policy will have it's own
file, named as such:

> <_filename_> is what you specify with the `--file` flag.

* <_filename_>-Security_HitCount.csv for security rules
* <_filename_>-NAT_HitCount.csv for NAT rules
* <_filename_>-PBF_HitCount.csv for PBF rules

To export hit count data from a firewall, execute the following command:

```
panco policy hitcount --file <name-of-Output-File> --type all
```

THe files will have the following columns defined:

`Name,Hit Count,First Hit,Last Hit,Last Reset,Rule Created,Rule Modified`