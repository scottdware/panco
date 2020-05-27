# CSV Structure - Policies

This guide will help show you the way to structure your CSV file(s) for use when working with the various
policy actions - importing or modifying rules, exporting rules, moving rules, grouping rules by tags.

## Importing (modify) A Security Policy

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
panco policy export --type security --file <file-to-output>
```

When you create rules, or want to modify existing values of a rule, you **_DO NOT_** need to have every column that is listed above filled out with a value. You still **_NEED_** them to be defined/listed, but they can be empty.

Any field that you want to add or modify you need to have a value there, but at the very least, you **_MUST_** have values in the following fields:

```
Name,Type,Action
```

## Importing (modify) A NAT Policy

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
panco policy export --type nat --file <file-to-output>
```

When you create rules, or want to modify existing values of a rule, you **_DO NOT_** need to have every column that is listed above filled out with a value. You still **_NEED_** them to be defined/listed, but they can be empty.

Any field that you want to add or modify you need to have a value there, but at the very least, you **_MUST_** have values in the following fields:

```
Name,Type,ToInterface
```

## Import (modify) A Policy-Based Forwarding Policy

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
panco policy export --type pbf --file <file-to-output>
```

## Moving Rules

When using the `panco policy move` command, here is the format that the CSV file must adhere to:

Column | Description
:--- | :---
`Rule Type` | Type of rule - `security`, `nat` or `pbf`
`Location` | ** Only used when ran against Panorama (`pre` or `post`); leave blank otherwise.
`Rule Name` | Name of the rule you wish to move.
`Destination` | Where to move the rule - `after`, `before`, `top` or `bottom`
`Target Rule` | Target rule where `Destination` is referencing.
`Device Group/Vsys` | Name of the Device Group or Vsys (defaults are: `shared` for Panorama, `vsys1` for a firewall).

Once you have specified what rules you need to move, you can execute it with the following command:

```
panco policy move --file <name-of-CSV-file>
```

## Group Rules By Tags

You can group multiple rules by tags, which allow you to "View the Rulebase as Groups" as shown in Panorama and on the firewall Policy tab. To do so, you need to structure your CSV file with the following two columns:

Column | Description
:--- | :---
`Rule Name` | Name of the rule you wish to order-by tag.
`Tag` | Name of the tag you wish to group rules by - MUST be pre-existing on the device.

If you want to do this on an existing rulebase, the easiest way is to first export the policy that you want, then, remove all of the other columns outside of the `Name` column and then add in what Tags you want applied to each rule to group them by. Once you have your file all set, run the following command:

```
panco policy group --file <name-of-CSV-file> --type <security|nat>
```