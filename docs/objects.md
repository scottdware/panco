# Objects Command

```
Usage:
  panco objects [flags]
  panco objects [command]

Available Commands:
  cli         Generate CLI set commands from a CSV file
  duplicates  Find duplicate address and service objects
  export      Export address, service and tag objects
  import      Import (create, modify) address, service, custom URL and tag objects

Flags:
  -h, --help   help for objects

Use "panco objects [command] --help" for more information about a command.
```

## Overview

The `objects` command allows you to import (create, modify) address, service and tag objects, along with exporting them from
the device. You can also find duplicate address and service objects, which will output to an Excel file for you to review.

You can also add/remove URL entries from a custom URL category object.

Using the `cli` command, you can convert the CSV object entries to CLI set commands that can be pasted into the device.

**_Important_**: Please refer to the [CSV Structure](https://panco.dev/csv_objects.html) page on how to structure your CSV files when importing/modifying objects.

Click on any one of the available commands to view the full documentation and usage:

* [cli](https://panco.dev/objects_cli.html)
* [duplicates](https://panco.dev/objects_duplicates.html)
* [export](https://panco.dev/objects_export.html)
* [import](https://panco.dev/objects_import.html)
