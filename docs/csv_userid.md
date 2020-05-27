## CSV Structure - Userid Functions

This guide will help show you the way to structure your CSV file(s) for use when working with the different
user-id based functions.

The current actions you can perform are as follows:

* Tagging/untagging an IP addresses for use in dynamic address groups (no commit required).
* Manually registering/mapping an IP address to a user-id, as well as un-registering/removing the mapping from them.

## Tagging/Untagging IP Addresses

Here is the format that your CSV file must adhere to when tagging or untagging an IP address:

Column | Description
:--- | :---
IP Address | This can **_ONLY_** be a single IP address. Subnets are not supported at this time.
Action | Action to perform - `tag` or `untag`
Tag | Name of the tag(s) you wish to apply.

* When specifying the IP address, it does **_NOT_** need to pre-exist as an object on the device.
* If you wish to apply multiple tags to an IP address, they must be enclosed in quotes and separated by a comma - e.g.: `"Malicious, Block-IP"`
  * If you are using Excel to modify the CSV file, then it will automatically enclose the comma-separated list in quotes.

## Logging In/Logging Out Users

When logging in (registering) or logging out (unregistering) users, here is the format that your CSV file must adhere to:

Column | Description
:--- | :---
User | Name of the user to login/logout.
Action | Action to perform - `login` or `logout`
IP Address | IP address to assign to the user.

* You can also include a domain for the user if necessary - e.g.: `MYDOMAIN\user`