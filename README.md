# panco
Command-line tool that interacts with Palo Alto firewalls and Panorama.

For a detailed explanation of commands, and how they are used, click on any one of the
command names below.

[`panco help`](https://github.com/scottdware/panco#usage)
[`panco example`][example-doc]
[`panco import`][import-doc]
[`panco logs`][logs-doc]

### Usage

```
Usage:
  panco [command]

Available Commands:
  example     Create example CSV files for import reference
  help        Help about any command
  import      Import CSV files that will create and/or modify objects
  logs        Retrieve logs from the device and export them to a CSV file

Flags:
  -h, --help   help for panco

Use "panco [command] --help" for more information about a command.
```

#### panco example

```
Usage:
  panco example [flags]

Flags:
  -h, --help   help for example
```

This command will create two sample, reference CSV files for use with the
`import` command. The files will be placed in the location where you are running
the command from, and are named as such:

`example-create.csv` and `example-modify.csv`

#### panco import [flags]

```
Usage:
  panco import [flags]

Flags:
  -c, --create string   Name of the CSV file to create objects with
  -d, --device string   Firewall or Panorama device to connect to
  -h, --help            help for import
  -m, --modify string   Name of the CSV file to modify groups with
  -u, --user string     User to connect to the device as
```

The `import` command, given the spcific flag, will create or modify address and/or
service objects based on the information you have provided in your CSV file(s).

#### panco logs [flags]

```
Usage:
  panco logs [flags]

Flags:
  -d, --device string   Firewall or Panorama device to connect to
  -e, --export string   Name of the CSV file to export the logs to
  -h, --help            help for logs
  -n, --nlogs int       Number of logs to retrieve (default 20)
  -q, --query string    Critera to search the logs on
  -t, --type string     Log type to search under (default "traffic")
  -u, --user string     User to connect to the device as
  -w, --wait int        Wait time in seconds to delay retrieving logs - helpful for large queries (default 5)
```

You can query the device logs via this tool the same way you would on the GUI.
The different log types you can retrieve are:

`config`, `system`, `traffic`, `threat`, `wildfire`, `url`, `data`

When using the `--query` flag, be sure to enclose your search criteria in quotes `""` like so:

`--query "(addr.src in 10.0.0.0/8)"`

The default search type is `traffic`. Based on your query, and the device,
log retrieval and export could take a while.

[example-doc]: https://github.com/scottdware/panco#panco-example
[import-doc]: https://github.com/scottdware/panco#panco-import-flags
[logs-doc]: https://github.com/scottdware/panco#panco-logs-flags