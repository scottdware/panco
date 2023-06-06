# panco

Command-line tool that interacts with Palo Alto firewalls and Panorama using CSV files

パン粉 - Pronounced like the breadcrum!

Abilities include:

* Exporting objects from the device - address, service, tag
* Creating address, service and tag objects
* Renaming address, service and tag objects
* Adding or removing objects from address and service groups
* Creating, adding or removing URL entries from a custom URL category
* Tag multiple objects
* Deleting objects
* Exporting a security, NAT, Decryption or Policy-Based Forwarding (PBF) policy
* Creating security, NAT, Decryption or Policy-Based Forwarding (PBF) rules
* Editing security, NAT, Decryption or Policy-Based Forwarding (PBF) rules - e.g. adding a Log Profile to all rules
* Generate CLI set commands from a CSV file used for object actions

## Installation

Installation of extremely easy, and just involves downloading the correct binary for your OS. You can download them by clicking on the [release](https://github.com/scottdware/panco/releases) tab above.

Current support OS's:

* Windows
* Mac OS
* Linux

Just download the `panco-<OS>.zip` file, extract the binary and place it somewhere in your PATH.

## Getting Started/Documentation

* Visit the [panco Documentation](https://panco.dev) site.

There you will find in-depth documentation and examples on how to structure the CSV file(s), as well as using the different commands.

## Usage

```
Command-line tool that interacts with Palo Alto firewalls and Panorama using CSV files

See https://panco.dev for complete documentation

Usage:
  panco [command]

Available Commands:
  help        Help about any command
  objects     Commands to work with address, service and tag objects
  policy      Commands to work with security, NAT and Policy-Based Forwarding policies
  userid      Commands to interact with User-ID functions
  version     Version information for panco

Flags:
  -h, --help   help for panco

Use "panco [command] --help" for more information about a command.
```
