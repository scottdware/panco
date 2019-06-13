# panco
Command-line tool that interacts with Palo Alto firewalls and Panorama.

Some of the features include:

* Import or export address, service objects and groups
  * Add existing objects to existing groups.
  * Remove address objects from an existing address group.
* Import or export an entire security policy.
  * Create/add new rules to a policy.
  * Modify existing rules.
* Provision a device using IronSkillet or from a different configuration file.

More features will continue to be added in the future.

For a detailed explanation of commands, and how they are used, please visit the [Wiki](https://github.com/scottdware/panco/wiki) page or click on any one of the command names below (takes you to their respective Wiki page).

`panco` [`help`](https://github.com/scottdware/panco#usage), [`objects`][objects-doc], [`policy`][policy-doc], [`provision`][provision-doc]

## Installation

Installation of extremely easy, and just involves downloading the correct binary for your OS. You can download them by clicking on the [release](https://github.com/scottdware/panco/releases) tab above.

Current support OS's:

* Windows
* Mac OS

Just download the `panco-<OS>.zip` file, extract the binary and place it somewhere in your PATH.

## Getting Started

* Visit the [Wiki](https://github.com/scottdware/panco/wiki)

There you will find in-depth documentation and examples on how to structure the CSV file(s) when working with objects
 and policies.

## Usage

```
Command-line tool that interacts with Palo Alto firewalls and Panorama.

See https://github.com/scottdware/panco/Wiki for more information

Usage:
  panco [command]

Available Commands:
  example     Create example CSV files for import reference
  help        Help about any command
  objects     Import and export address and service objects
  policy      Import and export a security policy
  provision   Provision a device using IronSkillet or a local or remote (HTTP) file
  version     Prints the version number of panco

Flags:
  -h, --help   help for panco

Use "panco [command] --help" for more information about a command.
```

[objects-doc]: https://github.com/scottdware/panco/wiki/Objects
[policy-doc]: https://github.com/scottdware/panco/wiki/Policy
[provision-doc]: https://github.com/scottdware/panco/wiki/Provisioning