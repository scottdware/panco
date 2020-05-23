# panco Documentation

```
Usage:
  panco [command]

Available Commands:
  config      Configure a device using IronSkillet or a local or remote (HTTP) file; export a device configuration
  help        Help about any command
  objects     Commands to work with address, service, and tag objects
  policy      Commands to work with security, NAT and PBF policies
  userid      Commands to interact with user-id functions
  version     Version information for panco

Flags:
  -h, --help   help for panco

Use "panco [command] --help" for more information about a command.
```

## Overview

This command-line tool helps to automate large tasks such as creating and modifying lots of objects,
tagging multiple objects, importing a policy (Security, NAT or PBF), changing values on rules,
finding duplicate address and service objects, and more. WHen it comes to creating and modifying
(importing) objects, all of this can be done by using a CSV file...and you can accomplish multiple tasks within the same file.

Please use the below link to the CSV structure as a guide when using the import function on the `objects`
, `poilcy` and `userid` commands.

* [CSV Structure](https://scottdware.github.io/panco/csv.html)

Click on any one of the available commands to view the full documentation and usage.

## Available Commands

* [config](config.html)
* [objects](objects.html)
* [policy](policy.html)
* [userid](userid.html)