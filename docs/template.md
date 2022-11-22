# Template Command

```
Usage:
  panco template [flags]

Flags:
  -h, --help          help for template
  -t, --type string   Type of the template to generate <objects|policy|all> (default "all")
```

## Overview

The `template` command generates CSV files which can be used as a guide (template) when importing objects and policies. The CSV file(s) will be generated in the location where you are running the `panco` command from.

For more information regarding the CSV file structure, please see the following:

* [CSV Structure - Objects](https://panco.dev/csv_objects.html)
* [CSV Structure - Policies](https://panco.dev/csv_policy.html)
