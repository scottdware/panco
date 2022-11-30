# About panco

* `panco` is written using the [Go](https://golang.org) programming language
* The underlying library behind it Palo Alto's [pango library](https://github.com/PaloAltoNetworks/pango), as well as a few custom API calls.

## How It Works

When `panco` creates or modifies an object, it uses the "set" API action against
the device. When modifying existing rules, the `panco policy modify` command does use the "edit" API action. There are very important differences between
 the two as described below.

> _Set and edit actions differ in two important ways:_
> * _Set actions add, update, or merge configuration nodes, while **edit actions replace configuration nodes**._
> * _Set actions are **non-destructive and are only additive**, while **edit actions can be destructive**._

The last bullet point is important! Please use caution when running the `modify` command against policy. More information on this can be found on the [modify](https://panco.dev/policy_modify.html) command page. 

For more infomation on these actions, please refer to the following guide from Palo Alto:

[Actions for Modifying a Configuration](https://docs.paloaltonetworks.com/pan-os/9-0/pan-os-panorama-api/pan-os-xml-api-request-types/pan-os-xml-api-request-types-and-actions/configuration-actions/actions-for-modifying-a-configuration.html)

## Contact Me

If you run into issues, or need assistance, please submit an issue on the [main Github repository](https://github.com/scottdware/panco).

## References

[PAN-OS® and Panorama™API Usage Guide](https://docs.paloaltonetworks.com/pan-os/9-0/pan-os-panorama-api.html)

[Palo Alto pango library](https://github.com/PaloAltoNetworks/pango)

[Go Programming Language](https://golang.org)