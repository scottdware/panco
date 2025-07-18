/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/addr"
	"github.com/PaloAltoNetworks/pango/objs/addrgrp"
	"github.com/PaloAltoNetworks/pango/objs/custom/url"
	"github.com/PaloAltoNetworks/pango/objs/srvc"
	"github.com/PaloAltoNetworks/pango/objs/srvcgrp"
	"github.com/PaloAltoNetworks/pango/objs/tags"
	"github.com/Songmu/prompter"
	easycsv "github.com/scottdware/go-easycsv"
	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
)

// importCmd represents the import command
var objectsImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import (create, modify) address, service, custom URL and tag objects",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var delay time.Duration
		resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
		keyrexp := regexp.MustCompile(`key=([0-9A-Za-z\=]+).*`)
		passwd := prompter.Password(fmt.Sprintf("Password for %s", user))
		_ = passwd

		if p == "" {
			delay, _ = time.ParseDuration("100ms")
		} else {
			delay, _ = time.ParseDuration(fmt.Sprintf("%sms", p))
		}

		cl := pango.Client{
			Hostname: device,
			Username: user,
			Password: passwd,
			Logging:  pango.LogQuiet,
		}

		con, err := pango.Connect(cl)
		if err != nil {
			log.Printf("Failed to connect: %s", err)
			os.Exit(1)
		}

		switch c := con.(type) {
		case *pango.Firewall:
			lines, err := easycsv.Open(f)
			if err != nil {
				log.Printf("CSV file error - %s", err)
				os.Exit(1)
			}

			lc := len(lines)
			log.Printf("Running actions on %d lines - this might take a few of minutes if you have a lot of objects", lc)

			timeoutCount := 0
			timeoutData := []string{}

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
					if value == "delete" {
						err = c.Objects.Address.Delete(vsys, name)
						if err != nil {
							if strings.Contains(err.Error(), "Client.Timeout") {
								timeoutCount++
								timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
							} else {
								log.Printf("Line %d - failed to delete %s: %s", i+1, name, err)
							}
						}
					} else {
						if len(name) > 63 {
							log.Printf("Line %d - failed to create %s: name is over the max 63 characters", i+1, name)
						} else {
							e := addr.Entry{
								Name:        name,
								Value:       strings.TrimSpace(value),
								Type:        addr.IpNetmask,
								Description: desc,
								Tags:        stringToSlice(tg),
							}

							err = c.Objects.Address.Set(vsys, e)
							if err != nil {
								if strings.Contains(err.Error(), "Client.Timeout") {
									timeoutCount++
									timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
								} else {
									log.Printf("Line %d - failed to create %s: %s", i+1, name, err)
								}
							}
						}
					}
				case "range", "IP Range", "ip-range":
					if value == "delete" {
						err = c.Objects.Address.Delete(vsys, name)
						if err != nil {
							if strings.Contains(err.Error(), "Client.Timeout") {
								timeoutCount++
								timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
							} else {
								log.Printf("Line %d - failed to delete %s: %s", i+1, name, err)
							}
						}
					} else {
						if len(name) > 63 {
							log.Printf("Line %d - failed to create %s: name is over the max 63 characters", i+1, name)
						} else {
							e := addr.Entry{
								Name:        name,
								Value:       strings.TrimSpace(value),
								Type:        addr.IpRange,
								Description: desc,
								Tags:        stringToSlice(tg),
							}

							err = c.Objects.Address.Set(vsys, e)
							if err != nil {
								if strings.Contains(err.Error(), "Client.Timeout") {
									timeoutCount++
									timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
								} else {
									log.Printf("Line %d - failed to create %s: %s", i+1, name, err)
								}
							}
						}
					}
				case "fqdn", "FQDN", "Fqdn":
					if value == "delete" {
						err = c.Objects.Address.Delete(vsys, name)
						if err != nil {
							if strings.Contains(err.Error(), "Client.Timeout") {
								timeoutCount++
								timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
							} else {
								log.Printf("Line %d - failed to delete %s: %s", i+1, name, err)
							}
						}
					} else {
						if len(name) > 63 {
							log.Printf("Line %d - failed to create %s: name is over the max 63 characters", i+1, name)
						} else {
							e := addr.Entry{
								Name:        name,
								Value:       strings.TrimSpace(value),
								Type:        addr.Fqdn,
								Description: desc,
								Tags:        stringToSlice(tg),
							}

							err = c.Objects.Address.Set(vsys, e)
							if err != nil {
								if strings.Contains(err.Error(), "Client.Timeout") {
									timeoutCount++
									timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
								} else {
									log.Printf("Line %d - failed to create %s: %s", i+1, name, err)
								}
							}
						}
					}
				case "tcp", "udp":
					if value == "delete" {
						err = c.Objects.Services.Delete(vsys, name)
						if err != nil {
							if strings.Contains(err.Error(), "Client.Timeout") {
								timeoutCount++
								timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
							} else {
								log.Printf("Line %d - failed to delete %s: %s", i+1, name, err)
							}
						}
					} else {
						if len(name) > 63 {
							log.Printf("Line %d - failed to create %s: name is over the max 63 characters", i+1, name)
						} else {
							e := srvc.Entry{
								Name:            name,
								Description:     desc,
								Protocol:        otype,
								DestinationPort: strings.TrimSpace(value),
								Tags:            stringToSlice(tg),
							}

							err = c.Objects.Services.Set(vsys, e)
							if err != nil {
								if strings.Contains(err.Error(), "Client.Timeout") {
									timeoutCount++
									timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
								} else {
									log.Printf("Line %d - failed to create %s: %s", i+1, name, err)
								}
							}
						}
					}
				case "service":
					if value == "delete" {
						err = c.Objects.ServiceGroup.Delete(vsys, name)
						if err != nil {
							if strings.Contains(err.Error(), "Client.Timeout") {
								timeoutCount++
								timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
							} else {
								log.Printf("Line %d - failed to delete %s: %s", i+1, name, err)
							}
						}
					} else {
						if len(name) > 63 {
							log.Printf("Line %d - failed to create %s: name is over the max 63 characters", i+1, name)
						} else {
							e := srvcgrp.Entry{
								Name:     name,
								Services: stringToSlice(value),
								Tags:     stringToSlice(tg),
							}

							err = c.Objects.ServiceGroup.Set(vsys, e)
							if err != nil {
								if strings.Contains(err.Error(), "Client.Timeout") {
									timeoutCount++
									timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
								} else {
									log.Printf("Line %d - failed to create/update %s: %s", i+1, name, err)
								}
							}
						}
					}
				case "static":
					if value == "delete" {
						err = c.Objects.AddressGroup.Delete(vsys, name)
						if err != nil {
							if strings.Contains(err.Error(), "Client.Timeout") {
								timeoutCount++
								timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
							} else {
								log.Printf("Line %d - failed to delete %s: %s", i+1, name, err)
							}
						}
					} else {
						if len(name) > 63 {
							log.Printf("Line %d - failed to create %s: name is over the max 63 characters", i+1, name)
						} else {
							groupLen := len(stringToSlice(value))

							e := addrgrp.Entry{
								Name:            name,
								Description:     desc,
								StaticAddresses: stringToSlice(value),
								Tags:            stringToSlice(tg),
							}

							err = c.Objects.AddressGroup.Set(vsys, e)
							if err != nil {
								if strings.Contains(err.Error(), "Client.Timeout") {
									timeoutCount++
									timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
								} else {
									log.Printf("Line %d - failed to create/update %s: %s", i+1, name, err)
									if groupLen > 40 {
										log.Printf("Line %d - address group %s is over 40 members, try to add/create/breakup the group with a smaller number of members (20-30)", i+1, name)
									}
								}
							}
						}
					}
				case "dynamic":
					if value == "delete" {
						err = c.Objects.AddressGroup.Delete(vsys, name)
						if err != nil {
							if strings.Contains(err.Error(), "Client.Timeout") {
								timeoutCount++
								timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
							} else {
								log.Printf("Line %d - failed to delete %s: %s", i+1, name, err)
							}
						}
					} else {
						if len(name) > 63 {
							log.Printf("Line %d - failed to create %s: name is over the max 63 characters", i+1, name)
						} else {
							e := addrgrp.Entry{
								Name:         name,
								Description:  desc,
								DynamicMatch: value,
								Tags:         stringToSlice(tg),
							}

							err = c.Objects.AddressGroup.Set(vsys, e)
							if err != nil {
								if strings.Contains(err.Error(), "Client.Timeout") {
									timeoutCount++
									timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
								} else {
									log.Printf("Line %d - failed to create %s: %s", i+1, name, err)
								}
							}
						}
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
						if strings.Contains(err.Error(), "Client.Timeout") {
							timeoutCount++
							timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
						} else {
							log.Printf("Line %d - failed to update %s: %s", i+1, name, err)
						}
					}
				case "remove-service":
					if len(value) <= 0 {
						log.Printf("Line %d - you must specify a value to remove from group: %s", i+1, name)
					}

					remove := stringToSlice(value)
					cur, err := c.Objects.ServiceGroup.Get(vsys, name)
					if err != nil {
						log.Printf("Line %d - could not retrieve object: %s", i+1, err)
					}

					newsrvcs := cur.Services

					for idx, ev := range cur.Services {
						for _, rv := range remove {
							if ev == rv {
								newsrvcs = append(newsrvcs[:idx], newsrvcs[idx+1:]...)
								break
							}
						}
					}

					e := srvcgrp.Entry{
						Name:     name,
						Services: newsrvcs,
					}

					err = c.Objects.ServiceGroup.Edit(vsys, e)
					if err != nil {
						if strings.Contains(err.Error(), "Client.Timeout") {
							timeoutCount++
							timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
						} else {
							log.Printf("Line %d - failed to update %s: %s", i+1, name, err)
						}
					}
				case "rename-address":
					var xpath string

					xpath = fmt.Sprintf("/config/devices/entry[@name='localhost.localdomain']/vsys/entry[@name='%s']/address/entry[@name='%s']", vsys, name)

					_, err := resty.R().Get(fmt.Sprintf("https://%s/api/?type=config&action=rename&xpath=%s&newname=%s&key=%s", device, xpath, value, c.ApiKey))
					if err != nil {
						if strings.Contains(err.Error(), "Client.Timeout") {
							timeoutCount++
							timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
						} else {
							formatkey := keyrexp.ReplaceAllString(err.Error(), "key=********")
							log.Printf("Line %d - failed to rename object %s: %s", i+1, name, formatkey)
						}
					}
				case "rename-addressgroup":
					var xpath string

					xpath = fmt.Sprintf("/config/devices/entry[@name='localhost.localdomain']/vsys/entry[@name='%s']/address-group/entry[@name='%s']", vsys, name)

					_, err := resty.R().Get(fmt.Sprintf("https://%s/api/?type=config&action=rename&xpath=%s&newname=%s&key=%s", device, xpath, value, c.ApiKey))
					if err != nil {
						if strings.Contains(err.Error(), "Client.Timeout") {
							timeoutCount++
							timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
						} else {
							formatkey := keyrexp.ReplaceAllString(err.Error(), "key=********")
							log.Printf("Line %d - failed to rename object %s: %s", i+1, name, formatkey)
						}
					}
				case "rename-service":
					var xpath string

					xpath = fmt.Sprintf("/config/devices/entry[@name='localhost.localdomain']/vsys/entry[@name='%s']/service/entry[@name='%s']", vsys, name)

					_, err := resty.R().Get(fmt.Sprintf("https://%s/api/?type=config&action=rename&xpath=%s&newname=%s&key=%s", device, xpath, value, c.ApiKey))
					if err != nil {
						if strings.Contains(err.Error(), "Client.Timeout") {
							timeoutCount++
							timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
						} else {
							formatkey := keyrexp.ReplaceAllString(err.Error(), "key=********")
							log.Printf("Line %d - failed to rename object %s: %s", i+1, name, formatkey)
						}
					}
				case "rename-servicegroup":
					var xpath string

					xpath = fmt.Sprintf("/config/devices/entry[@name='localhost.localdomain']/vsys/entry[@name='%s']/service-group/entry[@name='%s']", vsys, name)

					_, err := resty.R().Get(fmt.Sprintf("https://%s/api/?type=config&action=rename&xpath=%s&newname=%s&key=%s", device, xpath, value, c.ApiKey))
					if err != nil {
						if strings.Contains(err.Error(), "Client.Timeout") {
							timeoutCount++
							timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
						} else {
							formatkey := keyrexp.ReplaceAllString(err.Error(), "key=********")
							log.Printf("Line %d - failed to rename object %s: %s", i+1, name, formatkey)
						}
					}
				case "urlcreate":
					if len(name) > 31 {
						log.Printf("Line %d - failed to create %s: name is over the max 31 characters", i+1, name)
					} else {
						e := url.Entry{
							Name:        name,
							Description: desc,
							Sites:       urlStringToSlice(value),
							Type:        "URL List",
						}

						err = c.Objects.CustomUrlCategory.Set(vsys, e)
						if err != nil {
							if strings.Contains(err.Error(), "Client.Timeout") {
								timeoutCount++
								timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
							} else {
								log.Printf("Line %d - failed to create URL category %s: %s", i+1, name, err)
							}
						}
					}
				case "urladd":
					e := url.Entry{
						Name:        name,
						Description: desc,
						Sites:       urlStringToSlice(value),
						Type:        "URL List",
					}

					err = c.Objects.CustomUrlCategory.Set(vsys, e)
					if err != nil {
						if strings.Contains(err.Error(), "Client.Timeout") {
							timeoutCount++
							timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
						} else {
							log.Printf("Line %d - failed to update URL category %s: %s", i+1, name, err)
						}
					}

					// err = c.Objects.CustomUrlCategory.SetSite(vsys, name, value)
					// if err != nil {
					// 	log.Printf("Line %d - failed to add %s to %s: %s", i+1, value, name, err)
					// }
				case "urlremove":
					urls := urlStringToSlice(value)

					for _, url := range urls {
						err = c.Objects.CustomUrlCategory.DeleteSite(vsys, name, url)
						if err != nil {
							if strings.Contains(err.Error(), "Client.Timeout") {
								timeoutCount++
								timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
							} else {
								log.Printf("Line %d - failed to remove %s from %s: %s", i+1, value, name, err)
							}
						}
					}
				case "tag":
					if value == "delete" {
						err = c.Objects.Tags.Delete(v, name)
						if err != nil {
							if strings.Contains(err.Error(), "Client.Timeout") {
								timeoutCount++
								timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
							} else {
								log.Printf("Line %d - failed to delete %s: %s", i+1, name, err)
							}
						}
					} else {
						if len(name) > 63 {
							log.Printf("Line %d - failed to create %s: name is over the max 63 characters", i+1, name)
						} else {
							e := tags.Entry{}
							if value == "None" || value == "none" || value == "color0" || value == "" {
								e = tags.Entry{
									Name:    name,
									Comment: desc,
								}
							} else {
								e = tags.Entry{
									Name:    name,
									Color:   tag2color[value],
									Comment: desc,
								}
							}

							err = c.Objects.Tags.Set(v, e)
							if err != nil {
								if strings.Contains(err.Error(), "Client.Timeout") {
									timeoutCount++
									timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
								} else {
									log.Printf("Line %d - failed to create %s: %s", i+1, name, err)
								}
							}
						}
					}
				}

				time.Sleep(delay * time.Millisecond)
			}

			if timeoutCount > 0 {
				log.Printf("There were %d API timeout errors during import. Please verify the following have been imported/deleted/modified:\n\n", timeoutCount)
				for _, data := range timeoutData {
					info := strings.Split(data, ":")
					fmt.Printf("Line %s: Host/URL %s\n", info[0], info[1])
				}
			}
		case *pango.Panorama:
			lines, err := easycsv.Open(f)
			if err != nil {
				log.Printf("CSV file error - %s", err)
				os.Exit(1)
			}

			lc := len(lines)
			log.Printf("Running actions on %d lines - this might take a few of minutes if you have a lot of objects", lc)

			timeoutCount := 0
			timeoutData := []string{}

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
					if value == "delete" {
						err = c.Objects.Address.Delete(dgroup, name)
						if err != nil {
							if strings.Contains(err.Error(), "Client.Timeout") {
								timeoutCount++
								timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
							} else {
								log.Printf("Line %d - failed to delete %s: %s", i+1, name, err)
							}
						}
					} else {
						if len(name) > 63 {
							log.Printf("Line %d - failed to create %s: name is over the max 63 characters", i+1, name)
						} else {
							e := addr.Entry{
								Name:        name,
								Value:       strings.TrimSpace(value),
								Type:        addr.IpNetmask,
								Description: desc,
								Tags:        stringToSlice(tg),
							}

							err = c.Objects.Address.Set(dgroup, e)
							if err != nil {
								if strings.Contains(err.Error(), "Client.Timeout") {
									timeoutCount++
									timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
								} else {
									log.Printf("Line %d - failed to create %s: %s", i+1, name, err)
								}
							}
						}
					}
				case "range", "IP Range", "ip-range":
					if value == "delete" {
						err = c.Objects.Address.Delete(dgroup, name)
						if err != nil {
							if strings.Contains(err.Error(), "Client.Timeout") {
								timeoutCount++
								timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
							} else {
								log.Printf("Line %d - failed to delete %s: %s", i+1, name, err)
							}
						}
					} else {
						if len(name) > 63 {
							log.Printf("Line %d - failed to create %s: name is over the max 63 characters", i+1, name)
						} else {
							e := addr.Entry{
								Name:        name,
								Value:       strings.TrimSpace(value),
								Type:        addr.IpRange,
								Description: desc,
								Tags:        stringToSlice(tg),
							}

							err = c.Objects.Address.Set(dgroup, e)
							if err != nil {
								if strings.Contains(err.Error(), "Client.Timeout") {
									timeoutCount++
									timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
								} else {
									log.Printf("Line %d - failed to create %s: %s", i+1, name, err)
								}
							}
						}
					}
				case "fqdn", "FQDN", "Fqdn":
					if value == "delete" {
						err = c.Objects.Address.Delete(dgroup, name)
						if err != nil {
							if strings.Contains(err.Error(), "Client.Timeout") {
								timeoutCount++
								timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
							} else {
								log.Printf("Line %d - failed to delete %s: %s", i+1, name, err)
							}
						}
					} else {
						if len(name) > 63 {
							log.Printf("Line %d - failed to create %s: name is over the max 63 characters", i+1, name)
						} else {
							e := addr.Entry{
								Name:        name,
								Value:       strings.TrimSpace(value),
								Type:        addr.Fqdn,
								Description: desc,
								Tags:        stringToSlice(tg),
							}

							err = c.Objects.Address.Set(dgroup, e)
							if err != nil {
								if strings.Contains(err.Error(), "Client.Timeout") {
									timeoutCount++
									timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
								} else {
									log.Printf("Line %d - failed to create %s: %s", i+1, name, err)
								}
							}
						}
					}
				case "tcp", "udp":
					if value == "delete" {
						err = c.Objects.Services.Delete(dgroup, name)
						if err != nil {
							if strings.Contains(err.Error(), "Client.Timeout") {
								timeoutCount++
								timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
							} else {
								log.Printf("Line %d - failed to delete %s: %s", i+1, name, err)
							}
						}
					} else {
						if len(name) > 63 {
							log.Printf("Line %d - failed to create %s: name is over the max 63 characters", i+1, name)
						} else {
							e := srvc.Entry{
								Name:            name,
								Description:     desc,
								Protocol:        otype,
								DestinationPort: strings.TrimSpace(value),
								Tags:            stringToSlice(tg),
							}

							err = c.Objects.Services.Set(dgroup, e)
							if err != nil {
								if strings.Contains(err.Error(), "Client.Timeout") {
									timeoutCount++
									timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
								} else {
									log.Printf("Line %d - failed to create %s: %s", i+1, name, err)
								}
							}
						}
					}
				case "service":
					if value == "delete" {
						err = c.Objects.ServiceGroup.Delete(dgroup, name)
						if err != nil {
							if strings.Contains(err.Error(), "Client.Timeout") {
								timeoutCount++
								timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
							} else {
								log.Printf("Line %d - failed to delete %s: %s", i+1, name, err)
							}
						}
					} else {
						if len(name) > 63 {
							log.Printf("Line %d - failed to create %s: name is over the max 63 characters", i+1, name)
						} else {
							e := srvcgrp.Entry{
								Name:     name,
								Services: stringToSlice(value),
								Tags:     stringToSlice(tg),
							}

							err = c.Objects.ServiceGroup.Set(dgroup, e)
							if err != nil {
								if strings.Contains(err.Error(), "Client.Timeout") {
									timeoutCount++
									timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
								} else {
									log.Printf("Line %d - failed to create/update %s: %s", i+1, name, err)
								}
							}
						}
					}
				case "static":
					if value == "delete" {
						err = c.Objects.AddressGroup.Delete(dgroup, name)
						if err != nil {
							if strings.Contains(err.Error(), "Client.Timeout") {
								timeoutCount++
								timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
							} else {
								log.Printf("Line %d - failed to delete %s: %s", i+1, name, err)
							}
						}
					} else {
						if len(name) > 63 {
							log.Printf("Line %d - failed to create %s: name is over the max 63 characters", i+1, name)
						} else {
							groupLen := len(stringToSlice(value))

							e := addrgrp.Entry{
								Name:            name,
								Description:     desc,
								StaticAddresses: stringToSlice(value),
								Tags:            stringToSlice(tg),
							}

							err = c.Objects.AddressGroup.Set(dgroup, e)
							if err != nil {
								if strings.Contains(err.Error(), "Client.Timeout") {
									timeoutCount++
									timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
								} else {
									log.Printf("Line %d - failed to create/update %s: %s", i+1, name, err)
									if groupLen > 40 {
										log.Printf("Line %d - address group %s is over 40 members, try to add/create/breakup the group with a smaller number of members (20-30)", i+1, name)
									}
								}
							}
						}
					}
				case "dynamic":
					if value == "delete" {
						err = c.Objects.AddressGroup.Delete(dgroup, name)
						if err != nil {
							if strings.Contains(err.Error(), "Client.Timeout") {
								timeoutCount++
								timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
							} else {
								log.Printf("Line %d - failed to delete %s: %s", i+1, name, err)
							}
						}
					} else {
						if len(name) > 63 {
							log.Printf("Line %d - failed to create %s: name is over the max 63 characters", i+1, name)
						} else {
							e := addrgrp.Entry{
								Name:         name,
								Description:  desc,
								DynamicMatch: value,
								Tags:         stringToSlice(tg),
							}

							err = c.Objects.AddressGroup.Set(dgroup, e)
							if err != nil {
								if strings.Contains(err.Error(), "Client.Timeout") {
									timeoutCount++
									timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
								} else {
									log.Printf("Line %d - failed to create %s: %s", i+1, name, err)
								}
							}
						}
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
						if strings.Contains(err.Error(), "Client.Timeout") {
							timeoutCount++
							timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
						} else {
							log.Printf("Line %d - failed to update %s: %s", i+1, name, err)
						}
					}
				case "remove-service":
					if len(value) <= 0 {
						log.Printf("Line %d - you must specify a value to remove from group: %s", i+1, name)
					}

					remove := stringToSlice(value)
					cur, err := c.Objects.ServiceGroup.Get(dgroup, name)
					if err != nil {
						log.Printf("Line %d - could not retrieve object: %s", i+1, err)
					}

					newsrvcs := cur.Services

					for idx, ev := range cur.Services {
						for _, rv := range remove {
							if ev == rv {
								newsrvcs = append(newsrvcs[:idx], newsrvcs[idx+1:]...)
								break
							}
						}
					}

					e := srvcgrp.Entry{
						Name:     name,
						Services: newsrvcs,
					}

					err = c.Objects.ServiceGroup.Edit(dgroup, e)
					if err != nil {
						if strings.Contains(err.Error(), "Client.Timeout") {
							timeoutCount++
							timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
						} else {
							log.Printf("Line %d - failed to update %s: %s", i+1, name, err)
						}
					}
				case "rename-address":
					var xpath string

					if dgroup == "shared" {
						xpath = fmt.Sprintf("/config/shared/address/entry[@name='%s']", name)
					}

					if dgroup != "shared" {
						xpath = fmt.Sprintf("/config/devices/entry[@name='localhost.localdomain']/device-group/entry[@name='%s']/address/entry[@name='%s']", dgroup, name)
					}

					_, err := resty.R().Get(fmt.Sprintf("https://%s/api/?type=config&action=rename&xpath=%s&newname=%s&key=%s", device, xpath, value, c.ApiKey))
					if err != nil {
						if strings.Contains(err.Error(), "Client.Timeout") {
							timeoutCount++
							timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
						} else {
							formatkey := keyrexp.ReplaceAllString(err.Error(), "key=********")
							log.Printf("Line %d - failed to rename object %s: %s", i+1, name, formatkey)
						}
					}
				case "rename-addressgroup":
					var xpath string

					if dgroup == "shared" {
						xpath = fmt.Sprintf("/config/shared/address-group/entry[@name='%s']", name)
					}

					if dgroup != "shared" {
						xpath = fmt.Sprintf("/config/devices/entry[@name='localhost.localdomain']/device-group/entry[@name='%s']/address-group/entry[@name='%s']", dgroup, name)
					}

					_, err := resty.R().Get(fmt.Sprintf("https://%s/api/?type=config&action=rename&xpath=%s&newname=%s&key=%s", device, xpath, value, c.ApiKey))
					if err != nil {
						if strings.Contains(err.Error(), "Client.Timeout") {
							timeoutCount++
							timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
						} else {
							formatkey := keyrexp.ReplaceAllString(err.Error(), "key=********")
							log.Printf("Line %d - failed to rename object %s: %s", i+1, name, formatkey)
						}
					}
				case "rename-service":
					var xpath string

					if dgroup == "shared" {
						xpath = fmt.Sprintf("/config/shared/service/entry[@name='%s']", name)
					}

					if dgroup != "shared" {
						xpath = fmt.Sprintf("/config/devices/entry[@name='localhost.localdomain']/device-group/entry[@name='%s']/service/entry[@name='%s']", dgroup, name)
					}

					_, err := resty.R().Get(fmt.Sprintf("https://%s/api/?type=config&action=rename&xpath=%s&newname=%s&key=%s", device, xpath, value, c.ApiKey))
					if err != nil {
						if strings.Contains(err.Error(), "Client.Timeout") {
							timeoutCount++
							timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
						} else {
							formatkey := keyrexp.ReplaceAllString(err.Error(), "key=********")
							log.Printf("Line %d - failed to rename object %s: %s", i+1, name, formatkey)
						}
					}
				case "rename-servicegroup":
					var xpath string

					if dgroup == "shared" {
						xpath = fmt.Sprintf("/config/shared/service-group/entry[@name='%s']", name)
					}

					if dgroup != "shared" {
						xpath = fmt.Sprintf("/config/devices/entry[@name='localhost.localdomain']/device-group/entry[@name='%s']/service-group/entry[@name='%s']", dgroup, name)
					}

					_, err := resty.R().Get(fmt.Sprintf("https://%s/api/?type=config&action=rename&xpath=%s&newname=%s&key=%s", device, xpath, value, c.ApiKey))
					if err != nil {
						if strings.Contains(err.Error(), "Client.Timeout") {
							timeoutCount++
							timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
						} else {
							formatkey := keyrexp.ReplaceAllString(err.Error(), "key=********")
							log.Printf("Line %d - failed to rename object %s: %s", i+1, name, formatkey)
						}
					}
				case "urlcreate":
					if len(name) > 31 {
						log.Printf("Line %d - failed to create %s: name is over the max 31 characters", i+1, name)
					} else {
						e := url.Entry{
							Name:        name,
							Description: desc,
							Sites:       urlStringToSlice(value),
							Type:        "URL List",
						}

						err = c.Objects.CustomUrlCategory.Set(dgroup, e)
						if err != nil {
							if strings.Contains(err.Error(), "Client.Timeout") {
								timeoutCount++
								timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
							} else {
								log.Printf("Line %d - failed to create URL category %s: %s", i+1, name, err)
							}
						}
					}
				case "urladd":
					e := url.Entry{
						Name:        name,
						Description: desc,
						Sites:       urlStringToSlice(value),
						Type:        "URL List",
					}

					err = c.Objects.CustomUrlCategory.Set(dgroup, e)
					if err != nil {
						if strings.Contains(err.Error(), "Client.Timeout") {
							timeoutCount++
							timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
						} else {
							log.Printf("Line %d - failed to update URL category %s: %s", i+1, name, err)
						}
					}

					// err = c.Objects.CustomUrlCategory.SetSite(dgroup, name, value)
					// if err != nil {
					// 	log.Printf("Line %d - failed to add %s to %s: %s", i+1, value, name, err)
					// }
				case "urlremove":
					urls := urlStringToSlice(value)

					for _, url := range urls {
						err = c.Objects.CustomUrlCategory.DeleteSite(dgroup, name, url)
						if err != nil {
							if strings.Contains(err.Error(), "Client.Timeout") {
								timeoutCount++
								timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
							} else {
								log.Printf("Line %d - failed to remove %s from %s: %s", i+1, value, name, err)
							}
						}
					}
				case "tag":
					if value == "delete" {
						err = c.Objects.Tags.Delete(dgroup, name)
						if err != nil {
							if strings.Contains(err.Error(), "Client.Timeout") {
								timeoutCount++
								timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
							} else {
								log.Printf("Line %d - failed to delete %s: %s", i+1, name, err)
							}
						}
					} else {
						if len(name) > 63 {
							log.Printf("Line %d - failed to create %s: name is over the max 63 characters", i+1, name)
						} else {
							e := tags.Entry{}
							if value == "None" || value == "none" || value == "color0" || value == "" {
								e = tags.Entry{
									Name:    name,
									Comment: desc,
								}
							} else {
								e = tags.Entry{
									Name:    name,
									Color:   tag2color[value],
									Comment: desc,
								}
							}

							err = c.Objects.Tags.Set(dgroup, e)
							if err != nil {
								if strings.Contains(err.Error(), "Client.Timeout") {
									timeoutCount++
									timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, name))
								} else {
									log.Printf("Line %d - failed to create %s: %s", i+1, name, err)
								}
							}
						}
					}
				}

				time.Sleep(delay * time.Millisecond)
			}

			if timeoutCount > 0 {
				log.Printf("There were %d API timeout errors during import. Please verify the following have been imported/deleted/modified:\n\n", timeoutCount)
				for _, data := range timeoutData {
					info := strings.Split(data, ":")
					fmt.Printf("Line %s: Host/URL %s\n", info[0], info[1])
				}
			}
		}
	},
}

func init() {
	objectsCmd.AddCommand(objectsImportCmd)

	objectsImportCmd.Flags().StringVarP(&user, "user", "u", "", "User to connect to the device as")
	objectsImportCmd.Flags().StringVarP(&p, "delay", "p", "100", "Delay (in milliseconds) to pause between each API call")
	objectsImportCmd.Flags().StringVarP(&device, "device", "d", "", "Device to connect to")
	objectsImportCmd.Flags().StringVarP(&f, "file", "f", "", "Name of the CSV file to import")
	// objectsImportCmd.Flags().StringVarP(&dg, "devicegroup", "g", "shared", "Device Group name when exporting from Panorama")
	// objectsImportCmd.Flags().StringVarP(&v, "vsys", "v", "vsys1", "Vsys name when exporting from a firewall")
	objectsImportCmd.MarkFlagRequired("user")
	// objectsImportCmd.MarkFlagRequired("pass")
	objectsImportCmd.MarkFlagRequired("device")
	objectsImportCmd.MarkFlagRequired("file")
}
