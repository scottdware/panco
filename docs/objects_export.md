[panco.dev](https://panco.dev) \| [Documentation Home](https://panco.dev/docs.html) \| [Objects Command](https://panco.dev/objects.html)

# Exporting Objects

```
Usage:
  panco objects export [flags]

Flags:
  -p, --delay string         Delay (in milliseconds) to pause between each API call (default "100")
  -d, --device string        Device to connect to
  -g, --devicegroup string   Device Group name (default "shared")
  -f, --file string          Name of the CSV file you'd like to export to (default "PaloAltoObjects")
  -h, --help                 help for export
  -t, --type string          <address|addressgroup|service|servicegroup|tags|all>
  -u, --user string          User to connect to the device as
  -v, --vsys string          Vsys name (default "vsys1")
  ```

## Overview

This command allows you to export all of the addres, service and tag objects on a device. You can choose to specify
certain objects, or all of them. Each type will have it's own CSV file named as such:

> <_filename_> is what you specify with the `--file` flag.

* <_filename_>_Address.csv for address objects
* <_filename_>_AddressGroup.csv for address group objects
* <_filename_>_Service.csv for service objects
* <_filename_>_ServiceGroup.csv for service group objects
* <_filename_>_Tags.csv for tag objects

Run the following command to export all objects from the device:

```
panco objects export --device 10.1.1.1 --user admin --type all --file My-Objects.csv
```
