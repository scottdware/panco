# panco

Command-line tool that interacts with Palo Alto firewalls and Panorama.

Abilities include:

* Import (create, update) and export the following:
  * Address objects (IP, Range, FQDN), address groups (static and dynamic), service objects (TCP, UDP, port-ranges, etc.), service groups.

* Import (create, update) and export a security policy:
  * Add new rules to a policy.
  * Modify existing rules to update values (e.g. update lots of rules with a Log Profile, or Security Profile group).

* Import (create, update) and export a NAT policy.

* Move rules in a policy
  * Move a single rule, or multiple rules (using a CSV file) anywhere within a security policy.

* Configure a device using [IronSkillet](https://github.com/PaloAltoNetworks/iron-skillet) [`loadable_configs`](https://github.com/PaloAltoNetworks/iron-skillet/tree/panos_v8.0/loadable_configs) or from a different configuration file.
  * Can use a local file, or pull one from a remote HTTP location.

* Export a device's configuration to a file (XML).

* Group multiple rules (security, NAT) by a tag.

For a detailed explanation of commands, and how they are used, please visit the [Wiki](https://github.com/scottdware/panco/wiki) page or click on any one of the command names below (takes you to their respective Wiki page).

`panco` [`help`](https://github.com/scottdware/panco#usage), [`objects`][objects-doc], [`policy`][policy-doc], [`config`][config-doc]

## Installation

Installation of extremely easy, and just involves downloading the correct binary for your OS. You can download them by clicking on the [release](https://github.com/scottdware/panco/releases) tab above.

Current support OS's:

* Windows
* Mac OS

Just download the `panco-<OS>.zip` file, extract the binary and place it somewhere in your PATH.

### Build Option

You can also build the binaries yourself by cloning this repo, and running `go build -o <name of binary> panco\main.go`

## Getting Started

* Visit the [Wiki](https://github.com/scottdware/panco/wiki)

There you will find in-depth documentation and examples on how to structure the CSV file(s) when working with objects
 and policies.

## Usage

>**Note on passwords**: Due to some changes in Go code for certain functions this library uses, there were issues when inputing your
password on Windows, where it wouldn't respond unless you hit `CTRL-C`, which would exit you out of the program.
>The fix for this, is the new `--pass (-p)` flag attached to every command. I understand that this allows your password to be viewed
in plain text on the CLI, but I am working on finding other solutions for this.

```
Command-line tool that interacts with Palo Alto firewalls and Panorama.

See https://github.com/scottdware/panco/Wiki for more information

Usage:
  panco [command]

Available Commands:
  config      Configure a device using IronSkillet or a local or remote (HTTP) file; export a device configuration
  example     Create example CSV files for import reference
  help        Help about any command
  objects     Import and export address and service objects
  policy      Import/export a security policy, move rules
  version     Version information for panco

Flags:
  -h, --help   help for panco

Use "panco [command] --help" for more information about a command.
```

[objects-doc]: https://github.com/scottdware/panco/wiki/Objects
[policy-doc]: https://github.com/scottdware/panco/wiki/Policy
[config-doc]: https://github.com/scottdware/panco/wiki/Config