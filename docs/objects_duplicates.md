# Finding Duplicate Address and Service Objects

```
Find duplicate address and service objects

Usage:
  panco objects duplicates [flags]

Flags:
  -d, --device string        Device to connect to
  -g, --devicegroup string   Device Group name when exporting from Panorama (default "shared")
  -f, --file string          Name of the output file (you don't need an extension) (default "PaloAltoDuplicates")
  -h, --help                 help for duplicates
  -p, --pass string          Password for the user account specified
  -t, --type string          <address|service|all>
  -u, --user string          User to connect to the device as
  -v, --vsys string          Vsys name when exporting from a firewall (default "vsys1")
  ```

## Overview

Finding duplicate address and service objects is a quick and easy task with this command. You can choose either
address or service objects to find duplicates on, or both. When specifying either option, an Excel file is
created which includes the results in separate tabs. For example, to find _all_ duplicate address and service
objects, run the following command:

`panco objects duplicates --type all --file Dups --device 192.168.1.1 --user admin --password 'paloalto' --vsys vsys1`

What this will do is create a file called `Dups.xlsx` which will have four tabs:

* address-Unique
* address-Duplicates
* service-Unique
* service-Duplicates

If you were to just specify `--type address` for example, then you would just have the first two tabs listed above. Now
that we have our file, if you open it up it will look like the following: