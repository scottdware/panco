# CSV Structure - Policies

This guide will help show you the way to structure your CSV file(s) for use when working with the various
policy actions - importing or modifying rules, exporting rules, moving rules, grouping rules by tags.

## Importing A Security Policy

When importing a CSV file to create security rules or modify them, the file **_MUST_** have the following fields in this order:

```
Name,Type,Description,Tags,SourceZones,SourceAddresses,NegateSource,SourceUsers,HipProfiles,
DestinationZones,DestinationAddresses,NegateDestination,Applications,Services,Categories,Action,
LogSetting,LogStart,LogEnd,Disabled,Schedule,IcmpUnreachable,DisableServerResponseInspection,
Group,Virus,Spyware,Vulnerability,UrlFiltering,FileBlocking,WildFireAnalysis,DataFiltering
```

> **_TIP_**: The easiest way to make the modifications all while adhering to this order, is to export the security policy first (using the `panco policy export` command),
> and then modifying the output file. For example:

```
panco policy export -d firewall -u admin -g "Device-Group" --type security --file <file-to-output>
```

> *NOTE:* When ran against Panorama, be sure to use the `--location` flag to specify which rulebase to import/create the rules on. By default
> this is the "pre" rulebase.

When you create rules, or want to modify existing values of a rule, you **_DO NOT_** need to have every column that is listed above filled out with a value. You still **_NEED_** them to be defined/listed, but they can be empty.

Any field that you want to add or modify you need to have a value there, but at the very least, you **_MUST_** have values in the following fields:

```
Name,Type,Action
```

## Importing A NAT Policy

When importing a CSV file to create NAT rules or modify them, the file **_MUST_** have the following columns in this order:

```
Name,Type,Description,Tags,SourceZones,DestinationZone,ToInterface,Service,SourceAddresses,
DestinationAddresses,SatType,SatAddressType,SatTranslatedAddresses,SatInterface,SatIpAddress,
SatFallbackType,SatFallbackTranslatedAddresses,SatFallbackInterface,SatFallbackIpType,
SatFallbackIpAddress,SatStaticTranslatedAddress,SatStaticBiDirectional,DatType,DatAddress,
DatPort,DatDynamicDistribution,Disabled
```

> **_TIP_**: The easiest way to make the modifications all while adhering to this order, is to export the NAT policy first (using the `panco policy export` command),
> and then modifying the output file. For example:

```
panco policy export -d firewall -u admin -g "Device-Group" --type nat --file <file-to-output>
```

>*NOTE:* When ran against Panorama, be sure to use the `--location` flag to specify which rulebase to import/create the rules on. By default
> this is the "post" rulebase.

When you create rules, or want to modify existing values of a rule, you **_DO NOT_** need to have every column that is listed above filled out with a value. You still **_NEED_** them to be defined/listed, but they can be empty.

Any field that you want to add or modify you need to have a value there, but at the very least, you **_MUST_** have values in the following fields:

```
Name,Type,ToInterface
```

## Import A Policy-Based Forwarding Policy

When importing a CSV file to create policy-based forwarding rules or modify them, the file **_MUST_** have the following columns in this order:

```
Name,Description,Tags,FromType,FromValues,SourceAddresses,SourceUsers,NegateSource,
DestinationAddresses,NegateDestination,Applications,Services,Schedule,Disabled,Action,
ForwardVsys,ForwardEgressInterface,ForwardNextHopType,ForwardNextHopValue,ForwardMonitorProfile,
ForwardMonitorIpAddress,ForwardMonitorDisableIfUnreachable,EnableEnforceSymmetricReturn,
SymmetricReturnAddresses,ActiveActiveDeviceBinding,NegateTarget,Uuid
```

**_TIP_**: The easiest way to make the modifications all while adhering to this order, is to export the NAT policy first (using the `panco policy export` command),
and then modifying the output file. For example:

```
panco policy export -d firewall -u admin -g "Device-Group" --type pbf --file <file-to-output>
```

> *NOTE:* When ran against Panorama, be sure to use the `--location` flag to specify which rulebase to import/create the rules on. By default
> this is the "post" rulebase.

## Modifying A Security, NAT or Policy-Based Forwarding Policy/Rules -- IMPORTANT

When you modify (edit) rules using the `panco policy modify` command, there are a few things to be aware of.  The `modify` command uses the Palo Alto API `edit` action, instead of the `set` action that is used when using the `import` command. You can read more about the differences of the `edit` and `set` [here](edit-set)

Set and edit actions differ in two important ways:
* Set actions add, update, or merge configuration nodes, while **_edit actions replace configuration nodes_**.
* Set actions are non-destructive and are only additive, while **_edit actions can be destructive_**.

> **_IMPORTANT_**: Please read and understand the above actions when using the `panco policy modify` command vs `panco policy import`.

Using the `modify` command will ultimately be the best way to make changes to rules, such as adding/removing address objects, applications, services, etc.. Similar to the `import` command, the best way to preserve the current state of the rule(s) you are modifying, is to first export the policy/rules you need to modify using the below command:

```
panco policy export -d firewall -u admin -g "Device-Group" --type security --file <file-to-output>
```

Once you have exported the rules, then you can add/remove values from the different fields as needed, before then running the `panco policy modify` command on the CSV file you just edited.


## Moving Rules

When using the `panco policy move` command, here is the format that the CSV file must adhere to:

Column | Description
:--- | :---
Rule Type | Type of rule - `security`, `nat` or `pbf`
Location | ** Only used when ran against Panorama (`pre` or `post`); leave blank otherwise.
Rule Name | Name of the rule you wish to move.
Destination | Where to move the rule - `after`, `before`, `top` or `bottom`
Target Rule | Target rule where `Destination` is referencing.
Device Group/Vsys | Name of the Device Group or Vsys (defaults are: `shared` for Panorama, `vsys1` for a firewall).

Once you have specified what rules you need to move, you can execute it with the following command:

```
panco policy move --file <name-of-CSV-file>
```

## Group Rules By Tags

You can group multiple rules by tags, which allow you to "View the Rulebase as Groups" as shown in Panorama and on the firewall Policy tab. To do so, you need to structure your CSV file with the following two columns:

Column | Description
:--- | :---
Rule Name | Name of the rule you wish to order-by tag.
Tag | Name of the tag you wish to group rules by - MUST be pre-existing on the device.

If you want to do this on an existing rulebase, the easiest way is to first export the policy that you want, then, remove all of the other columns outside of the `Name` and `Tag` columns and then add in what tags you want applied to each rule to group them by. Once you have your file all set, run the following command:

```
panco policy group --file <name-of-CSV-file> --type <security|nat>
```

[edit-set](https://docs.paloaltonetworks.com/pan-os/10-2/pan-os-panorama-api/pan-os-xml-api-request-types/pan-os-xml-api-request-types-and-actions/configuration-actions/actions-for-modifying-a-configuration#id44705ad2-4f22-4b6c-bb94-caea78a6d510)
