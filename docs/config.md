# Config Command

```
This command will allow you to configure your device with best practice templates that have
been setup by Palo Alto Networks, in their IronSkillet Github repository. You can also choose to load a
config file locally, or from a remote source using HTTP. Only the pre-built "loadable_configs" are an
option at this time.

You can also export the configuration from a device when the source is set to "export."

When using IronSkillet as the source (--source ironskillet), you will need to specifcy a couple of options:

'--os' must be either "panos" or "panorama"
'--config' can be one of: "aws", "azure", "gcp", "dhcp" or "static"

If you do not specify the --load flag, the the configuration will only be transferred to the device, where you
will have to load it manually. If you specify the --load flag, then the configuration will be loaded automatically,
but will still need to be manually committed.

See https://github.com/scottdware/panco/Wiki for more information

Usage:
  panco config [flags]

Flags:
  -d, --device string   Device to connect to
  -f, --file string     Name of the XML config file - only used when source is local, remote or export
  -h, --help            help for config
  -l, --load            Load the file into the running-config - USE WITH CAUTION
  -o, --os string       Device OS - only used when source is ironskillet - <panos|panorama>
  -p, --pass string     Password for the user account specified
  -s, --source string   Source of config - <ironskillet|local|remote|export>
  -t, --type string     IronSkillet config to use - <aws|azure|gcp|dhcp|static> (default "static")
  -u, --user string     User to connect to the device as
```

## Overview

Using the `config` command allows you to use default, best practice configuration templates to load into a device. These templates were created by Palo Alto Networks and are a part of their [IronSkillet](https://github.com/PaloAltoNetworks/iron-skillet) repository. You can also use a local file on your machine, or fetch a remote one using HTTP.

When using [IronSkillet](https://github.com/PaloAltoNetworks/iron-skillet) configs to provision a device, currently only the pre-built [`loadable_configs`](https://github.com/PaloAltoNetworks/iron-skillet/tree/panos_v8.0/loadable_configs) are an option at this time. Defining custom values will have to be done manually (e.g. can be changed in the GUI once the config is loaded), but are coming in a future release of `panco`.

## Examples

### IronSkillet Configs

In this example, we'll configure our firewall using one of the [IronSkillet](https://github.com/PaloAltoNetworks/iron-skillet) pre-built configs. Let's choose the `dhcp` option (management interface uses DHCP):

```
panco config --source ironskillet --os panos --config dhcp --device pa-vm --user admin
```

Behind the scenes, this is what the command will do:

* Go to the [`sample-mgmt-dhcp`](https://github.com/PaloAltoNetworks/iron-skillet/tree/panos_v8.0/loadable_configs/sample-mgmt-dhcp/panos) folder in the [IronSkillet](https://github.com/PaloAltoNetworks/iron-skillet) repository (`iron-skillet/loadable_configs/sample-mgmt-dhcp/panos`) and download the `iron_skillet_panos_full.xml` file.
* The file is then transferred to the device we specified, ready to be loaded into the running-config.

If we take a look at the screenshot, you can see the config file we just downloaded has been saved to the firewall, and is ready to be used:

<img src="https://github.com/scottdware/panco-examples/blob/master/provision_loaded_config.png" alt="IronSkillet loaded config"/>
<!-- ![Screenshot: IronSkillet loaded config](https://github.com/scottdware/panco-examples/blob/master/provision_loaded_config.png) -->

If we were to include the `--load` option, then the configuration file will also be loaded into the running-config, ready to be committed.

### Local and Remote Files

Choosing a local or a remote file to provision a device with works very similar as the above example. For these options, the parameters change as you do not need to specify your OS type, as well as a couple others:

```
panco config --source local --file /path/to/config.xml --device pa-vm --user admin
```

Remote files work the same way, except you would specify the URL to where the file is located:

```
panco config --source remote --file https://sdubs.org/path/to/config.xml --device pa-vm --user admin
```

Again, you can choose to add the `--load` flag if you wish to have the configuration automatically loaded into the running-config.

### Export Device Configuration

```
panco config --source export --file firewall-config.xml --device pa-vm --user admin
```