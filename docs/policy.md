# `policy` Command

This page will help to show you how the CSV file is formatted when used to import and update rules in a policy.

* [Usage](https://github.com/scottdware/panco/wiki/Policy#usage)
* [CSV Structure](https://github.com/scottdware/panco/wiki/Policy#csv-structure)
* [Examples](https://github.com/scottdware/panco/wiki/Policy#examples)

## Usage

```
This command will allow you to import and export an entire security policy, along
with moving rules within the policy. When importing, this allows you to create new rules,
or modify existing values in rules.

When moving rules, if you are only doing one at a time, you do not need to specify a CSV file
or the '--movemultiple' flag. However, if you are wanting to move multiple rules around, then
you will want to use a CSV file, and it must include the '--movemultiple' flag.

See https://github.com/scottdware/panco/Wiki for more information

Usage:
  panco policy [flags]

Flags:
  -a, --action string        Action to perform - import, export, groupbytag, or move (only for security policy)
  -d, --device string        Firewall or Panorama device to connect to
  -g, --devicegroup string   Device Group name; only needed when ran against Panorama
  -f, --file string          Name of the CSV file to import/export to
  -h, --help                 help for policy
  -l, --location string      Rule location; pre or post when ran against Panorama (default "post")
  -m, --movemultiple         Specifies you wish to move multiple security rules; use only with --file (default true)
  -x, --nat                  Run the given action on the NAT policy
  -p, --pass string          Password for the user account specified
  -w, --ruledest string      Where to move the rule - after, before, top, or bottom
  -n, --rulename string      Name of the security rule you wish to move
  -t, --targetrule string    Name of the rule 'ruledest' is referencing
  -u, --user string          User to connect to the device as
  -v, --vsys string          Vsys name when ran against a firewall (default "vsys1")
```

## CSV Structure

### Importing a Security Policy

When importing a CSV to create security rules or modify them, the file **_MUST_** have the following fields in this order:

`Name,Type,Description,Tags,SourceZones,SourceAddresses,NegateSource,SourceUsers,HipProfiles,DestinationZones,DestinationAddresses,NegateDestination,Applications,Services,Categories,Action,LogSetting,LogStart,LogEnd,Disabled,Schedule,IcmpUnreachable,DisableServerResponseInspection,Group,Virus,Spyware,Vulnerability,UrlFiltering,FileBlocking,WildFireAnalysis,DataFiltering`

The easiest way to make the modifications all while adhering to this order, is to export the security policy first (using the `--action export`), and then modifying that file.

Here is an example CSV file of a security policy that has been exported:

[PA-VM_Policy.csv](https://github.com/scottdware/panco-examples/blob/master/PA-VM_Policy.csv)

And here is a screenshot of this policy:

[Screenshot: Security Policy](https://github.com/scottdware/panco-examples/blob/master/sec-policy.png)

### Importing a NAT Policy

Similar to importing a security policy, when using a CSV file to create or modify NAT rules, the CSV file **_MUST_** have the following columns in this order:

`Name,Description,Type,SourceZones,DestinationZone,ToInterface,Service,SourceAddresses,DestinationAddresses,SatType,SatAddressType,SatTranslatedAddresses,SatInterface,SatIpAddress,SatFallbackType,SatFallbackTranslatedAddresses,SatFallbackInterface,SatFallbackIpType,SatFallbackIpAddress,SatStaticTranslatedAddress,SatStaticBiDirectional,DatType,DatAddress,DatPort,DatDynamicDistribution,Disabled,Tags`

The easiest way to make the modifications all while adhering to this order, is to export the NAT policy first (using the `--action export --xlate` flags), and then modifying that file.

Here is an example CSV file of a security policy that has been exported:

[PA-VM_NAT.csv](https://github.com/scottdware/panco-examples/blob/master/PA-VM_NAT.csv)

And here is a screenshot of this policy:

[Screenshot: NAT Policy](https://github.com/scottdware/panco-examples/blob/master/nat-policy.png)

### Moving Rules

When you specify the option to move multiple rules (`--movemultiple` flag in conjunction with `--file`) using a CSV file, here is the format that the CSV file must adhere to:

`RuleName,RuleDestination,TargetRule,<blank>,<blank>,Device group/Vsys`

Here is an example CSV file which has 2 rules we want to move:

[PA-VM_MoveRules.csv](https://github.com/scottdware/panco-examples/blob/master/PA-VM_MoveRules.csv)

## Importing and Modifying Rules

When you import (create) rules, or want to modify existing values of a rule, you **_DO NOT_** need to have every column that is listed above filled out with a value.  You still **_NEED_** them to be defined/listed, but they can be empty.

Any field that you want to add or modify you need to have a value there, but at the very least, you **_MUST_** have values in the following fields:

**Security Policy:**

`Name`, `Type`, `Action`

**NAT Policy:**

`Name`, `Type`, `ToInterface`

## Group Rules By Tag

You can group multiple rules by tags, which allow you to "View the Rulebase as Groups" as shown in Panorama and on the firewall Policy tab. To do so, you need to structure your CSV file with the following two columns:

`Rule name,Tag`

If you want to do this on an existing rulebase, probably the easiest way is to first export the policy that you want, then, just remove all of the other columns outside of the `Name` column and then add in what Tags you want applied to each rule to group them by.

## Examples

### Add Profiles

Based on the policy above, only three of the rules have security profiles (URL filtering) configured. Let's say we want to add more (AV, Vulnerability, Wildfire, etc.), as a security profile group to a couple of rules. Here is a CSV file that we will use to accomplise this:

[PA-VM_AddProfiles.csv](https://github.com/scottdware/panco-examples/blob/master/PA-VM_AddProfiles.csv)

As you can see in the file, we only have the `Name`, `Type`, `Action` and profile/group fields defined with the values we want to add. Once `panco` has imported this file, our policy should now reflect the security profiles that we defined:

[Screenshot: Policy with profiles](https://github.com/scottdware/panco-examples/blob/master/sec-policy-group.png)

---

### Add Rules

Now let's add a couple of rules to our policy, and then add tags to a couple of existing ones. Here is the CSV file we will be using for this task:

[PA-VM_AddSecRules.csv](https://github.com/scottdware/panco-examples/blob/master/PA-VM_AddSecRules.csv)

As you can see in this file, we are doing the following:

* Add a rule that will deny traffic to the `gaming` URL category from desktop clients, using the `Gaming_Apps` application filter we have defined.
* Add a rule that will deny access to Reddit.
* Add the `Network-Services` tag to two existing rules: `SSH_Outbound` and `RADIUS_Desktop_Clients`

Once we import this using `panco`, our policy will now reflect all of the changes we have done:

[Screenshot: Policy with new rules](https://github.com/scottdware/panco-examples/blob/master/sec-policy-newrules.png)

---

### Move Rules

Based on the above example, we can see that the newly added rules aren't really in the location that we want. Let's move `Block_Gaming`, `Block_Reddit` before the `Browsing_Desktop_Clients-Auth` rule so we know we will match the traffic correctly. Here's the CSV file we will use to accomplish this:

[PA-VM_MoveRules.csv](https://github.com/scottdware/panco-examples/blob/master/PA-VM_MoveRules.csv)

We can import this using `panco` by running the following command:

`panco policy --action move --file PA-VM_MoveRules.csv --movemultiple --device pa-vm --user admin`

Once we import this, our rule base will now look like the following:

[Screenshot: Policy with all changes](https://github.com/scottdware/panco-examples/blob/master/sec-policy-postmove.png)

#### Moving a Single Rule

If you wish to only move one rule at a time, you can do that by using the following options:

`panco policy --action move --rulename <existing rule to move> --ruledest <after|before|top|bottom> --targetrule <rule that references --ruledest> --device pa-vm --user admin`

So say you wish to move an `SSH` rule above a `Ping_Traceroute` rule, your command would look like:

`panco policy --action move --rulename SSH --ruledest before --targetrule Ping_Traceroute --device pa-vm --user admin`
