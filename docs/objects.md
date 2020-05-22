# Objects Command

This page will help to show you how to structure your CSV file for importing and updating address and service objects.

* [Usage](https://github.com/scottdware/panco/wiki/Objects#usage)
* [CSV Structure](https://github.com/scottdware/panco/wiki/Objects#csv-structure)
* [Example Files](https://github.com/scottdware/panco/wiki/Objects#example-files)

## Usage

```
This command allows you to import and export address and service objects. You
can also add objects to groups, as well as remove addreses objects from address groups.

Please run "panco example" for sample CSV file to use as a reference when importing.

See https://github.com/scottdware/panco/Wiki for more information

Usage:
  panco objects [flags]

Flags:
  -a, --action string        Action to perform; import or export
  -d, --device string        Device to connect to
  -g, --devicegroup string   Device Group name when exporting from Panorama (default "shared")
  -f, --file string          Name of the CSV file to import/export to
  -h, --help                 help for objects
  -p, --pass string          Password for the user account specified
  -u, --user string          User to connect to the device as
  -v, --vsys string          Vsys name when exporting from a firewall (default "vsys1")
```

## CSV Structure

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

`Value` must contain a single port number (443), range (1023-3000), or comma separated list of ports, enclosed in quotes, e.g.: `"80, 443, 8080"`.

### Creating Service Groups - or Add to Existing

Column | Description
:--- | :---
`Name` | Name of the object you wish to create.
`Type` | `service`
`Value` | ** See below
`Description` | Not used - leave blank (not available on service groups).
`Tags` | (Optional) Name of a pre-existing tag or tags on the device to apply. Separate multiple using a comma or semicolon.
`Device Group/Vsys` | Name of the device-group, or **shared** if creating a shared object.

** `Value` must contain a comma or semicolon separated list of service objects to add to the group, enclosed in quotes `""`, e.g.: `"tcp_8080, udp_666; tcp_range"`.

### Creating Tags

Column | Description
:--- | :---
`Name` | Name of the tag you wish to create.
`Type` | `tag`
`Value` | ** See below
`Description` | (Optional) A description of the tag.
`Tags` | Not used - leave blank
`Device Group/Vsys` | Name of the device-group, or **shared** if creating a shared object.

** `Value` is the color that you want the tag to represent. Below are the following colors available for use:

`None`, `Red`, `Green`, `Blue`, `Yellow`, `Copper`, `Orange`, `Purple`, `Gray`, `Light Green`, `Cyan`, `Light Gray`, `Blue Gray`, `Lime`, `Black`, `Gold`, `Brown`, `Olive`, `Maroon`, `Red-Orange`, `Yellow-Orange`, `Forest Green`, `Turquoise Blue`, `Azure Blue`, `Cerulean Blue`, `Midnight Blue`, `Medium Blue`, `Cobalt Blue`, `Blue Violet`, `Medium Violet`, `Medium Rose`, `Lavender`, `Orchid`, `Thistle`, `Peach`, `Salmon`, `Magenta`, `Red Violet`, `Mahogany`, `Burnt Sienna`, `Chestnut`

### Renaming Objects

Column | Description
:--- | :---
`Name` | Name of the object you wish to rename.
`Type` | One of: `rename-address`, `rename-addressgroup`, `rename-service`, `rename-servicegroup`
`Value` | New name of the object you wish to rename to.
`Description` | Not used for rename - leave blank.
`Tags` | Not used for rename - leave blank.
`Device Group/Vsys` | Name of the device-group/vsys, or **shared** if renaming a shared object.

## Example Files

* [objects.csv](https://github.com/scottdware/panco-examples/blob/master/objects.csv) - Sample file that will create address and service objects.
* [update.csv](https://github.com/scottdware/panco-examples/blob/master/update.csv) - Sample file that will remove an address object from a group, and add a service object to one.