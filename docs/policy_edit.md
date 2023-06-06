[Home](https://panco.dev) > [Policy Command](https://panco.dev/policy.html)

# Editing A Policy

```
Edit a security, NAT, Decryption or Policy-Based Forwarding policy

Usage:
  panco policy edit [flags]

Flags:
  -d, --device string        Device to connect to
  -g, --devicegroup string   Device Group name when importing to Panorama (default "shared")
  -f, --file string          Name of the CSV file to export to
  -h, --help                 help for modify
  -l, --location string      Location of the rulebase - <pre|post> (default "pre")
  -t, --type string          Type of policy to import - <security|nat|decrypt|pbf>
  -u, --user string          User to connect to the device as
  -v, --vsys string          Vsys name when importing to a firewall (default "vsys1")
```

## Overview

Using the `edit` command allows you to edit existing rules, by adding or removing entries from
each of the rule fields. You can modify/edit the following types of policies at this time:

* Security
* NAT
* Decryption
* Policy-Based Forwarding (PBF)

Please use the below link as a guide on how to structure your CSV file when modifying rules:

[CSV Structure - Policies](https://panco.dev/csv_policy.html)

## Important Tips

When you edit rules using the `panco policy edit` command, there are a few things to be aware of.  The `edit` command uses the Palo Alto API `edit` action, instead of the `set` action that is used when using the `import` command. You can read more about the differences of the `edit` and `set` on [Palo Alto's API request types documentation](https://docs.paloaltonetworks.com/pan-os/10-2/pan-os-panorama-api/pan-os-xml-api-request-types/pan-os-xml-api-request-types-and-actions/configuration-actions/actions-for-modifying-a-configuration#id44705ad2-4f22-4b6c-bb94-caea78a6d510) page.

Set and edit actions differ in two important ways:
* Set actions add, update, or merge configuration nodes, while **_edit actions replace configuration nodes_**.
* Set actions are non-destructive and are only additive, while **_edit actions can be destructive_**.

> **_IMPORTANT_**: Please read and understand the above actions when using the `panco policy edit` command vs `panco policy import`.

Using the `edit` command will ultimately be the best way to make changes to rules, such as adding/removing address objects, applications, services, etc.. Similar to the `import` command, the best way to preserve the current state of the rule(s) you are modifying, is to first export the policy/rules you need to modify using the below command:

```
panco policy export -d firewall -u admin -g "Device-Group" --type security --file <file-to-output>
```

Once you have exported the rules, then you can add/remove values from the different fields as needed, before then running the `panco policy edit` command on the CSV file you just edited.

