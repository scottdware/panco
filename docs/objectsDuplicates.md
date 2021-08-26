# Finding Duplicate Address and Service Objects

```
Usage:
  panco objects duplicates [flags]

Flags:
  -d, --device string        Device to connect to
  -g, --devicegroup string   Device Group name when exporting from Panorama (default "shared")
  -f, --file string          Name of the output file (you don't need an extension) (default "PaloAltoDuplicates")
  -h, --help                 help for duplicates
  -t, --type string          <address|service|all>
  -u, --user string          User to connect to the device as
  -v, --vsys string          Vsys name when exporting from a firewall (default "vsys1")
  ```

## Overview

Finding duplicate address and service objects is a quick and easy task with this command. You can choose either
address or service objects to find duplicates on, or both. When specifying either option, an Excel file is
created which includes the results in separate tabs. For example, to find _all_ duplicate address and service
objects, run the following command:

```
panco objects duplicates --type all --file Duplicate-Objects --vsys vsys1
```

What this will do is create a file called `Duplicate-Objects.xlsx` which will have four tabs:

* address-Unique
* address-Duplicates
* service-Unique
* service-Duplicates

The duplicates are found based on the value of the object. So for example, if you have two objects with the same
IP address:

* **Desktop-PC** has the IP of **192.168.80.10**
* **Client-Machine** has the IP of **192.168.80.10**

Either one can show up in the duplicates tab. You'll just have to figure out which one you want to take action
on.

* Currently, this command just gives you a list of the duplicate objects, and there is no way to automatically
take action to resolve (delete) the dupllicates. This functionality will be coming in a later release.

<!-- If you were to just specify `--type address` for example, then you would just have the first two tabs listed above. Here
is an example of how this file looks:

<img src="duplicates.png" alt="Duplicate objects"/> -->