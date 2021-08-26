# Importing Objects

```
Usage:
  panco objects import [flags]

Flags:
  -d, --device string   Device to connect to
  -f, --file string     Name of the CSV file to import
  -h, --help            help for import
  -u, --user string     User to connect to the device as
  ```

## Overview

Using the import command allows you to create or modify address and service objects, tag objects, and more. For each type of action you want to perform,
you will need to make sure your CSV file is structured accordingly. Please use the following guide for this purpose:

[CSV Structure for Object Actions](https://panco.dev/csvObjects.html)

Once your CSV file is all set, you can execute the following command to apply the changes:

```
panco objects import --file <name-of-CSV-file>
```