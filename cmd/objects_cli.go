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
	"log"
	"os"
	"strings"

	easycsv "github.com/scottdware/go-easycsv"
	"github.com/spf13/cobra"
)

// cliCmd represents the cli command
var objectsCliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Convert CSV file entries to CLI set commands",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		lines, err := easycsv.Open(f)
		if err != nil {
			log.Printf("CSV file error - %s", err)
			os.Exit(1)
		}

		lc := len(lines)
		log.Printf("Converting %d lines - this might take a few of minutes if you have a lot of objects", lc)

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
					if strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\ndelete address %s", name)
					}

					if !strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\ndelete device-group %s address %s", devtype, name)
					}
				} else {
					if strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\nset address %s ip-netmask %s", name, value)

						if desc != "" {
							command += fmt.Sprintf("\nset address %s description \"%s\"", name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("\nset address %s tag [ %s ]", name, strings.Join(tags, " "))
						}
					}

					if !strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\nset device-group %s address %s ip-netmask %s", devtype, name, value)

						if desc != "" {
							command += fmt.Sprintf("\nset device-group %s address %s description \"%s\"", devtype, name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("\nset device-group %s address %s tag [ %s ]", devtype, name, strings.Join(tags, " "))
						}
					}
				}

				fmt.Printf("%s", command)
			case "range", "IP Range", "ip-range":
				if value == "delete" {
					if strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\ndelete address %s", name)
					}

					if !strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\ndelete device-group %s address %s", devtype, name)
					}
				} else {
					if strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\nset address %s ip-range %s", name, value)

						if desc != "" {
							command += fmt.Sprintf("\nset address %s description \"%s\"", name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("\nset address %s tag [ %s ]", name, strings.Join(tags, " "))
						}
					}

					if !strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\nset device-group %s address %s ip-range %s", devtype, name, value)

						if desc != "" {
							command += fmt.Sprintf("\nset device-group %s address %s description \"%s\"", devtype, name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("\nset device-group %s address %s tag [ %s ]", devtype, name, strings.Join(tags, " "))
						}
					}
				}

				fmt.Printf("%s", command)
			case "fqdn", "FQDN", "Fqdn":
				if value == "delete" {
					if strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\ndelete address %s", name)
					}

					if !strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\ndelete device-group %s address %s", devtype, name)
					}
				} else {
					if strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\nset address %s fqdn %s", name, value)

						if desc != "" {
							command += fmt.Sprintf("\nset address %s description \"%s\"", name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("\nset address %s tag [ %s ]", name, strings.Join(tags, " "))
						}
					}

					if !strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\nset device-group %s address %s fqdn %s", devtype, name, value)

						if desc != "" {
							command += fmt.Sprintf("\nset device-group %s address %s description \"%s\"", devtype, name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("\nset device-group %s address %s tag [ %s ]", devtype, name, strings.Join(tags, " "))
						}
					}
				}

				fmt.Printf("%s", command)
			case "tcp", "udp":
				if value == "delete" {
					if strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\ndelete service %s", name)
					}

					if !strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\ndelete device-group %s service %s", devtype, name)
					}
				} else {
					if strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\nset service %s protocol %s %s", name, otype, value)

						if desc != "" {
							command += fmt.Sprintf("\nset service %s description \"%s\"", name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("\nset service %s tag [ %s ]", name, strings.Join(tags, " "))
						}
					}

					if !strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\nset device-group %s service %s protocol %s %s", devtype, name, otype, value)

						if desc != "" {
							command += fmt.Sprintf("\nset device-group %s service %s description \"%s\"", devtype, name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("\nset device-group %s service %s tag [ %s ]", devtype, name, strings.Join(tags, " "))
						}
					}
				}

				fmt.Printf("%s", command)
			case "service":
				if value == "delete" {
					if strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\ndelete service-group %s", name)
					}

					if !strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\ndelete device-group %s service-group %s", devtype, name)
					}
				} else {
					if strings.Contains(devtype, "vsys") {
						members := stringToSlice(value)
						command = fmt.Sprintf("\nset service-group %s members [ %s ]", name, strings.Join(members, " "))

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("\nset service-group %s tag [ %s ]", name, strings.Join(tags, " "))
						}
					}

					if !strings.Contains(devtype, "vsys") {
						members := stringToSlice(value)
						command = fmt.Sprintf("\nset device-group %s service-group %s members [ %s ]", devtype, name, strings.Join(members, " "))

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("\nset device-group %s service-group %s tag [ %s ]", devtype, name, strings.Join(tags, " "))
						}
					}
				}

				fmt.Printf("%s", command)
			case "static":
				if value == "delete" {
					if strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\ndelete address-group %s", name)
					}

					if !strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\ndelete device-group %s address-group %s", devtype, name)
					}
				} else {
					if strings.Contains(devtype, "vsys") {
						members := stringToSlice(value)
						command = fmt.Sprintf("\nset address-group %s static [ %s ]", name, strings.Join(members, " "))

						if desc != "" {
							command += fmt.Sprintf("\nset address-group %s description \"%s\"", name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("\nset address-group %s tag [ %s ]", name, strings.Join(tags, " "))
						}
					}

					if !strings.Contains(devtype, "vsys") {
						members := stringToSlice(value)
						command = fmt.Sprintf("\nset device-group %s service-group %s static [ %s ]", devtype, name, strings.Join(members, " "))

						if desc != "" {
							command += fmt.Sprintf("\nset device-group %s address-group %s description \"%s\"", devtype, name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("\nset device-group %s address-group %s tag [ %s ]", devtype, name, strings.Join(tags, " "))
						}
					}
				}

				fmt.Printf("%s", command)
			case "dynamic":
				if value == "delete" {
					if strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\ndelete address-group %s", name)
					}

					if !strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\ndelete device-group %s address-group %s", devtype, name)
					}
				} else {
					if strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\nset address-group %s dynamic filter \"%s\"", name, value)

						if desc != "" {
							command += fmt.Sprintf("\nset address-group %s description \"%s\"", name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("\nset address-group %s tag [ %s ]", name, strings.Join(tags, " "))
						}
					}

					if !strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\nset device-group %s address-group %s dynamic filter \"%s\"", devtype, name, value)

						if desc != "" {
							command += fmt.Sprintf("\nset device-group %s address-group %s description \"%s\"", devtype, name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("\nset device-group %s address-group %s tag [ %s ]", devtype, name, strings.Join(tags, " "))
						}
					}
				}

				fmt.Printf("%s", command)
			case "remove-address":
				if len(value) <= 0 {
					log.Printf("Line %d - you must specify a value to remove from group: %s", i+1, name)
				}

				// remove := stringToSlice(value)

			case "remove-service":
				if len(value) <= 0 {
					log.Printf("Line %d - you must specify a value to remove from group: %s", i+1, name)
				}

				// remove := stringToSlice(value)

			case "rename-address":

			case "rename-addressgroup":

			case "rename-service":

			case "rename-servicegroup":

			case "tag":
				if value == "delete" {

				} else {
					if value == "None" || value == "none" || value == "color0" || value == "" {

					} else {

					}

				}
			}
		}
	},
}

func init() {
	objectsCmd.AddCommand(objectsCliCmd)

	objectsCliCmd.Flags().StringVarP(&f, "file", "f", "", "Name of the CSV file to convert")

	objectsCliCmd.MarkFlagRequired("file")
}

func genDeleteCli(name, obj, devtype string) string {
	var command string

	if strings.Contains(devtype, "vsys") {
		command = fmt.Sprintf("\ndelete %s %s", obj, name)
	}

	if !strings.Contains(devtype, "vsys") {
		command = fmt.Sprintf("\ndelete device-group %s %s %s", devtype, obj, name)
	}

	return command
}
