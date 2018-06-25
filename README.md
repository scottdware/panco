# panco
Command-line tool that interacts with Palo Alto firewalls and Panorama.

Primarily, this tool is used for creating objects and security rules via CSV files, as
well as exporting the above information, as well as:

* Log exporting: Query any log type just like you would in the GUI, and export them to a CSV file locally on your machine.
* Session table dump: Query the entire session table on a firewall, and export it to a CSV file.
  * You can also use filters just as you would on the command line.

All of the backend functions are from the [go-panos](https://github.com/scottdware/go-panos) package.

More features will continue to be added on a regular basis.

For a detailed explanation of commands, and how they are used, click on any one of the
command names below.

`panco` [`help`](https://github.com/scottdware/panco#usage), [`example`][example-doc], [`logs`][logs-doc], [`objects`][objects-doc], [`policy`][policy-doc], [`sessions`][sessions-doc], [`version`][version-doc]

## Installation

Installation of extremely easy, and just involves downloading the correct binary for your OS. You can download them by clicking on the [release](https://github.com/scottdware/panco/releases) tab above.

Current support OS's:

* Windows
* Mac OS
* Linux

## Usage

```
Usage:
  panco [command]

Available Commands:
  example     Create example CSV files for import reference
  help        Help about any command
  import      Import CSV files that will create and/or modify objects
  logs        Retrieve logs from the device and export them to a CSV file
  sessions    Query the session table on a firewall, and export it to a CSV file
  version     Prints the version number of panco

Flags:
  -h, --help   help for panco

Use "panco [command] --help" for more information about a command.
```

## panco example

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

The sections below describe these files and how to structure them in more detail.

### example-create.csv

The CSV file for object creation should be organized with the following columns:

`name,type,value,description (optional),tag (optional),device-group`.

> **_<span style="color:red">NOTE</span>_**: Here are a few things to keep in mind when creating objects:
> * For the name of the object, it cannot be longer than 63 characters, and must only include letters, numbers, spaces, hyphens, and underscores.
> * If you are tagging an object upon creation, please make sure that the tags exist prior to creating the objects.
> * When creating service groups, you DO NOT need to specify a description, as they do not have that capability.
> * When you create address or service groups, I would place them at the bottom of the CSV file, that way you don't risk adding a member that doesn't exist.
> * When creating objects on a local firewall, and not Panorama, you can leave the device-group column blank.

**Creating Address Objects**

Column | Description
:--- | :---
`name` | Name of the object you wish to create.
`type` | **ip**, **range**, or **fqdn**
`value` | Must contain the IP address, FQDN, or IP range of the object.
`description` | (Optional) A description of the object.
`tag` | (Optional) Name of a pre-existing tag on the device to apply.
`device-group` | Name of the device-group, or **shared** if creating a shared object.

When creating address groups:

Column | Description
:--- | :---
`name` | Name of the address group you wish to create.
`type` | **static** or **dynamic**
`value` | * See below explanation
`description` | (Optional) A description of the object.
`tag` | (Optional) Name of a pre-existing tag on the device to apply.
`device-group` | Name of the device-group, or **shared** if creating a shared object.

For a **_static_** address group, `value` must contain a list of members to add to the group, separated by a space, e.g.:

`ip-host1 ip-net1 fqdn-example.com`

For a **_dynamic_** address group, `value` must contain the criteria (tags) to match on, e.g.:

`web-servers or db-servers and linux`

**Creating Service Objects**

Column | Description
:--- | :---
`name` | Name of the object you wish to create.
`type` | **tcp** or **udp**
`value` | * See below
`description` | (Optional) A description of the object.
`tag` | (Optional) Name of a pre-existing tag on the device to apply.
`device-group` | Name of the device-group, or **shared** if creating a shared object.

* `value` must contain a single port number, range (1023-3000), or comma-separated list of ports, e.g.: `80, 443, 2000`.

When creating service groups:

Column | Description
:--- | :---
`name` | Name of the object you wish to create.
`type` | **service**
`value` | * See below
`description` | Not available on service groups.
`tag` | (Optional) Name of a pre-existing tag on the device to apply.
`device-group` | Name of the device-group, or **shared** if creating a shared object.

* `value` must contain a list of service objects to add to the group, separated by a space, e.g.: `tcp_8080 udp_666 tcp_range`.

Example:

![alt-text](https://raw.githubusercontent.com/scottdware/images/master/example-create.png "example-create.csv")

### example-modify.csv

The CSV file for modifying groups should be organized with the following columns:

`grouptype,action,object-name,group-name,device-group`.

Column | Description
:--- | :---
`grouptype` | **address** or **service**
`action` | **add** or **remove**
`object-name` | Name of the object to add or remove from group.
`group-name` | Name of the group to modify.
`device-group` | Name of the device-group, or **shared** if creating a shared object.

Example:

![alt-text](https://raw.githubusercontent.com/scottdware/images/master/example-modify.png "example-modify.csv")

## panco logs [flags]

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

> **NOTE:** If you do not get any results, you might want to try using the `--wait` flag and increasing the delay time.

[Here](https://github.com/scottdware/panco/blob/master/traffic_log_example.csv) is an example of an export of traffic logs.

## panco objects [flags]

```
Usage:
  panco objects [flags]

Flags:
  -a, --action string        Action to perform - export, import, or modify
  -d, --device string        Firewall or Panorama device to connect to
  -g, --devicegroup string   Device group - only needed when ran against Panorama
  -f, --file string          Name of the CSV file to export/import or modify
  -h, --help                 help for objects
  -u, --user string          User to connect to the device as
```

This command allows you to perform the following actions on address and service objects:
export, import, and modify groups. When you select the export option (`--action export`), there are
two files that will be created. One will old all of the address objects, and the other will hold all of the service objects.

Using the modify action, allows you to add or remove objects from groups, based on the data you have within your CSV file.

Please see the [`example`][example-doc] command documentation above on how the CSV files should be structured.

## panco policy [flags]

```
Usage:
  panco policy [flags]

Flags:
  -a, --action string        Action to perform - export or import
  -d, --device string        Firewall or Panorama device to connect to
  -g, --devicegroup string   Device group - only needed when ran against Panorama
  -f, --file string          Name of the CSV file to export/import
  -h, --help                 help for policy
  -u, --user string          User to connect to the device as
```

This command will allow you to export and import an entire security policy. If you are running this
against a Panorama device, it can be really helpful if you want to clone an entire policy,
as you can export it from one device-group, modify it if needed, then import the poilcy into a different device-group.

You must always specify the action you want to take via the --action flag. Actions are either export or import.

## panco sessions [flags]

```
Usage:
  panco sessions [flags]

Flags:
  -d, --device string   Firewall or Panorama device to connect to
  -e, --export string   Name of the CSV file to export the session table to
  -f, --filter string   Filter string to include sessions only matching the criteria
  -h, --help            help for sessions
  -u, --user string     User to connect to the device as
```

This command will dump the entire session table on a firewall to the CSV file
that you specify. You can optionally define a filter, and use the same criteria as you would
on the command line. The filter query must be enclosed in quotes "", and the format is:

option=value (e.g. `--filter "application=ssl"`)

Your filter can include multiple items, and each group must be separated by a comma, e.g.:

`--filter "application=ssl, ssl-decrypt=yes, protocol=tcp"`

Depending on the number of sessions, the export could take some time.

## panco version

```
Usage:
  panco version [flags]

Flags:
  -h, --help   help for version
```

Version information for panco.

[example-doc]: https://github.com/scottdware/panco#panco-example
[import-doc]: https://github.com/scottdware/panco#panco-import-flags
[logs-doc]: https://github.com/scottdware/panco#panco-logs-flags
[sessions-doc]: https://github.com/scottdware/panco#panco-sessions-flags
[policy-doc]: https://github.com/scottdware/panco#panco-policy-flags
[objects-doc]: https://github.com/scottdware/panco#panco-objects-flags
[version-doc]: https://github.com/scottdware/panco#panco-version