# panco

パン粉 - Pronounced like the breadcrum!

## What Is It

`panco` is a command-line tool which helps to automate bulk tasks when working with [Palo Alto Networks](https://paloaltonetworks.com) firewalls
and Panorama, which normally would be cumbersome having to do them one-by-one. Features of this tool include:

* Exporting objects from the device - address, service, tag
* Creating address, service and tag objects
* Renaming address, service and tag objects
* Adding or removing objects from address and service groups
* Finding duplicate address and service objects
* Tag multiple objects
* Deleting objects
* Exporting a security, NAT or Policy-Based Forwarding (PBF) policy
* Creating security, NAT or Policy-Based Forwarding (PBF) rules
* Modifying security, NAT or Policy-Based Forwarding (PBF) rules - e.g. adding a Log Profile to all rules
* Group security or NAT rules by tags
* Move multiple security, NAT or Policy-Based Forwarding (PBF) rules at a time
* Get the hit count data on security, NAT or Policy-Based Forwarding (PBF) rules
* Tag/untag IP addresses for use in dynamic address groups
* Manually login/logout a user and map them to an IP address
* Generate CLI set commands from a CSV file used for importing objects, policies

> **_NOTE_**: Your account must have API access to the devices in order to use this tool

Check out the [About panco](https://panco.dev/about.html) page for more info.

## How Do I Get It

You can install `panco` by clicking on the below link and downloading the binary for your specific OS
(Windows, Mac OS and Linux are supported). Links are available on the main [Github repo](https://github.com/scottdware/panco) page as well.

* [Download panco!](https://github.com/scottdware/panco/releases)

Once you download the binary, place it in your `PATH` environment variable, or run it from it's current location.

## Documentation

Access the full package documentation below!

* [panco Documentation](https://panco.dev/docs.html)

## About

Check out the [About panco](https://panco.dev/about.html) page for more info.

## Help & Support

If you run into issues, or need assistance, please submit an issue on the [main Github repository](https://github.com/scottdware/panco),
or drop me a line on Twitter [@scottdware](https://twitter.com/scottdware).