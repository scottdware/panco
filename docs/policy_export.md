# Exporting A Policy

```
Usage:
  panco policy export [flags]

Flags:
  -d, --device string        Device to connect to
  -g, --devicegroup string   Device Group name when exporting from Panorama (default "shared")
  -f, --file string          Name of the CSV file you'd like to export to (default "PaloAltoPolicy")
  -h, --help                 help for export
  -p, --pass string          Password for the user account specified
  -t, --type string          Type of policy to export - <security|nat|pbf|all>
  -u, --user string          User to connect to the device as
  -v, --vsys string          Vsys name when exporting from a firewall (default "vsys1")
```

## Overview

This command allows you to export a seurity, NAT or policy-based forwarding (PBF) policy to a CSV format. You can choose
to specify them separately, or export them all at once. Please note, that given the size of your rulebase, it could take
a couple of minutes to export all of the rules.

Each policies CSV file will be formatted differently. Below are the formats for each of them:

**Security**

`Name,Type,Description,Tags,SourceZones,SourceAddresses,NegateSource,SourceUsers,HipProfiles,DestinationZones,DestinationAddresses,NegateDestination,Applications,Services,Categories,Action,LogSetting,LogStart,LogEnd,Disabled,Schedule,IcmpUnreachable,DisableServerResponseInspection,Group,Virus,Spyware,Vulnerability,UrlFiltering,FileBlocking,WildFireAnalysis,DataFiltering`

**NAT**

`Name,Type,Description,Tags,SourceZones,DestinationZone,ToInterface,Service,SourceAddresses,DestinationAddresses,SatType,SatAddressType,SatTranslatedAddresses,SatInterface,SatIpAddress,SatFallbackType,SatFallbackTranslatedAddresses,SatFallbackInterface,SatFallbackIpType,SatFallbackIpAddress,SatStaticTranslatedAddress,SatStaticBiDirectional,DatType,DatAddress,DatPort,DatDynamicDistribution,Disabled`

**Policy-Based Forwarding (PBF)**

`Name,Description,Tags,FromType,FromValues,SourceAddresses,SourceUsers,NegateSource,DestinationAddresses,NegateDestination,Applications,Services,Schedule,Disabled,Action,ForwardVsys,ForwardEgressInterface,ForwardNextHopType,ForwardNextHopValue,ForwardMonitorProfile,ForwardMonitorIpAddress,ForwardMonitorDisableIfUnreachable,EnableEnforceSymmetricReturn,SymmetricReturnAddresses,ActiveActiveDeviceBinding,NegateTarget,Uuid`

To export all policies, execute the following command:

```
panco policy export --file <name-of-output-file> --type all
```