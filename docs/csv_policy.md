# CSV Structure - Policies

This guide will help show you the way to structure your CSV file(s) for use when working with the various
policy actions - importing or modifying rules, exporting rules, moving rules, grouping rules by tags.

## Importing (modify) A Security Policy

When importing a CSV file to create security rules or modify them, the file **_MUST_** have the following fields in this order:

`Name,Type,Description,Tags,SourceZones,SourceAddresses,NegateSource,SourceUsers,HipProfiles,DestinationZones,DestinationAddresses,NegateDestination,Applications,Services,Categories,Action,LogSetting,LogStart,LogEnd,Disabled,Schedule,IcmpUnreachable,DisableServerResponseInspection,Group,Virus,Spyware,Vulnerability,UrlFiltering,FileBlocking,WildFireAnalysis,DataFiltering`

**_TIP_**: The easiest way to make the modifications all while adhering to this order, is to export the security policy first (using the `panco policy export` command),
and then modifying the output file. For example:

```
panco policy export --type security --file <file-to-output>
```

## Importing (modify) A NAT Policy

When importing a CSV file to create NAT rules or modify them, the file **_MUST_** have the following columns in this order:

`Name,Type,Description,Tags,SourceZones,DestinationZone,ToInterface,Service,SourceAddresses,DestinationAddresses,SatType,SatAddressType,SatTranslatedAddresses,SatInterface,SatIpAddress,SatFallbackType,SatFallbackTranslatedAddresses,SatFallbackInterface,SatFallbackIpType,SatFallbackIpAddress,SatStaticTranslatedAddress,SatStaticBiDirectional,DatType,DatAddress,DatPort,DatDynamicDistribution,Disabled`

**_TIP_**: The easiest way to make the modifications all while adhering to this order, is to export the NAT policy first (using the `panco policy export` command),
and then modifying the output file. For example:

```
panco policy export --type nat --file <file-to-output>
```

## Import (modify) A Policy-Based Forwarding Policy

When importing a CSV file to create policy-based forwarding rules or modify them, the file **_MUST_** have the following columns in this order:

`Name,Description,Tags,FromType,FromValues,SourceAddresses,SourceUsers,NegateSource,DestinationAddresses,NegateDestination,Applications,Services,Schedule,Disabled,Action,ForwardVsys,ForwardEgressInterface,ForwardNextHopType,ForwardNextHopValue,ForwardMonitorProfile,ForwardMonitorIpAddress,ForwardMonitorDisableIfUnreachable,EnableEnforceSymmetricReturn,SymmetricReturnAddresses,ActiveActiveDeviceBinding,NegateTarget,Uuid`

**_TIP_**: The easiest way to make the modifications all while adhering to this order, is to export the NAT policy first (using the `panco policy export` command),
and then modifying the output file. For example:

```
panco policy export --type pbf --file <file-to-output>
```

## Moving Rules

When you specify the option to move multiple rules (`--movemultiple` flag in conjunction with `--file`) using a CSV file, here is the format that the CSV file must adhere to:

`RuleName,RuleDestination,TargetRule,<blank>,<blank>,Device group/Vsys`

## Importing and Modifying Rules

When you import (create) rules, or want to modify existing values of a rule, you **_DO NOT_** need to have every column that is listed above filled out with a value.  You still **_NEED_** them to be defined/listed, but they can be empty.

Any field that you want to add or modify you need to have a value there, but at the very least, you **_MUST_** have values in the following fields:

**Security Policy:**

`Name`, `Type`, `Action`

**NAT Policy:**

`Name`, `Type`, `ToInterface`

## Group Rules By Tags

You can group multiple rules by tags, which allow you to "View the Rulebase as Groups" as shown in Panorama and on the firewall Policy tab. To do so, you need to structure your CSV file with the following two columns:

`Rule Name,Tag`

If you want to do this on an existing rulebase, the easiest way is to first export the policy that you want, then, remove all of the other columns outside of the `Name` column and then add in what Tags you want applied to each rule to group them by. Once you have your file all set, run the following command:

```
panco policy group --file <name-of-CSV-file>
```

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
