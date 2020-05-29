# About panco

* `panco` is written using the [Go](https://golang.org) programming language
* The underlying library behind it Palo Alto's [pango library](https://github.com/PaloAltoNetworks/pango)

## How It Works

When `panco` creates or modifies an object or rule, it uses the "set" API action against
the device. I chose this method over the "edit" API action for the following reasons (taken from the below link)

> _Set and edit actions differ in two important ways:_
> * _Set actions add, update, or merge configuration nodes, while edit actions replace configuration nodes._
> * _Set actions are non-destructive and are only additive, while edit actions can be destructive._

The last bullet point is important - I like to err on the side of caution :)

Now, 99.9% of the time, you will not see anything out of the ordinary when using `panco` to add or modify
objects or rules. Examples where you might see strange behavior are:

* If you want to modify a rule, and the rule name is incorrect (typo) - it might create a completely
new rule with only the tag or modification you were trying to do.
* If you are trying to tag rules and specify the wrong rule location on Panorama (e.g. "pre" or "post" rules), you
will end up creating rules in said location with just the tag.

During my testing, those are the only cases where I've seen "oddities" happen. But in all cases, **_NOTHING_** is
deleted or removed fromt the configuration (unless you choose to remove something).

For more infomation on these actions, please refer to the following guide from Palo Alto:

[Actions for Modifying a Configuration](https://docs.paloaltonetworks.com/pan-os/9-0/pan-os-panorama-api/pan-os-xml-api-request-types/pan-os-xml-api-request-types-and-actions/configuration-actions/actions-for-modifying-a-configuration.html)

If you run into issues, or need assistance, please submit an issue on the [main Github repository](https://github.com/scottdware/panco),
or drop me a line on Twitter [@scottdware](https://twitter.com/scottdware).

## References

[PAN-OS® and Panorama™API Usage Guide](https://docs.paloaltonetworks.com/pan-os/9-0/pan-os-panorama-api.html)

[Palo Alto pango library](https://github.com/PaloAltoNetworks/pango)

[Go Programming Language](https://golang.org)