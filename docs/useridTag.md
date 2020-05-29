# Tagging IP Addresses

```
Usage:
  panco userid tag [flags]

Flags:
  -d, --device string   Device to connect to
  -f, --file string     Name of the CSV file to import
  -h, --help            help for tag
  -p, --pass string     Password for the user account specified
  -u, --user string     User to connect to the device as
  -v, --vsys string     Vsys name (default "vsys1")
```

## Overview

Using the `tag` command allows you to dynamically tag (untag) an IP address which can be used in a dynamic address group. This
can be extremely helpful for example, when blocking malicious IP addresses quickly. There is no commit required
on the firewall when tagging an IP address. You just have to make sure that the dynamic address group is created and the
configuration has already been committed.

* Currently, you can **_ONLY_** tag an IP address on a firewall, not in Panorama.
* Only single IP addresses are supported, no subnets or CIDR blocks.
* The IP addresses that you tag/untag do not have to pre-exist as objects.

Please refer to the below link as a guide on how to format your CSV file:

[CSV Structure - Userid Functions](https://panco.dev/csvUserid.html)