[panco.dev](https://panco.dev) \| [Documentation Home](https://panco.dev/docs.html) \| [Policy Command](https://panco.dev/policy.html)

# Exporting A Policy

```
Usage:
  panco policy export [flags]

Flags:
  -d, --device string        Device to connect to
  -g, --devicegroup string   Device Group name when exporting from Panorama (default "shared")
  -f, --file string          Name of the CSV file you'd like to export to (default "PaloAltoPolicy")
  -h, --help                 help for export
  -l, --location string      Location of the rulebase - <pre|post> (default "pre")
  -r, --rules string         [OPTIONAL] Only export these specific rules - specify text file
  -t, --type string          Type of policy to export - <security|nat|pbf|decrypt|all>
  -u, --user string          User to connect to the device as
  -v, --vsys string          Vsys name when exporting from a firewall (default "vsys1")
```

## Overview

This command allows you to export a seurity, NAT or policy-based forwarding (PBF) policy to a CSV format. You can choose
to specify them separately, or export them all at once. Please note, that given the size of your rulebase, it could take
a couple of minutes to export all of the rules.

## Specifying Only Certain Rules to Export

You can optionally specify certain rules to export, instead of them all. This can be done by placing the rule names
you'd like to export in a text file, one on each line, and then referencing that text file with the `-r` flag.

**Example Text File Contents**

Say you have the following file called `rules.txt` with the below contents:

```
Allow-DNS
Block Malicious Sites
VPN traffic
Allow-Social-Media
```

With the above text file, you can run the below command against a firewall and it will only export the four (4) rules that are listed in the `rules.txt`
file:

`panco policy export -d 10.1.1.1 -u admin -t security -r rules.txt -f SpecificRules_from_Policy.csv`

## Exported Rules CSV Format

Each policies CSV file will be formatted differently. Below are the formats for each of them:

**Security**

```
Name,Type,Description,Tags,SourceZones,SourceAddresses,NegateSource,SourceUsers,HipProfiles,
DestinationZones,DestinationAddresses,NegateDestination,Applications,Services,Categories,Action,
LogSetting,LogStart,LogEnd,Disabled,Schedule,IcmpUnreachable,DisableServerResponseInspection,
Group,Virus,Spyware,Vulnerability,UrlFiltering,FileBlocking,WildFireAnalysis,DataFiltering
```

**NAT**

```
Name,Type,Description,Tags,SourceZones,DestinationZone,ToInterface,Service,SourceAddresses,
DestinationAddresses,SatType,SatAddressType,SatTranslatedAddresses,SatInterface,SatIpAddress,
SatFallbackType,SatFallbackTranslatedAddresses,SatFallbackInterface,SatFallbackIpType,
SatFallbackIpAddress,SatStaticTranslatedAddress,SatStaticBiDirectional,DatType,DatAddress,
DatPort,DatDynamicDistribution,Disabled
```

**Policy-Based Forwarding (PBF)**

```
Name,Description,Tags,FromType,FromValues,SourceAddresses,SourceUsers,NegateSource,
DestinationAddresses,NegateDestination,Applications,Services,Schedule,Disabled,Action,
ForwardVsys,ForwardEgressInterface,ForwardNextHopType,ForwardNextHopValue,ForwardMonitorProfile,
ForwardMonitorIpAddress,ForwardMonitorDisableIfUnreachable,EnableEnforceSymmetricReturn,
SymmetricReturnAddresses,ActiveActiveDeviceBinding,NegateTarget,Uuid
```

To export all policies, execute the following command:

```
panco policy export --file <name-of-output-file> --type all
```
