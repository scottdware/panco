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
					if strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\ndelete address %s", name)
					}

					if strings.ToLower(devtype) == "shared" {
						command = fmt.Sprintf("\ndelete shared address %s", name)
					}

					if !strings.Contains(devtype, "vsys") || !strings.Contains(devtype, "hared") {
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

					if strings.ToLower(devtype) == "shared" {
						command = fmt.Sprintf("\nset shared address %s ip-netmask %s", name, value)

						if desc != "" {
							command += fmt.Sprintf("\nset shared address %s description \"%s\"", name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("\nset shared address %s tag [ %s ]", name, strings.Join(tags, " "))
						}
					}

					if !strings.Contains(devtype, "vsys") || !strings.Contains(devtype, "hared") {
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

				// fmt.Printf("%s", command)
				_, err = io.WriteString(txtfile, command)

				if err != nil {
					log.Printf("Failed to write to TXT file - %s", err)
				}
			case "range", "IP Range", "ip-range":
				if value == "delete" {
					if strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\ndelete address %s", name)
					}

					if strings.ToLower(devtype) == "shared" {
						command = fmt.Sprintf("\ndelete shared address %s", name)
					}

					if !strings.Contains(devtype, "vsys") || !strings.Contains(devtype, "hared") {
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

					if strings.ToLower(devtype) == "shared" {
						command = fmt.Sprintf("\nset shared address %s ip-range %s", name, value)

						if desc != "" {
							command += fmt.Sprintf("\nset shared address %s description \"%s\"", name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("\nset shared address %s tag [ %s ]", name, strings.Join(tags, " "))
						}
					}

					if !strings.Contains(devtype, "vsys") || !strings.Contains(devtype, "hared") {
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

				// fmt.Printf("%s", command)
				_, err = io.WriteString(txtfile, command)

				if err != nil {
					log.Printf("Failed to write to TXT file - %s", err)
				}
			case "fqdn", "FQDN", "Fqdn":
				if value == "delete" {
					if strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\ndelete address %s", name)
					}

					if strings.ToLower(devtype) == "shared" {
						command = fmt.Sprintf("\ndelete shared address %s", name)
					}

					if !strings.Contains(devtype, "vsys") || !strings.Contains(devtype, "hared") {
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

					if strings.ToLower(devtype) == "shared" {
						command = fmt.Sprintf("\nset shared address %s fqdn %s", name, value)

						if desc != "" {
							command += fmt.Sprintf("\nset shared address %s description \"%s\"", name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("\nset shared address %s tag [ %s ]", name, strings.Join(tags, " "))
						}
					}

					if !strings.Contains(devtype, "vsys") || !strings.Contains(devtype, "hared") {
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

				// fmt.Printf("%s", command)
				_, err = io.WriteString(txtfile, command)

				if err != nil {
					log.Printf("Failed to write to TXT file - %s", err)
				}
			case "tcp", "udp":
				if value == "delete" {
					if strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\ndelete service %s", name)
					}

					if strings.ToLower(devtype) == "shared" {
						command = fmt.Sprintf("\ndelete shared service %s", name)
					}

					if !strings.Contains(devtype, "vsys") || !strings.Contains(devtype, "hared") {
						command = fmt.Sprintf("\ndelete device-group %s service %s", devtype, name)
					}
				} else {
					if strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\nset service %s protocol %s port %s", name, otype, value)

						if desc != "" {
							command += fmt.Sprintf("\nset service %s description \"%s\"", name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("\nset service %s tag [ %s ]", name, strings.Join(tags, " "))
						}
					}

					if strings.ToLower(devtype) == "shared" {
						command = fmt.Sprintf("\nset shared service %s protocol %s port %s", name, otype, value)

						if desc != "" {
							command += fmt.Sprintf("\nset shared service %s description \"%s\"", name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("\nset shared service %s tag [ %s ]", name, strings.Join(tags, " "))
						}
					}

					if !strings.Contains(devtype, "vsys") || !strings.Contains(devtype, "hared") {
						command = fmt.Sprintf("\nset device-group %s service %s protocol %s port %s", devtype, name, otype, value)

						if desc != "" {
							command += fmt.Sprintf("\nset device-group %s service %s description \"%s\"", devtype, name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("\nset device-group %s service %s tag [ %s ]", devtype, name, strings.Join(tags, " "))
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
					if strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\ndelete service-group %s", name)
					}

					if strings.ToLower(devtype) == "shared" {
						command = fmt.Sprintf("\ndelete shared service-group %s", name)
					}

					if !strings.Contains(devtype, "vsys") || !strings.Contains(devtype, "hared") {
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

					if strings.ToLower(devtype) == "shared" {
						members := stringToSlice(value)
						command = fmt.Sprintf("\nset shared service-group %s members [ %s ]", name, strings.Join(members, " "))

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("\nset shared service-group %s tag [ %s ]", name, strings.Join(tags, " "))
						}
					}

					if !strings.Contains(devtype, "vsys") || !strings.Contains(devtype, "hared") {
						members := stringToSlice(value)
						command = fmt.Sprintf("\nset device-group %s service-group %s members [ %s ]", devtype, name, strings.Join(members, " "))

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("\nset device-group %s service-group %s tag [ %s ]", devtype, name, strings.Join(tags, " "))
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
					if strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\ndelete address-group %s", name)
					}

					if strings.ToLower(devtype) == "shared" {
						command = fmt.Sprintf("\ndelete shared address-group %s", name)
					}

					if !strings.Contains(devtype, "vsys") || !strings.Contains(devtype, "hared") {
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

					if strings.ToLower(devtype) == "shared" {
						members := stringToSlice(value)
						command = fmt.Sprintf("\nset shared address-group %s static [ %s ]", name, strings.Join(members, " "))

						if desc != "" {
							command += fmt.Sprintf("\nset shared address-group %s description \"%s\"", name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("\nset shared address-group %s tag [ %s ]", name, strings.Join(tags, " "))
						}
					}

					if !strings.Contains(devtype, "vsys") || !strings.Contains(devtype, "hared") {
						members := stringToSlice(value)
						command = fmt.Sprintf("\nset device-group %s address-group %s static [ %s ]", devtype, name, strings.Join(members, " "))

						if desc != "" {
							command += fmt.Sprintf("\nset device-group %s address-group %s description \"%s\"", devtype, name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("\nset device-group %s address-group %s tag [ %s ]", devtype, name, strings.Join(tags, " "))
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
					if strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\ndelete address-group %s", name)
					}

					if strings.ToLower(devtype) == "shared" {
						command = fmt.Sprintf("\ndelete shared address-group %s", name)
					}

					if !strings.Contains(devtype, "vsys") || !strings.Contains(devtype, "hared") {
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

					if strings.ToLower(devtype) == "shared" {
						command = fmt.Sprintf("\nset shared address-group %s dynamic filter \"%s\"", name, value)

						if desc != "" {
							command += fmt.Sprintf("\nset shared address-group %s description \"%s\"", name, desc)
						}

						if len(tg) > 0 {
							tags := stringToSlice(tg)
							command += fmt.Sprintf("\nset shared address-group %s tag [ %s ]", name, strings.Join(tags, " "))
						}
					}

					if !strings.Contains(devtype, "vsys") || !strings.Contains(devtype, "hared") {
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

				// fmt.Printf("%s", command)
				_, err = io.WriteString(txtfile, command)

				if err != nil {
					log.Printf("Failed to write to TXT file - %s", err)
				}
			case "remove-address":
				if len(value) <= 0 {
					log.Printf("Line %d - you must specify a value to remove from group: %s", i+1, name)
				}

				if strings.Contains(devtype, "vsys") {
					command = fmt.Sprintf("\ndelete address-group %s static %s", name, value)
				}

				if strings.ToLower(devtype) == "shared" {
					command = fmt.Sprintf("\ndelete shared address-group %s static %s", name, value)
				}

				if !strings.Contains(devtype, "vsys") || !strings.Contains(devtype, "hared") {
					command = fmt.Sprintf("\ndelete device-group %s address-group %s static %s", devtype, name, value)
				}

				// fmt.Printf("%s", command)
				_, err = io.WriteString(txtfile, command)

				if err != nil {
					log.Printf("Failed to write to TXT file - %s", err)
				}
			case "remove-service":
				if len(value) <= 0 {
					log.Printf("Line %d - you must specify a value to remove from group: %s", i+1, name)
				}

				if strings.Contains(devtype, "vsys") {
					command = fmt.Sprintf("\ndelete service-group %s members %s", name, value)
				}

				if strings.ToLower(devtype) == "shared" {
					command = fmt.Sprintf("\ndelete shared service-group %s members %s", name, value)
				}

				if !strings.Contains(devtype, "vsys") || !strings.Contains(devtype, "hared") {
					command = fmt.Sprintf("\ndelete device-group %s service-group %s members %s", devtype, name, value)
				}

				// fmt.Printf("%s", command)
				_, err = io.WriteString(txtfile, command)

				if err != nil {
					log.Printf("Failed to write to TXT file - %s", err)
				}
			case "rename-address":
				if strings.Contains(devtype, "vsys") {
					command = fmt.Sprintf("\nrename address %s to %s", name, value)
				}

				if strings.ToLower(devtype) == "shared" {
					command = fmt.Sprintf("\nrename shared address %s to %s", name, value)
				}

				if !strings.Contains(devtype, "vsys") || !strings.Contains(devtype, "hared") {
					command = fmt.Sprintf("\nrename device-group %s address %s to %s", devtype, name, value)
				}

				// fmt.Printf("%s", command)
				_, err = io.WriteString(txtfile, command)

				if err != nil {
					log.Printf("Failed to write to TXT file - %s", err)
				}
			case "rename-addressgroup":
				if strings.Contains(devtype, "vsys") {
					command = fmt.Sprintf("\nrename address-group %s to %s", name, value)
				}

				if strings.ToLower(devtype) == "shared" {
					command = fmt.Sprintf("\nrename shared address-group %s to %s", name, value)
				}

				if !strings.Contains(devtype, "vsys") || !strings.Contains(devtype, "hared") {
					command = fmt.Sprintf("\nrename device-group %s address-group %s to %s", devtype, name, value)
				}

				// fmt.Printf("%s", command)
				_, err = io.WriteString(txtfile, command)

				if err != nil {
					log.Printf("Failed to write to TXT file - %s", err)
				}
			case "rename-service":
				if strings.Contains(devtype, "vsys") {
					command = fmt.Sprintf("\nrename service %s to %s", name, value)
				}

				if strings.ToLower(devtype) == "shared" {
					command = fmt.Sprintf("\nrename shared service %s to %s", name, value)
				}

				if !strings.Contains(devtype, "vsys") || !strings.Contains(devtype, "hared") {
					command = fmt.Sprintf("\nrename device-group %s service %s to %s", devtype, name, value)
				}

				// fmt.Printf("%s", command)
				_, err = io.WriteString(txtfile, command)

				if err != nil {
					log.Printf("Failed to write to TXT file - %s", err)
				}
			case "rename-servicegroup":
				if strings.Contains(devtype, "vsys") {
					command = fmt.Sprintf("\nrename service-group %s to %s", name, value)
				}

				if strings.ToLower(devtype) == "shared" {
					command = fmt.Sprintf("\nrename shared service-group %s to %s", name, value)
				}

				if !strings.Contains(devtype, "vsys") || !strings.Contains(devtype, "hared") {
					command = fmt.Sprintf("\nrename device-group %s service-group %s to %s", devtype, name, value)
				}

				// fmt.Printf("%s", command)
				_, err = io.WriteString(txtfile, command)

				if err != nil {
					log.Printf("Failed to write to TXT file - %s", err)
				}
			case "tag":
				if value == "delete" {
					if strings.Contains(devtype, "vsys") {
						command = fmt.Sprintf("\ndelete tag %s", name)
					}

					if strings.ToLower(devtype) == "shared" {
						command = fmt.Sprintf("\ndelete shared tag %s", name)
					}

					if !strings.Contains(devtype, "vsys") || !strings.Contains(devtype, "hared") {
						command = fmt.Sprintf("\ndelete device-group %s tag %s", devtype, name)
					}
				} else {
					if value == "None" || value == "none" || value == "color0" || value == "" {
						if strings.Contains(devtype, "vsys") {
							command = fmt.Sprintf("\nset tag %s color color0", name)

							if desc != "" {
								command += fmt.Sprintf("\nset tag %s comments \"%s\"", name, desc)
							}
						}

						if strings.ToLower(devtype) == "shared" {
							command = fmt.Sprintf("\nset shared tag %s color color0", name)

							if desc != "" {
								command += fmt.Sprintf("\nset shared tag %s comments \"%s\"", name, desc)
							}
						}

						if !strings.Contains(devtype, "vsys") || !strings.Contains(devtype, "hared") {
							command = fmt.Sprintf("\nset device-group %s tag %s color color0", devtype, name)

							if desc != "" {
								command += fmt.Sprintf("\nset device-group %s tag %s comments \"%s\"", devtype, name, desc)
							}
						}
					} else {
						if strings.Contains(devtype, "vsys") {
							command = fmt.Sprintf("\nset tag %s color %s", name, tag2color[value])

							if desc != "" {
								command += fmt.Sprintf("\nset tag %s comments \"%s\"", name, desc)
							}
						}

						if strings.ToLower(devtype) == "shared" {
							command = fmt.Sprintf("\nset shared tag %s color %s", name, tag2color[value])

							if desc != "" {
								command += fmt.Sprintf("\nset shared tag %s comments \"%s\"", name, desc)
							}
						}

						if !strings.Contains(devtype, "vsys") || !strings.Contains(devtype, "hared") {
							command = fmt.Sprintf("\nset device-group %s tag %s color %s", devtype, name, tag2color[value])

							if desc != "" {
								command += fmt.Sprintf("\nset device-group %s tag %s comments \"%s\"", devtype, name, desc)
							}
						}
					}
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
