# Generate CLI Commands from a CSV File

```
Usage:
  panco objects cli [flags]

Flags:
  -f, --file string   Name of the CSV file to convert
  -h, --help          help for cli
  -t, --txt string    Name of the TXT file to output SET commands
  ```

## Overview

Using the cli command allows you to generate CLI set commands from the CSV file that can be pasted into the device. For each type of action you want to perform,
you will need to make sure your CSV file is structured accordingly. Please use the following guide for this purpose:

[CSV Structure for Object Actions](https://panco.dev/csvObjects.html)

Once your CSV file is all set, you can execute the following command to generate a text file with the commands:

```
panco objects cli --file <name-of-CSV-file> --txt <name-of-TXT-file>
```