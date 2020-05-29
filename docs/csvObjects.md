# CSV Structure - Objects

This guide will help show you the way to structure your CSV file(s) for use when working with the various
objects types - address, service and tag.

When creating (importing) address, service and tag objects, the CSV file needs to have the following columns:

```
Name,Type,Value,Description,Tags,Device Group/Vsys
```

* The `Description` and `Tags` fields are optional, however you **_MUST_** still include them even if they are blank in your file!
* If any line begins with a hashtag `#`, it WILL be ignored!

**_TIP_**: A good example is to first [export](https://panco.dev/objectsExport.html) address, service or tag objects from the device.
That way, you get a good idea of how the CSV file is laid out.

> **_NOTE_**: Here are a few things to keep in mind when creating objects:
> * For the name of the object, it cannot be longer than 32 characters, and must only include letters, numbers, spaces, hyphens, and underscores.
> * If you are tagging an object upon creation, please make sure that the tags exist prior to creating the objects.
> * When creating service groups, you DO NOT need to specify a description, as they do not have that capability.
> * When ran against a local firewall, the default value for `Vsys` is "vsys1" if you do not specify one. When ran against Panorama, the default value for `Device Group` is "shared."

> **_WARNING_**: If an existing address or service object has the same name as one you are creating, it's value will be overwritten with what you specify.
> Please refer to the [`duplicates`](https://panco.dev/objectsDuplicates.html) command to help alleviate this.

Once your CSV file is organized with any of the following options, you can execute the changes using the following command:

```
panco objects import --file <name-of-CSV-file>
```

## Creating or Modifying Address Objects

Column | Description
:--- | :---
Name | Name of the address object you wish to create.
Type | **ip**, **ip-netmask**, **range** or **fqdn**
Value | Must contain the IP address, FQDN or IP range of the object.
Description | (Optional) A description of the object.
Tags | (Optional) Name of a pre-existing tag on the device to apply to the object.
Device Group/Vsys | Name of the Device Group or Vsys (defaults are: `shared` for Panorama, `vsys1` for a firewall).

## Creating Address Groups - or Add to Existing

Column | Description
:--- | :---
Name | Name of the address group you wish to create or add to.
Type | `static` or `dynamic`
Value | ** See below explanation
Description | (Optional) A description of the object.
Tags | (Optional) Name of a pre-existing tag or tags on the device to apply. Separate multiple using a comma or semicolon.
Device Group/Vsys | Name of the Device Group or Vsys (defaults are: `shared` for Panorama, `vsys1` for a firewall).

For a `static` address group, `Value` must contain a comma, or semicolon separated list of members to add to the group, enclosed in quotes `""`, e.g.:

```
"ip-host1, ip-net1; fqdn-example.com"
```

For a `dynamic` address group, `Value` must contain the criteria (tags) to match on. This **_MUST_** be enclosed in quotes `""`, and
each criteria (tag) must be surrounded by single-quotes `'`, e.g.:

```
"'Servers' or 'Web-Servers' and 'DMZ'"
```

## Removing Objects From Address Groups

Column | Description
:--- | :---
Name | Name of the address group you wish to remove object(s) from.
Type | `remove-address`
Value | Must contain a comma, or semicolon separated list of members to remove from group, enclosed in quotes `""`.
Description | Not used - leave blank.
Tags | Not used - leave blank.
Device Group/Vsys | Name of the Device Group or Vsys (defaults are: `shared` for Panorama, `vsys1` for a firewall).

## Creating or Modifying Service Objects

Column | Description
:--- | :---
Name | Name of the service object you wish to create.
Type | `tcp` or `udp`
Value | ** See below
Description | (Optional) A description of the object.
Tags | (Optional) Name of a pre-existing tag or tags on the device to apply. Separate multiple using a comma or semicolon.
Device Group/Vsys | Name of the device-group, or **shared** if creating a shared object.

`Value` must contain a single port number (443), range (1023-3000), or comma separated list of ports, enclosed in quotes, e.g.:

```
"80, 443, 8080"
```

## Creating Service Groups - or Add to Existing

Column | Description
:--- | :---
Name | Name of the service group you wish to create or add to.
Type | `service`
Value | ** See below
Description | Not used - leave blank (not available on service groups).
Tags | (Optional) Name of a pre-existing tag or tags on the device to apply. Separate multiple using a comma or semicolon.
Device Group/Vsys | Name of the device-group, or **shared** if creating a shared object.

** `Value` must contain a comma or semicolon separated list of service objects to add to the group, enclosed in quotes `""`, e.g.:

```
"tcp_8080, udp_666; tcp_range"
```

## Removing Objects From Service Groups

Column | Description
:--- | :---
Name | Name of the service group you wish to remove object(s) from.
Type | `remove-service`
Value | Must contain a comma, or semicolon separated list of members to remove from group, enclosed in quotes `""`.
Description | Not used - leave blank.
Tags | Not used - leave blank.
Device Group/Vsys | Name of the Device Group or Vsys (defaults are: `shared` for Panorama, `vsys1` for a firewall).

## Creating Tags

Column | Description
:--- | :---
Name | Name of the tag you wish to create.
Type | `tag`
Value | ** See below
Description | (Optional) A description of the tag.
Tags | Not used - leave blank
Device Group/Vsys | Name of the device-group, or **shared** if creating a shared object.

** `Value` is the color that you want the tag to represent. Below are the following colors available for use:

```
None, Red, Green, Blue, Yellow, Copper, Orange, Purple, Gray, Light Green, Cyan, Light Gray,
Blue Gray, Lime, Black, Gold, Brown, Olive, Maroon, Red-Orange, Yellow-Orange, Forest Green,
Turquoise Blue, Azure Blue, Cerulean Blue, Midnight Blue, Medium Blue, Cobalt Blue, Blue Violet,
Medium Violet, Medium Rose, Lavender, Orchid, Thistle, Peach, Salmon, Magenta, Red Violet,
Mahogany, Burnt Sienna, Chestnut
```

## Renaming Objects

Column | Description
:--- | :---
Name | Name of the object you wish to rename.
Type | One of: `rename-address`, `rename-addressgroup`, `rename-service`, `rename-servicegroup`
Value | New name of the object you wish to rename to.
Description | Not used for rename - leave blank.
Tags | Not used for rename - leave blank.
Device Group/Vsys | Name of the device-group/vsys, or **shared** if renaming a shared object.

## Tagging Objects

When applying tags to objects, use the same CSV format as you would when creating or modifying, but
in the `tags` column, enter in the name of a pre-existing tag you wish to apply.

If you wish to apply multiple tags, use a comma or semicolon separated list of pre-existing tags,
enclosed in quotes `""`, e.g.:

```
"Windows-Server, Internet-Access"
```