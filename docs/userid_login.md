# Manually Logging In/Out Users

```
Usage:
  panco userid login [flags]

Flags:
  -d, --device string   Device to connect to
  -f, --file string     Name of the CSV file to import
  -h, --help            help for login
  -p, --pass string     Password for the user account specified
  -u, --user string     User to connect to the device as
  -v, --vsys string     Vsys name (default "vsys1")
```

## Overview

Using the `login` command allows you to manually register/map a user-ID to an IP address. You can also unregister them
to release the IP address mapping. This can be helpful if you have no way of reading user-ID logs, say from a domain controller,
but still want to write rules around users, not just IP addresses.

* Currently, you can **_ONLY_** login/logout a user on a firewall, not in Panorama.

Please refer to the below link as a guide on how to format your CSV file:

[CSV Structure - Userid Functions](https://panco.dev/csv_userid.html)