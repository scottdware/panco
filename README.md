# panco
Command-line tool that interacts with Palo Alto firewalls and Panorama.

Primarily, this tool is used for importing and exporting objects using CSV files. Current abilities include:

* Import or export address, service objects and groups
  * You can also add existing objects to existing groups.
  * You can remove address objects from an existing group.
* Import or export an entire security policy.

More features will continue to be added in the future.

For a detailed explanation of commands, and how they are used, click on any one of the
command names below.

`panco` [`help`](https://github.com/scottdware/panco#usage), [`example`][example-doc], [`objects`][objects-doc], [`policy`][policy-doc], [`version`][version-doc]

## Installation

Installation of extremely easy, and just involves downloading the correct binary for your OS. You can download them by clicking on the [release](https://github.com/scottdware/panco/releases) tab above.

Current support OS's:

* Windows
* Mac OS

Just download the `panco-<OS>.zip` file, extract the binary and place it somewhere in your PATH.

## Usage

```
Usage:
  panco [command]

Available Commands:
  example     Create example CSV files for import reference
  help        Help about any command
  objects     Import and export address and service objects
  policy      Import and export a security policy
  version     Prints the version number of panco

Flags:
  -h, --help   help for panco

Use "panco [command] --help" for more information about a command.
```

### panco example

```
Usage:
  panco example [flags]

Flags:
  -h, --help   help for example
```

This command will create a sample CSV file for use with the import command. The 
files will be placed in the location where you are running the command from, and are named as such:
	
`panco-example-import.csv`

## Structuring Your CSV File

The CSV file for object creation (import) should be organized with the following columns:

`Name,Type,Value,Description,Tags,Device Group/Vsys`

* The `Description` and `Tags` fields are optional, however you **MUST** still include them even if they are blank in your file!
* If any line begins with a hashtag `#`, it WILL be ignored!

> **_NOTE_**: Here are a few things to keep in mind when creating objects:
> * For the name of the object, it cannot be longer than 32 characters, and must only include letters, numbers, spaces, hyphens, and underscores.
> * If you are tagging an object upon creation, please make sure that the tags exist prior to creating the objects.
> * When creating service groups, you DO NOT need to specify a description, as they do not have that capability.
> * When ran against a local firewall, the default value for `Vsys` is "vsys1" if you do not specify one. When ran against Panorama, the default value for `Device Group` is "shared."

> **_WARNING_**: If an existing address or service object has the same name as one you are creating, it's value will be overwritten with what you specify.

### Creating Address Objects

Column | Description
:--- | :---
`Name` | Name of the object you wish to create.
`Type` | **ip**, **range**, or **fqdn**
`Value` | Must contain the IP address, FQDN, or IP range of the object.
`Description` | (Optional) A description of the object.
`Tags` | (Optional) Name of a pre-existing tag on the device to apply.
`Device Group/Vsys` | Name of the Device Group or Vsys (defaults are: `shared` for Panorama, `vsys1` for a firewall).

### Creating Address Groups - or Add to Existing

Column | Description
:--- | :---
`Name` | Name of the address group you wish to create.
`Type` | `static` or `dynamic`
`Value` | ** See below explanation
`Description` | (Optional) A description of the object.
`Tags` | (Optional) Name of a pre-existing tag or tags on the device to apply. Separate multiple using a comma or semicolon.
`Device Group/Vsys` | Name of the Device Group or Vsys (defaults are: `shared` for Panorama, `vsys1` for a firewall).

For a `static` address group, `Value` must contain a comma, or semicolon separated list of members to add to the group, enclosed in quotes `""`, e.g.:

`"ip-host1, ip-net1; fqdn-example.com"`

For a `dynamic` address group, `Value` must contain the criteria (tags) to match on. This **_MUST_** be enclosed in quotes `""`, and
each criteria (tag) must be surrounded by single-quotes `'`, e.g.:

`"'Servers' or 'Web-Servers' and 'DMZ'"`

### Removing Objects From Address Groups

Column | Description
:--- | :---
`Name` | Name of the address group you wish to remove object(s) from.
`Type` | `remove-address`
`Value` | Must contain a comma, or semicolon separated list of members to remove from group, enclosed in quotes `""`.
`Description` | Not used - leave blank.
`Tags` | Not used - leave blank.
`Device Group/Vsys` | Name of the Device Group or Vsys (defaults are: `shared` for Panorama, `vsys1` for a firewall).

### Creating Service Objects

Column | Description
:--- | :---
`Name` | Name of the object you wish to create.
`Type` | `tcp` or `udp`
`Value` | ** See below
`Description` | (Optional) A description of the object.
`Tags` | (Optional) Name of a pre-existing tag or tags on the device to apply. Separate multiple using a comma or semicolon.
`Device Group/Vsys` | Name of the device-group, or **shared** if creating a shared object.

`Value` must contain a single port number (443), range (1023-3000), or comma-separated list of ports, enclosed in quotes, e.g.: `"80, 443, 8080"`.

### Creating Service Groups - or Add to Existing

Column | Description
:--- | :---
`Name` | Name of the object you wish to create.
`Type` | `service`
`Value` | ** See below
`Description` | Not used - leave blank (not available on service groups).
`Tags` | (Optional) Name of a pre-existing tag or tags on the device to apply. Separate multiple using a comma or semicolon.
`Device Group/Vsys` | Name of the device-group, or **shared** if creating a shared object.

** `Value` must contain a comma-separated list of service objects to add to the group, enclosed in quotes `""`, e.g.: `"tcp_8080, udp_666, tcp_range"`.

## panco objects [flags]

```
Usage:
  panco objects [flags]

Flags:
  -a, --action string        Action to perform; import or export
  -d, --device string        Device to connect to
  -g, --devicegroup string   Device Group name when exporting from Panorama (default "shared")
  -f, --file string          Name of the CSV file to import/export to
  -h, --help                 help for objects
  -u, --user string          User to connect to the device as
  -v, --vsys string          Vsys name when exporting from a firewall (default "vsys1")
```

This command allows you to import and export address and service objects.

Please run `panco example` for sample CSV file to use as a reference when importing.

## panco policy [flags]

```
Usage:
  panco policy [flags]

Flags:
  -a, --action string        Action to perform; import or export
  -d, --device string        Firewall or Panorama device to connect to
  -g, --devicegroup string   Device Group name; only needed when ran against Panorama
  -f, --file string          Name of the CSV file to import/export to
  -h, --help                 help for policy
  -l, --location string      Rule location; pre or post when ran against Panorama (default "post")
  -u, --user string          User to connect to the device as
  -v, --vsys string          Vsys name when ran against a firewall (default "vsys1")
```

This command will allow you to import and export an entire security policy. If
you are running this against a Panorama device, it can be really helpful if you want to clone
an entire policy, as you can export it from one device-group, modify it if needed, then import
the poilcy into a different device-group (or firewall).

## panco version

```
Usage:
  panco version [flags]

Flags:
  -h, --help   help for version
```

Version information for panco.

[example-doc]: https://github.com/scottdware/panco#panco-example
[policy-doc]: https://github.com/scottdware/panco#panco-policy-flags
[objects-doc]: https://github.com/scottdware/panco#panco-objects-flags
[version-doc]: https://github.com/scottdware/panco#panco-version