[panco.dev](https://panco.dev) \| [Documentation Home](https://panco.dev/docs.html) \| [Objects Command](https://panco.dev/objects.html)

# Generate CLI Commands from a CSV File

```
Usage:
  panco objects cli [flags]

Flags:
  -f, --file string     Name of the CSV file to convert
  -h, --help            help for cli
  -o, --output string   Name of the file to output SET commands to
  ```

## Overview

Using the cli command allows you to generate CLI set commands from the CSV file that can be pasted into the device. For each type of action you want to perform,
you will need to make sure your CSV file is structured accordingly. Please use the following guide for this purpose:

[CSV Structure for Object Actions](https://panco.dev/csv_objects.html)

Once your CSV file is all set, you can execute the following command to generate a text file with the commands:

```
panco objects cli --file <name-of-CSV-file> --output <name-of-TXT-file>
```
