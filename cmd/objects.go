// Copyright © 2018 Scott Ware <scottdware@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/addr"
	"github.com/PaloAltoNetworks/pango/objs/addrgrp"
	"github.com/PaloAltoNetworks/pango/objs/srvc"
	"github.com/PaloAltoNetworks/pango/objs/srvcgrp"
	easycsv "github.com/scottdware/go-easycsv"
	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
)

// objectsCmd represents the objects command
var objectsCmd = &cobra.Command{
	Use:   "objects",
	Short: "Import and export address and service objects",
	Long: `This command allows you to import and export address and service objects. You
can also add objects to groups, as well as remove addreses objects from address groups.

Please run "panco example" for sample CSV file to use as a reference when importing.

See https://github.com/scottdware/panco/Wiki for more information`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		// pass := passwd()
		resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

		cl := pango.Client{
			Hostname: device,
			Username: user,
			Password: pass,
			Logging:  pango.LogQuiet,
		}

		con, err := pango.Connect(cl)
		if err != nil {
			log.Printf("Failed to connect: %s", err)
			os.Exit(1)
		}

		switch c := con.(type) {
		case *pango.Firewall:
			if action == "export" {
				if v == "" {
					v = "vsys1"
				}

				fh = strings.TrimSuffix(fh, ".csv")
				afh := fmt.Sprintf("%s_addr.csv", fh)
				agfh := fmt.Sprintf("%s_addrgrp.csv", fh)
				sfh := fmt.Sprintf("%s_srvc.csv", fh)
				sgfh := fmt.Sprintf("%s_srvcgrp.csv", fh)

				log.Printf("Exporting objects - this might take a few of minutes if you have a lot of objects")

				// Address objects
				ac, err := easycsv.NewCSV(afh)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				addrs, err := c.Objects.Address.GetList(v)
				if err != nil {
					log.Printf("Failed to get the list of address objects: %s", err)
				}

				ac.Write("#Name,Type,Value,Description,Tags,Device Group/Vsys\n")
				for _, aentry := range addrs {
					a, err := c.Objects.Address.Get(v, aentry)
					if err != nil {
						log.Printf("Failed to retrieve object data for '%s': %s", aentry, err)
					}

					ac.Write(fmt.Sprintf("%s,%s,\"%s\",\"%s\",\"%s\",%s\n", a.Name, a.Type, a.Value, a.Description, sliceToString(a.Tags), v))
				}

				ac.End()

				// Address groups
				agc, err := easycsv.NewCSV(agfh)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				addrgrps, err := c.Objects.AddressGroup.GetList(v)
				if err != nil {
					log.Printf("Failed to get the list of address groups: %s", err)
				}

				agc.Write("#Name,Type,Value,Description,Tags,Device Group/Vsys\n")
				for _, agentry := range addrgrps {
					var gtype, val string
					a, err := c.Objects.AddressGroup.Get(v, agentry)
					if err != nil {
						log.Printf("Failed to retrieve object data for '%s': %s", agentry, err)
					}

					if len(a.StaticAddresses) <= 0 && len(a.DynamicMatch) > 0 {
						gtype = "dynamic"
						val = a.DynamicMatch
					}

					if len(a.DynamicMatch) <= 0 && len(a.StaticAddresses) > 0 {
						gtype = "static"
						val = sliceToString(a.StaticAddresses)
					}

					agc.Write(fmt.Sprintf("%s,%s,\"%s\",\"%s\",\"%s\",%s\n", a.Name, gtype, val, a.Description, sliceToString(a.Tags), v))
				}

				agc.End()

				// Service objects
				sc, err := easycsv.NewCSV(sfh)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				srvcs, err := c.Objects.Services.GetList(v)
				if err != nil {
					log.Printf("Failed to get the list of service objects: %s", err)
				}

				sc.Write("#Name,Type,Value,Description,Tags,Device Group/Vsys\n")
				for _, sentry := range srvcs {
					s, err := c.Objects.Services.Get(v, sentry)
					if err != nil {
						log.Printf("Failed to retrieve object data for '%s': %s", sentry, err)
					}

					sc.Write(fmt.Sprintf("%s,%s,\"%s\",\"%s\",\"%s\",%s\n", s.Name, s.Protocol, s.DestinationPort, s.Description, sliceToString(s.Tags), v))
				}

				sc.End()

				// Service groups
				sgc, err := easycsv.NewCSV(sgfh)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				srvcgrps, err := c.Objects.ServiceGroup.GetList(v)
				if err != nil {
					log.Printf("Failed to get the list of service groups: %s", err)
				}

				sgc.Write("#Name,Type,Value,Description,Tags,Device Group/Vsys\n")
				for _, sgentry := range srvcgrps {
					sg, err := c.Objects.ServiceGroup.Get(v, sgentry)
					if err != nil {
						log.Printf("Failed to retrieve object data for '%s': %s", sgentry, err)
					}

					sgc.Write(fmt.Sprintf("%s,service,\"%s\",,\"%s\",%s\n", sg.Name, sliceToString(sg.Services), sliceToString(sg.Tags), v))
				}

				sgc.End()
			}

			if action == "import" {
				lines, err := easycsv.Open(fh)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				lc := len(lines)
				log.Printf("Importing %d lines - this might take a few of minutes if you have a lot of objects", lc)

				for i, line := range lines {
					var vsys string
					llen := len(line)

					if llen <= 6 && len(line[5]) <= 0 {
						vsys = "vsys1"
					}

					name := line[0]
					otype := line[1]
					value := line[2]
					desc := line[3]
					tg := line[4]
					vsys = line[5]

					switch otype {
					case "ip", "IP Netmask", "ip-netmask":
						e := addr.Entry{
							Name:        name,
							Value:       strings.TrimSpace(value),
							Type:        addr.IpNetmask,
							Description: desc,
							Tags:        stringToSlice(tg),
						}

						err = c.Objects.Address.Set(vsys, e)
						if err != nil {
							log.Printf("Line %d - failed to create %s: %s", i+1, name, err)
						}
					case "range", "IP Range", "ip-range":
						e := addr.Entry{
							Name:        name,
							Value:       strings.TrimSpace(value),
							Type:        addr.IpRange,
							Description: desc,
							Tags:        stringToSlice(tg),
						}

						err = c.Objects.Address.Set(vsys, e)
						if err != nil {
							log.Printf("Line %d - failed to create %s: %s", i+1, name, err)
						}
					case "fqdn", "FQDN", "Fqdn":
						e := addr.Entry{
							Name:        name,
							Value:       strings.TrimSpace(value),
							Type:        addr.Fqdn,
							Description: desc,
							Tags:        stringToSlice(tg),
						}

						err = c.Objects.Address.Set(vsys, e)
						if err != nil {
							log.Printf("Line %d - failed to create %s: %s", i+1, name, err)
						}
					case "tcp", "udp":
						e := srvc.Entry{
							Name:            name,
							Description:     desc,
							Protocol:        otype,
							DestinationPort: strings.TrimSpace(value),
							Tags:            stringToSlice(tg),
						}

						err = c.Objects.Services.Set(vsys, e)
						if err != nil {
							log.Printf("Line %d - failed to create %s: %s", i+1, name, err)
						}
					case "service":
						e := srvcgrp.Entry{
							Name:     name,
							Services: stringToSlice(value),
							Tags:     stringToSlice(tg),
						}

						err = c.Objects.ServiceGroup.Set(vsys, e)
						if err != nil {
							log.Printf("Line %d - failed to create/update %s: %s", i+1, name, err)
						}
					case "static":
						e := addrgrp.Entry{
							Name:            name,
							Description:     desc,
							StaticAddresses: stringToSlice(value),
							Tags:            stringToSlice(tg),
						}

						err = c.Objects.AddressGroup.Set(vsys, e)
						if err != nil {
							log.Printf("Line %d - failed to create/update %s: %s", i+1, name, err)
						}
					case "dynamic":
						e := addrgrp.Entry{
							Name:         name,
							Description:  desc,
							DynamicMatch: value,
							Tags:         stringToSlice(tg),
						}

						err = c.Objects.AddressGroup.Set(vsys, e)
						if err != nil {
							log.Printf("Line %d - failed to create %s: %s", i+1, name, err)
						}
					case "remove-address":
						if len(value) <= 0 {
							log.Printf("Line %d - you must specify a value to remove from group: %s", i+1, name)
						}

						remove := stringToSlice(value)
						cur, err := c.Objects.AddressGroup.Get(vsys, name)
						if err != nil {
							log.Printf("Line %d - could not retrieve object: %s", i+1, err)
						}

						newaddrs := cur.StaticAddresses

						for idx, ev := range cur.StaticAddresses {
							for _, rv := range remove {
								if ev == rv {
									newaddrs = append(newaddrs[:idx], newaddrs[idx+1:]...)
									break
								}
							}
						}

						e := addrgrp.Entry{
							Name:            name,
							StaticAddresses: newaddrs,
						}

						err = c.Objects.AddressGroup.Edit(vsys, e)
						if err != nil {
							log.Printf("Line %d - failed to update %s: %s", i+1, name, err)
						}
					case "rename-address":
						var xpath string

						xpath = fmt.Sprintf("/config/devices/entry[@name='localhost.localdomain']/vsys/entry[@name='%s']/address/entry[@name='%s']", vsys, name)

						_, err := resty.R().Post(fmt.Sprintf("type=config&action=rename&xpath=%s&newname=%s&key=%s", xpath, value, c.ApiKey))
						if err != nil {
							log.Printf("Line %d - failed to rename object %s: %s", i+1, name, err)
						}
					case "rename-addressgroup":
						var xpath string

						xpath = fmt.Sprintf("/config/devices/entry[@name='localhost.localdomain']/vsys/entry[@name='%s']/address-group/entry[@name='%s']", vsys, name)

						_, err := resty.R().Post(fmt.Sprintf("type=config&action=rename&xpath=%s&newname=%s&key=%s", xpath, value, c.ApiKey))
						if err != nil {
							log.Printf("Line %d - failed to rename object %s: %s", i+1, name, err)
						}
					case "rename-service":
						var xpath string

						xpath = fmt.Sprintf("/config/devices/entry[@name='localhost.localdomain']/vsys/entry[@name='%s']/service/entry[@name='%s']", vsys, name)

						_, err := resty.R().Post(fmt.Sprintf("type=config&action=rename&xpath=%s&newname=%s&key=%s", xpath, value, c.ApiKey))
						if err != nil {
							log.Printf("Line %d - failed to rename object %s: %s", i+1, name, err)
						}
					case "rename-servicegroup":
						var xpath string

						xpath = fmt.Sprintf("/config/devices/entry[@name='localhost.localdomain']/vsys/entry[@name='%s']/service-group/entry[@name='%s']", vsys, name)

						_, err := resty.R().Post(fmt.Sprintf("type=config&action=rename&xpath=%s&newname=%s&key=%s", xpath, value, c.ApiKey))
						if err != nil {
							log.Printf("Line %d - failed to rename object %s: %s", i+1, name, err)
						}
					}
				}
			}
		case *pango.Panorama:
			if action == "export" {
				if dg == "" {
					dg = "shared"
				}

				fh = strings.TrimSuffix(fh, ".csv")
				afh := fmt.Sprintf("%s_addr.csv", fh)
				agfh := fmt.Sprintf("%s_addrgrp.csv", fh)
				sfh := fmt.Sprintf("%s_srvc.csv", fh)
				sgfh := fmt.Sprintf("%s_srvcgrp.csv", fh)

				log.Printf("Exporting objects - this might take a few of minutes if you have a lot of objects")

				// Address objects
				ac, err := easycsv.NewCSV(afh)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				addrs, err := c.Objects.Address.GetList(dg)
				if err != nil {
					log.Printf("Failed to get the list of address objects: %s", err)
				}

				ac.Write("#Name,Type,Value,Description,Tags,Device Group/Vsys\n")
				for _, aentry := range addrs {
					a, err := c.Objects.Address.Get(dg, aentry)
					if err != nil {
						log.Printf("Failed to retrieve object data for '%s': %s", aentry, err)
					}

					ac.Write(fmt.Sprintf("%s,%s,\"%s\",\"%s\",\"%s\",%s\n", a.Name, a.Type, a.Value, a.Description, sliceToString(a.Tags), dg))
				}

				ac.End()

				// Address groups
				agc, err := easycsv.NewCSV(agfh)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				addrgrps, err := c.Objects.AddressGroup.GetList(dg)
				if err != nil {
					log.Printf("Failed to get the list of address groups: %s", err)
				}

				agc.Write("#Name,Type,Value,Description,Tags,Device Group/Vsys\n")
				for _, agentry := range addrgrps {
					var gtype, val string
					a, err := c.Objects.AddressGroup.Get(dg, agentry)
					if err != nil {
						log.Printf("Failed to retrieve object data for '%s': %s", agentry, err)
					}

					if len(a.StaticAddresses) <= 0 && len(a.DynamicMatch) > 0 {
						gtype = "dynamic"
						val = a.DynamicMatch
					}

					if len(a.DynamicMatch) <= 0 && len(a.StaticAddresses) > 0 {
						gtype = "static"
						val = sliceToString(a.StaticAddresses)
					}

					agc.Write(fmt.Sprintf("%s,%s,\"%s\",\"%s\",\"%s\",%s\n", a.Name, gtype, val, a.Description, sliceToString(a.Tags), dg))
				}

				agc.End()

				// Service objects
				sc, err := easycsv.NewCSV(sfh)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				srvcs, err := c.Objects.Services.GetList(dg)
				if err != nil {
					log.Printf("Failed to get the list of service objects: %s", err)
				}

				sc.Write("#Name,Type,Value,Description,Tags,Device Group/Vsys\n")
				for _, sentry := range srvcs {
					s, err := c.Objects.Services.Get(dg, sentry)
					if err != nil {
						log.Printf("Failed to retrieve object data for '%s': %s", sentry, err)
					}

					sc.Write(fmt.Sprintf("%s,%s,\"%s\",\"%s\",\"%s\",%s\n", s.Name, s.Protocol, s.DestinationPort, s.Description, sliceToString(s.Tags), dg))
				}

				sc.End()

				// Service groups
				sgc, err := easycsv.NewCSV(sgfh)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				srvcgrps, err := c.Objects.ServiceGroup.GetList(dg)
				if err != nil {
					log.Printf("Failed to get the list of service groups: %s", err)
				}

				sgc.Write("#Name,Type,Value,Description,Tags,Device Group/Vsys\n")
				for _, sgentry := range srvcgrps {
					sg, err := c.Objects.ServiceGroup.Get(dg, sgentry)
					if err != nil {
						log.Printf("Failed to retrieve object data for '%s': %s", sgentry, err)
					}

					sgc.Write(fmt.Sprintf("%s,service,\"%s\",,\"%s\",%s\n", sg.Name, sliceToString(sg.Services), sliceToString(sg.Tags), dg))
				}

				sgc.End()
			}

			if action == "import" {
				lines, err := easycsv.Open(fh)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				lc := len(lines)
				log.Printf("Importing %d lines - this might take a few of minutes if you have a lot of objects", lc)

				for i, line := range lines {
					var dgroup string
					llen := len(line)

					if llen <= 6 && len(line[5]) <= 0 {
						dgroup = "shared"
					}

					name := line[0]
					otype := line[1]
					value := line[2]
					desc := line[3]
					tg := line[4]
					dgroup = line[5]

					switch otype {
					case "ip", "IP Netmask", "ip-netmask":
						e := addr.Entry{
							Name:        name,
							Value:       strings.TrimSpace(value),
							Type:        addr.IpNetmask,
							Description: desc,
							Tags:        stringToSlice(tg),
						}

						err = c.Objects.Address.Set(dgroup, e)
						if err != nil {
							log.Printf("Line %d - failed to create %s: %s", i+1, name, err)
						}
					case "range", "IP Range", "ip-range":
						e := addr.Entry{
							Name:        name,
							Value:       strings.TrimSpace(value),
							Type:        addr.IpRange,
							Description: desc,
							Tags:        stringToSlice(tg),
						}

						err = c.Objects.Address.Set(dgroup, e)
						if err != nil {
							log.Printf("Line %d - failed to create %s: %s", i+1, name, err)
						}
					case "fqdn", "FQDN", "Fqdn":
						e := addr.Entry{
							Name:        name,
							Value:       strings.TrimSpace(value),
							Type:        addr.Fqdn,
							Description: desc,
							Tags:        stringToSlice(tg),
						}

						err = c.Objects.Address.Set(dgroup, e)
						if err != nil {
							log.Printf("Line %d - failed to create %s: %s", i+1, name, err)
						}
					case "tcp", "udp":
						e := srvc.Entry{
							Name:            name,
							Description:     desc,
							Protocol:        otype,
							DestinationPort: strings.TrimSpace(value),
							Tags:            stringToSlice(tg),
						}

						err = c.Objects.Services.Set(dgroup, e)
						if err != nil {
							log.Printf("Line %d - failed to create %s: %s", i+1, name, err)
						}
					case "service":
						e := srvcgrp.Entry{
							Name:     name,
							Services: stringToSlice(value),
							Tags:     stringToSlice(tg),
						}

						err = c.Objects.ServiceGroup.Set(dgroup, e)
						if err != nil {
							log.Printf("Line %d - failed to create/update %s: %s", i+1, name, err)
						}
					case "static":
						e := addrgrp.Entry{
							Name:            name,
							Description:     desc,
							StaticAddresses: stringToSlice(value),
							Tags:            stringToSlice(tg),
						}

						err = c.Objects.AddressGroup.Set(dgroup, e)
						if err != nil {
							log.Printf("Line %d - failed to create/update %s: %s", i+1, name, err)
						}
					case "dynamic":
						e := addrgrp.Entry{
							Name:         name,
							Description:  desc,
							DynamicMatch: value,
							Tags:         stringToSlice(tg),
						}

						err = c.Objects.AddressGroup.Set(dgroup, e)
						if err != nil {
							log.Printf("Line %d - failed to create %s: %s", i+1, name, err)
						}
					case "remove-address":
						if len(value) <= 0 {
							log.Printf("Line %d - you must specify a value to remove from group: %s", i+1, name)
						}

						remove := stringToSlice(value)
						cur, err := c.Objects.AddressGroup.Get(dgroup, name)
						if err != nil {
							log.Printf("Line %d - could not retrieve object: %s", i+1, err)
						}

						newaddrs := cur.StaticAddresses

						for idx, ev := range cur.StaticAddresses {
							for _, rv := range remove {
								if ev == rv {
									newaddrs = append(newaddrs[:idx], newaddrs[idx+1:]...)
									break
								}
							}
						}

						e := addrgrp.Entry{
							Name:            name,
							StaticAddresses: newaddrs,
						}

						err = c.Objects.AddressGroup.Edit(dgroup, e)
						if err != nil {
							log.Printf("Line %d - failed to update %s: %s", i+1, name, err)
						}
					case "rename-address":
						var xpath string

						if dgroup == "shared" {
							xpath = fmt.Sprintf("/config/shared/address/entry[@name='%s']", name)
						}

						if dgroup != "shared" {
							xpath = fmt.Sprintf("/config/devices/entry[@name='localhost.localdomain']/device-group/entry[@name='%s']/address/entry[@name='%s']", dgroup, name)
						}

						_, err := resty.R().Post(fmt.Sprintf("type=config&action=rename&xpath=%s&newname=%s&key=%s", xpath, value, c.ApiKey))
						if err != nil {
							log.Printf("Line %d - failed to rename object %s: %s", i+1, name, err)
						}
					case "rename-addressgroup":
						var xpath string

						if dgroup == "shared" {
							xpath = fmt.Sprintf("/config/shared/address-group/entry[@name='%s']", name)
						}

						if dgroup != "shared" {
							xpath = fmt.Sprintf("/config/devices/entry[@name='localhost.localdomain']/device-group/entry[@name='%s']/address-group/entry[@name='%s']", dgroup, name)
						}

						_, err := resty.R().Post(fmt.Sprintf("type=config&action=rename&xpath=%s&newname=%s&key=%s", xpath, value, c.ApiKey))
						if err != nil {
							log.Printf("Line %d - failed to rename object %s: %s", i+1, name, err)
						}
					case "rename-service":
						var xpath string

						if dgroup == "shared" {
							xpath = fmt.Sprintf("/config/shared/service/entry[@name='%s']", name)
						}

						if dgroup != "shared" {
							xpath = fmt.Sprintf("/config/devices/entry[@name='localhost.localdomain']/device-group/entry[@name='%s']/service/entry[@name='%s']", dgroup, name)
						}

						_, err := resty.R().Post(fmt.Sprintf("type=config&action=rename&xpath=%s&newname=%s&key=%s", xpath, value, c.ApiKey))
						if err != nil {
							log.Printf("Line %d - failed to rename object %s: %s", i+1, name, err)
						}
					case "rename-servicegroup":
						var xpath string

						if dgroup == "shared" {
							xpath = fmt.Sprintf("/config/shared/service-group/entry[@name='%s']", name)
						}

						if dgroup != "shared" {
							xpath = fmt.Sprintf("/config/devices/entry[@name='localhost.localdomain']/device-group/entry[@name='%s']/service-group/entry[@name='%s']", dgroup, name)
						}

						_, err := resty.R().Post(fmt.Sprintf("type=config&action=rename&xpath=%s&newname=%s&key=%s", xpath, value, c.ApiKey))
						if err != nil {
							log.Printf("Line %d - failed to rename object %s: %s", i+1, name, err)
						}
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(objectsCmd)

	objectsCmd.Flags().StringVarP(&action, "action", "a", "", "Action to perform; import or export")
	objectsCmd.Flags().StringVarP(&fh, "file", "f", "", "Name of the CSV file to import/export to")
	objectsCmd.Flags().StringVarP(&user, "user", "u", "", "User to connect to the device as")
	objectsCmd.Flags().StringVarP(&pass, "pass", "p", "", "Password for the user account specified")
	objectsCmd.Flags().StringVarP(&device, "device", "d", "", "Device to connect to")
	objectsCmd.Flags().StringVarP(&dg, "devicegroup", "g", "shared", "Device Group name when exporting from Panorama")
	objectsCmd.Flags().StringVarP(&v, "vsys", "v", "vsys1", "Vsys name when exporting from a firewall")
	objectsCmd.MarkFlagRequired("user")
	objectsCmd.MarkFlagRequired("pass")
	objectsCmd.MarkFlagRequired("device")
	objectsCmd.MarkFlagRequired("action")
	objectsCmd.MarkFlagRequired("file")
}
