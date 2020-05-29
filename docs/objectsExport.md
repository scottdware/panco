# Exporting Objects

```
Usage:
  panco objects export [flags]

Flags:
  -d, --device string        Device to connect to
  -g, --devicegroup string   Device Group name (default "shared")
  -f, --file string          Name of the CSV file you'd like to export to (default "PaloAltoObjects")
  -h, --help                 help for export
  -p, --pass string          Password for the user account specified
  -t, --type string          <address|addressgroup|service|servicegroup|tags|all>
  -u, --user string          User to connect to the device as
  -v, --vsys string          Vsys name (default "vsys1")
  ```

## Overview

This command allows you to export all of the addres, service and tag objects on a device. You can choose to specify
certain objects, or all of them. Each type will have it's own CSV file named as such:

> <_filename_> is what you specify with the `--file` flag.

* <_filename_>-Addr.csv for address objects
* <_filename_>-Addrgrp.csv for address group objects
* <_filename_>-Srvc.csv for service objects
* <_filename_>-Srvcgrp.csv for service group objects
* <_filename_>-Tags.csv for tag objects

Run the following command to export all objects from the device:

```
panco objects export --type all --file My-Objects.csv
```