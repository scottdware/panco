# Moving Rules

```
Usage:
  panco policy move [flags]

Flags:
  -d, --device string   Firewall or Panorama device to connect to
  -f, --file string     Name of the CSV file
  -h, --help            help for move
  -p, --pass string     Password for the user account specified
  -u, --user string     User to connect to the device as
```

## Overview

Using this command, you can move rules in security, NAT or Policy-Based Forwarding (PBF) rulebases anywhere within their
respective policy. Within the same CSV file, you can specify the different types of policies, or specify a single type (e.g. security).

The destnination can be any of the four types:

* After
* Before
* Top
* Bottom

Please use the below link as a guide on how to structure your CSV file when moving rules:

[CSV Structure - Policies (Moving Rules)](https://panco.dev/csvPolicy.html#moving-rules)