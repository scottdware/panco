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
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	easycsv "github.com/scottdware/go-easycsv"
	"github.com/spf13/cobra"
)

// cliCmd represents the cli command
var objectsCliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Generate CLI set commands from a CSV file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		lines, err := easycsv.Open(f)
		if err != nil {
			log.Printf("CSV file error - %s", err)
			os.Exit(1)
		}

		txtfile, err := os.Create(txt)
		if err != nil {
			log.Printf("TXT file error - %s", err)
			os.Exit(1)
		}

		defer txtfile.Close()

		lc := len(lines)
		log.Printf("Converting %d lines - saving results to '%s'", lc, txt)

		for i, line := range lines {
			var command string
			name := line[0]
			otype := line[1]
			value := line[2]
			desc := line[3]
			tg := line[4]
			devtype := line[5]

			switch otype {
			case "ip", "IP Netmask", "ip-netmask":
				if value == "delete" {
					switch devtype {
					case "vsys1":
						command = fmt.Sprintf("delete address %s\n", name)
					case "shared":
						command = fmt.Sprintf("delete shared address %s\n", name)
					default:
						command = fmt.Sprintf("delete device-group %s address %s\n", devtype, name)
					}
				} else {
					switch devtype {
					case "vsys1":
						command = fmt.Sprintf("set address %s ip-netmask %s\n", name, value)

						if desc != "" {
							command += fmt.Sprintf("set address %s description \"%s\"\n", name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("set address %s tag [ %s ]\n", name, strings.Join(tags, " "))
						}
					case "shared":
						command = fmt.Sprintf("set shared address %s ip-netmask %s\n", name, value)

						if desc != "" {
							command += fmt.Sprintf("set shared address %s description \"%s\"\n", name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("set shared address %s tag [ %s ]\n", name, strings.Join(tags, " "))
						}
					default:
						command = fmt.Sprintf("set device-group %s address %s ip-netmask %s\n", devtype, name, value)

						if desc != "" {
							command += fmt.Sprintf("set device-group %s address %s description \"%s\"\n", devtype, name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("set device-group %s address %s tag [ %s ]\n", devtype, name, strings.Join(tags, " "))
						}
					}
				}

				// fmt.Printf("%s", command)
				_, err = io.WriteString(txtfile, command)

				if err != nil {
					log.Printf("Failed to write to TXT file - %s", err)
				}
			case "range", "IP Range", "ip-range":
				if value == "delete" {
					switch devtype {
					case "vsys1":
						command = fmt.Sprintf("delete address %s\n", name)
					case "shared":
						command = fmt.Sprintf("delete shared address %s\n", name)
					default:
						command = fmt.Sprintf("delete device-group %s address %s\n", devtype, name)
					}
				} else {
					switch devtype {
					case "vsys1":
						command = fmt.Sprintf("set address %s ip-range %s\n", name, value)

						if desc != "" {
							command += fmt.Sprintf("set address %s description \"%s\"\n", name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("set address %s tag [ %s ]\n", name, strings.Join(tags, " "))
						}
					case "shared":
						command = fmt.Sprintf("set shared address %s ip-range %s\n", name, value)

						if desc != "" {
							command += fmt.Sprintf("set shared address %s description \"%s\"\n", name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("set shared address %s tag [ %s ]\n", name, strings.Join(tags, " "))
						}
					default:
						command = fmt.Sprintf("set device-group %s address %s ip-range %s\n", devtype, name, value)

						if desc != "" {
							command += fmt.Sprintf("set device-group %s address %s description \"%s\"\n", devtype, name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("set device-group %s address %s tag [ %s ]\n", devtype, name, strings.Join(tags, " "))
						}
					}
				}

				// fmt.Printf("%s", command)
				_, err = io.WriteString(txtfile, command)

				if err != nil {
					log.Printf("Failed to write to TXT file - %s", err)
				}
			case "fqdn", "FQDN", "Fqdn":
				if value == "delete" {
					switch devtype {
					case "vsys1":
						command = fmt.Sprintf("delete address %s\n", name)
					case "shared":
						command = fmt.Sprintf("delete shared address %s\n", name)
					default:
						command = fmt.Sprintf("delete device-group %s address %s\n", devtype, name)
					}
				} else {
					switch devtype {
					case "vsys1":
						command = fmt.Sprintf("set address %s fqdn %s\n", name, value)

						if desc != "" {
							command += fmt.Sprintf("set address %s description \"%s\"\n", name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("set address %s tag [ %s ]\n", name, strings.Join(tags, " "))
						}
					case "shared":
						command = fmt.Sprintf("set shared address %s fqdn %s\n", name, value)

						if desc != "" {
							command += fmt.Sprintf("set shared address %s description \"%s\"\n", name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("set shared address %s tag [ %s ]\n", name, strings.Join(tags, " "))
						}
					default:
						command = fmt.Sprintf("set device-group %s address %s fqdn %s\n", devtype, name, value)

						if desc != "" {
							command += fmt.Sprintf("set device-group %s address %s description \"%s\"\n", devtype, name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("set device-group %s address %s tag [ %s ]\n", devtype, name, strings.Join(tags, " "))
						}
					}
				}

				// fmt.Printf("%s", command)
				_, err = io.WriteString(txtfile, command)

				if err != nil {
					log.Printf("Failed to write to TXT file - %s", err)
				}
			case "tcp", "udp":
				if value == "delete" {
					switch devtype {
					case "vsys1":
						command = fmt.Sprintf("delete service %s\n", name)
					case "shared":
						command = fmt.Sprintf("delete shared service %s\n", name)
					default:
						command = fmt.Sprintf("delete device-group %s service %s\n", devtype, name)
					}
				} else {
					switch devtype {
					case "vsys1":
						command = fmt.Sprintf("set service %s protocol %s port %s\n", name, otype, value)

						if desc != "" {
							command += fmt.Sprintf("set service %s description \"%s\"\n", name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("set service %s tag [ %s ]\n", name, strings.Join(tags, " "))
						}
					case "shared":
						command = fmt.Sprintf("set shared service %s protocol %s port %s\n", name, otype, value)

						if desc != "" {
							command += fmt.Sprintf("set shared service %s description \"%s\"\n", name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("set shared service %s tag [ %s ]\n", name, strings.Join(tags, " "))
						}
					default:
						command = fmt.Sprintf("set device-group %s service %s protocol %s port %s\n", devtype, name, otype, value)

						if desc != "" {
							command += fmt.Sprintf("set device-group %s service %s description \"%s\"\n", devtype, name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("set device-group %s service %s tag [ %s ]\n", devtype, name, strings.Join(tags, " "))
						}
					}
				}

				// fmt.Printf("%s", command)
				_, err = io.WriteString(txtfile, command)

				if err != nil {
					log.Printf("Failed to write to TXT file - %s", err)
				}
			case "service":
				if value == "delete" {
					switch devtype {
					case "vsys1":
						command = fmt.Sprintf("delete service-group %s\n", name)
					case "shared":
						command = fmt.Sprintf("delete shared service-group %s\n", name)
					default:
						command = fmt.Sprintf("delete device-group %s service-group %s\n", devtype, name)
					}
				} else {
					switch devtype {
					case "vsys1":
						members := stringToSlice(value)
						command = fmt.Sprintf("set service-group %s members [ %s ]\n", name, strings.Join(members, " "))

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("set service-group %s tag [ %s ]\n", name, strings.Join(tags, " "))
						}
					case "shared":
						members := stringToSlice(value)
						command = fmt.Sprintf("set shared service-group %s members [ %s ]\n", name, strings.Join(members, " "))

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("set shared service-group %s tag [ %s ]\n", name, strings.Join(tags, " "))
						}
					default:
						members := stringToSlice(value)
						command = fmt.Sprintf("set device-group %s service-group %s members [ %s ]\n", devtype, name, strings.Join(members, " "))

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("set device-group %s service-group %s tag [ %s ]\n", devtype, name, strings.Join(tags, " "))
						}
					}
				}

				// fmt.Printf("%s", command)
				_, err = io.WriteString(txtfile, command)

				if err != nil {
					log.Printf("Failed to write to TXT file - %s", err)
				}
			case "static":
				if value == "delete" {
					switch devtype {
					case "vsys1":
						command = fmt.Sprintf("delete address-group %s\n", name)
					case "shared":
						command = fmt.Sprintf("delete shared address-group %s\n", name)
					default:
						command = fmt.Sprintf("delete device-group %s address-group %s\n", devtype, name)
					}
				} else {
					switch devtype {
					case "vsys1":
						members := stringToSlice(value)
						command = fmt.Sprintf("set address-group %s static [ %s ]\n", name, strings.Join(members, " "))

						if desc != "" {
							command += fmt.Sprintf("set address-group %s description \"%s\"\n", name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("set address-group %s tag [ %s ]\n", name, strings.Join(tags, " "))
						}
					case "shared":
						members := stringToSlice(value)
						command = fmt.Sprintf("set shared address-group %s static [ %s ]\n", name, strings.Join(members, " "))

						if desc != "" {
							command += fmt.Sprintf("set shared address-group %s description \"%s\"\n", name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("set shared address-group %s tag [ %s ]\n", name, strings.Join(tags, " "))
						}
					default:
						members := stringToSlice(value)
						command = fmt.Sprintf("set device-group %s address-group %s static [ %s ]\n", devtype, name, strings.Join(members, " "))

						if desc != "" {
							command += fmt.Sprintf("set device-group %s address-group %s description \"%s\"\n", devtype, name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("set device-group %s address-group %s tag [ %s ]\n", devtype, name, strings.Join(tags, " "))
						}
					}
				}

				// fmt.Printf("%s", command)
				_, err = io.WriteString(txtfile, command)

				if err != nil {
					log.Printf("Failed to write to TXT file - %s", err)
				}
			case "dynamic":
				if value == "delete" {
					switch devtype {
					case "vsys1":
						command = fmt.Sprintf("delete address-group %s\n", name)
					case "shared":
						command = fmt.Sprintf("delete shared address-group %s\n", name)
					default:
						command = fmt.Sprintf("delete device-group %s address-group %s\n", devtype, name)
					}
				} else {
					switch devtype {
					case "vsys1":
						command = fmt.Sprintf("set address-group %s dynamic filter \"%s\"\n", name, value)

						if desc != "" {
							command += fmt.Sprintf("set address-group %s description \"%s\"\n", name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("set address-group %s tag [ %s ]\n", name, strings.Join(tags, " "))
						}
					case "shared":
						command = fmt.Sprintf("set shared address-group %s dynamic filter \"%s\"\n", name, value)

						if desc != "" {
							command += fmt.Sprintf("set shared address-group %s description \"%s\"\n", name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("set shared address-group %s tag [ %s ]\n", name, strings.Join(tags, " "))
						}
					default:
						command = fmt.Sprintf("set device-group %s address-group %s dynamic filter \"%s\"\n", devtype, name, value)

						if desc != "" {
							command += fmt.Sprintf("set device-group %s address-group %s description \"%s\"\n", devtype, name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("set device-group %s address-group %s tag [ %s ]\n", devtype, name, strings.Join(tags, " "))
						}
					}
				}

				// fmt.Printf("%s", command)
				_, err = io.WriteString(txtfile, command)

				if err != nil {
					log.Printf("Failed to write to TXT file - %s", err)
				}
			case "remove-address":
				if len(value) <= 0 {
					log.Printf("Line %d - you must specify a value to remove from group: %s", i+1, name)
				}

				switch devtype {
				case "vsys1":
					command = fmt.Sprintf("delete address-group %s static %s\n", name, value)
				case "shared":
					command = fmt.Sprintf("delete shared address-group %s static %s\n", name, value)
				default:
					command = fmt.Sprintf("delete device-group %s address-group %s static %s\n", devtype, name, value)
				}

				// fmt.Printf("%s", command)
				_, err = io.WriteString(txtfile, command)

				if err != nil {
					log.Printf("Failed to write to TXT file - %s\n", err)
				}
			case "remove-service":
				if len(value) <= 0 {
					log.Printf("Line %d - you must specify a value to remove from group: %s\n", i+1, name)
				}

				switch devtype {
				case "vsys1":
					command = fmt.Sprintf("delete service-group %s members %s\n", name, value)
				case "shared":
					command = fmt.Sprintf("delete shared service-group %s members %s\n", name, value)
				default:
					command = fmt.Sprintf("delete device-group %s service-group %s members %s\n", devtype, name, value)
				}

				// fmt.Printf("%s", command)
				_, err = io.WriteString(txtfile, command)

				if err != nil {
					log.Printf("Failed to write to TXT file - %s", err)
				}
			case "rename-address":
				switch devtype {
				case "vsys1":
					command = fmt.Sprintf("\nrename address %s to %s\n", name, value)
				case "shared":
					command = fmt.Sprintf("\nrename shared address %s to %s\n", name, value)
				default:
					command = fmt.Sprintf("\nrename device-group %s address %s to %s\n", devtype, name, value)
				}

				// fmt.Printf("%s", command)
				_, err = io.WriteString(txtfile, command)

				if err != nil {
					log.Printf("Failed to write to TXT file - %s", err)
				}
			case "rename-addressgroup":
				switch devtype {
				case "vsys1":
					command = fmt.Sprintf("\nrename address-group %s to %s\n", name, value)
				case "shared":
					command = fmt.Sprintf("\nrename shared address-group %s to %s\n", name, value)
				default:
					command = fmt.Sprintf("\nrename device-group %s address-group %s to %s\n", devtype, name, value)
				}

				// fmt.Printf("%s", command)
				_, err = io.WriteString(txtfile, command)

				if err != nil {
					log.Printf("Failed to write to TXT file - %s", err)
				}
			case "rename-service":
				switch devtype {
				case "vsys1":
					command = fmt.Sprintf("\nrename service %s to %s\n", name, value)
				case "shared":
					command = fmt.Sprintf("\nrename shared service %s to %s\n", name, value)
				default:
					command = fmt.Sprintf("\nrename device-group %s service %s to %s\n", devtype, name, value)
				}

				// fmt.Printf("%s", command)
				_, err = io.WriteString(txtfile, command)

				if err != nil {
					log.Printf("Failed to write to TXT file - %s", err)
				}
			case "rename-servicegroup":
				switch devtype {
				case "vsys1":
					command = fmt.Sprintf("\nrename service-group %s to %s\n", name, value)
				case "shared":
					command = fmt.Sprintf("\nrename shared service-group %s to %s\n", name, value)
				default:
					command = fmt.Sprintf("\nrename device-group %s service-group %s to %s\n", devtype, name, value)
				}

				// fmt.Printf("%s", command)
				_, err = io.WriteString(txtfile, command)

				if err != nil {
					log.Printf("Failed to write to TXT file - %s", err)
				}
			case "tag":
				if value == "delete" {
					switch devtype {
					case "vsys1":
						command = fmt.Sprintf("delete tag %s\n", name)
					case "shared":
						command = fmt.Sprintf("delete shared tag %s\n", name)
					default:
						command = fmt.Sprintf("delete device-group %s tag %s\n", devtype, name)
					}
				} else {
					if value == "None" || value == "none" || value == "color0" || value == "" {
						switch devtype {
						case "vsys1":
							command = fmt.Sprintf("set tag %s color color0\n", name)

							if desc != "" {
								command += fmt.Sprintf("set tag %s comments \"%s\"\n", name, desc)
							}
						case "shared":
							command = fmt.Sprintf("set shared tag %s color color0\n", name)

							if desc != "" {
								command += fmt.Sprintf("set shared tag %s comments \"%s\"\n", name, desc)
							}
						default:
							command = fmt.Sprintf("set device-group %s tag %s color color0\n", devtype, name)

							if desc != "" {
								command += fmt.Sprintf("set device-group %s tag %s comments \"%s\"\n", devtype, name, desc)
							}
						}
					} else {
						switch devtype {
						case "vsys1":
							command = fmt.Sprintf("set tag %s color %s\n", name, tag2color[value])

							if desc != "" {
								command += fmt.Sprintf("set tag %s comments \"%s\"\n", name, desc)
							}
						case "shared":
							command = fmt.Sprintf("set shared tag %s color %s\n", name, tag2color[value])

							if desc != "" {
								command += fmt.Sprintf("set shared tag %s comments \"%s\"\n", name, desc)
							}
						default:
							command = fmt.Sprintf("set device-group %s tag %s color %s\n", devtype, name, tag2color[value])

							if desc != "" {
								command += fmt.Sprintf("set device-group %s tag %s comments \"%s\"\n", devtype, name, desc)
							}
						}
					}
				}

				// fmt.Printf("%s", command)
				_, err = io.WriteString(txtfile, command)

				if err != nil {
					log.Printf("Failed to write to TXT file - %s", err)
				}
			case "urladd":
				var cmdBody string

				urls := stringToSlice(value)

				for _, url := range urls {
					cmdBody += fmt.Sprintf("%s ", url)
				}

				switch devtype {
				case "vsys1":
					// command = fmt.Printf("set profiles custom-url-category type \"URL List\"\n")
					command = fmt.Sprintf("set profiles custom-url-category %s list [ %s ]\n", name, strings.TrimRight(cmdBody, " "))
				case "shared":
					// command = fmt.Printf("set shared profiles custom-url-category type \"URL List\"\n")
					command = fmt.Sprintf("set shared profiles custom-url-category %s list [ %s ]\n", name, strings.TrimRight(cmdBody, " "))
				default:
					// command = fmt.Sprintf("set device-group %s profiles custom-url-category type \"URL List\"\n", name)
					command = fmt.Sprintf("set device-group %s profiles custom-url-category %s list [ %s ]\n", devtype, name, strings.TrimRight(cmdBody, " "))
				}

				// fmt.Printf("%s", command)
				_, err = io.WriteString(txtfile, command)

				if err != nil {
					log.Printf("Failed to write to TXT file - %s", err)
				}
			case "urlremove":
				switch devtype {
				case "vsys1":
					command = fmt.Sprintf("delete profiles custom-url-category %s list %s\n", name, value)
				case "shared":
					command = fmt.Sprintf("delete shared profiles custom-url-category %s list %s\n", name, value)
				default:
					command = fmt.Sprintf("delete device-group %s profiles custom-url-category %s list %s\n", devtype, name, value)
				}

				// fmt.Printf("%s", command)
				_, err = io.WriteString(txtfile, command)

				if err != nil {
					log.Printf("Failed to write to TXT file - %s", err)
				}
			}

			txtfile.Sync()
		}
	},
}

func init() {
	objectsCmd.AddCommand(objectsCliCmd)

	objectsCliCmd.Flags().StringVarP(&f, "file", "f", "", "Name of the CSV file to convert")
	objectsCliCmd.Flags().StringVarP(&txt, "output", "o", "", "Name of the file to output SET commands to")
	objectsCliCmd.MarkFlagRequired("file")
	objectsCliCmd.MarkFlagRequired("output")
}
