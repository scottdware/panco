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
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/PaloAltoNetworks/pango"
	"github.com/Songmu/prompter"
	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
)

// Local variables
var dupfind, dupremove string

// duplicatesCmd represents the duplicates command
var objectsDuplicatesCmd = &cobra.Command{
	Use:   "duplicates",
	Short: "Find duplicate address and service objects",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
		passwd := prompter.Password(fmt.Sprintf("Password for %s", user))
		_ = passwd

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
			if v == "" {
				v = "vsys1"
			}

			xlsx := excelize.NewFile()
			f = strings.TrimSuffix(f, ".csv")
			f = strings.TrimSuffix(f, "xlsx")

			// Address objects
			if t == "address" || t == "all" {
				xlsx.NewSheet("address-Unique")
				xlsx.NewSheet("address-Duplicates")
				xlsx.SetCellValue("address-Unique", "A1", "Object Name")
				xlsx.SetCellValue("address-Unique", "B1", "Object Value")
				xlsx.SetCellValue("address-Unique", "C1", "Device Group/Vsys")
				xlsx.SetCellValue("address-Unique", "D1", "Total Objects")
				xlsx.SetCellValue("address-Duplicates", "A1", "Object Name")
				xlsx.SetCellValue("address-Duplicates", "B1", "Object Value")
				xlsx.SetCellValue("address-Duplicates", "C1", "Device Group/Vsys")
				xlsx.SetCellValue("address-Duplicates", "D1", "Total Objects")
				xlsx.SetColWidth("address-Unique", "A", "C", 20)
				xlsx.SetColWidth("address-Unique", "D", "D", 10)
				xlsx.SetColWidth("address-Duplicates", "A", "C", 20)
				xlsx.SetColWidth("address-Duplicates", "D", "D", 10)

				addrs, err := c.Objects.Address.GetList(v)
				if err != nil {
					log.Printf("Failed to get the list of address objects: %s", err)
				}

				numobjs := len(addrs)
				log.Printf("Finding duplicates on %d address objects", numobjs)

				objs := map[string]string{}

				for _, aentry := range addrs {
					a, err := c.Objects.Address.Get(v, aentry)
					if err != nil {
						log.Printf("Failed to retrieve object data for '%s': %s", aentry, err)
					}

					objs[a.Name] = a.Value
				}

				unique, duplicates := duplicateObjects(objs)

				xlsx.SetCellValue("address-Unique", "D2", len(unique))
				xlsx.SetCellValue("address-Duplicates", "D2", len(duplicates))

				for i := 1; i <= len(unique); i++ {
					for key, val := range unique {
						xlsx.SetCellValue("address-Unique", fmt.Sprintf("A%d", i+1), key)
						xlsx.SetCellValue("address-Unique", fmt.Sprintf("B%d", i+1), val)
						xlsx.SetCellValue("address-Unique", fmt.Sprintf("C%d", i+1), v)
						i++
					}
				}

				for i := 1; i <= len(duplicates); i++ {
					for key, val := range duplicates {
						xlsx.SetCellValue("address-Duplicates", fmt.Sprintf("A%d", i+1), key)
						xlsx.SetCellValue("address-Duplicates", fmt.Sprintf("B%d", i+1), val)
						xlsx.SetCellValue("address-Duplicates", fmt.Sprintf("C%d", i+1), v)
						i++
					}
				}
			}

			// Address groups
			// if t == "addressgroup" || t == "all" {
			// 	xlsx.NewSheet("addrgrp-Unique")
			// 	xlsx.NewSheet("addrgrp-Duplicates")
			// 	xlsx.SetCellValue("addrgrp-Unique", "A1", "Object Name")
			// 	xlsx.SetCellValue("addrgrp-Unique", "B1", "Object Value")
			// 	xlsx.SetCellValue("addrgrp-Unique", "C1", "Device Group/Vsys")
			// 	xlsx.SetCellValue("addrgrp-Unique", "D1", "Total Objects")
			// 	xlsx.SetCellValue("addrgrp-Duplicates", "A1", "Object Name")
			// 	xlsx.SetCellValue("addrgrp-Duplicates", "B1", "Object Value")
			// 	xlsx.SetCellValue("addrgrp-Duplicates", "C1", "Device Group/Vsys")
			// 	xlsx.SetCellValue("addrgrp-Duplicates", "D1", "Total Objects")
			// 	xlsx.SetColWidth("addrgrp-Unique", "A", "C", 20)
			// 	xlsx.SetColWidth("addrgrp-Unique", "D", "D", 10)
			// 	xlsx.SetColWidth("addrgrp-Duplicates", "A", "C", 20)
			// 	xlsx.SetColWidth("addrgrp-Duplicates", "D", "D", 10)

			// 	addrgrps, err := c.Objects.AddressGroup.GetList(v)
			// 	if err != nil {
			// 		log.Printf("Failed to get the list of address groups: %s", err)
			// 	}

			// 	numobjs := len(addrgrps)
			// 	log.Printf("Finding duplicates on %d address group objects - this might take a few minutes if you have a lot", numobjs)

			// 	objs := map[string]string{}

			// 	for _, agentry := range addrgrps {
			// 		var val string
			// 		a, err := c.Objects.AddressGroup.Get(v, agentry)
			// 		if err != nil {
			// 			log.Printf("Failed to retrieve object data for '%s': %s", agentry, err)
			// 		}

			// 		if len(a.StaticAddresses) <= 0 && len(a.DynamicMatch) > 0 {
			// 			// gtype = "dynamic"
			// 			val = a.DynamicMatch
			// 		}

			// 		if len(a.DynamicMatch) <= 0 && len(a.StaticAddresses) > 0 {
			// 			// gtype = "static"
			// 			val = sliceToString(a.StaticAddresses)
			// 		}

			// 		objs[a.Name] = val
			// 	}

			// 	unique, duplicates := duplicateObjects(objs)

			// 	xlsx.SetCellValue("addrgrp-Unique", "D2", len(unique))
			// 	xlsx.SetCellValue("addrgrp-Duplicates", "D2", len(duplicates))

			// 	for i := 2; i < len(unique); i++ {
			// 		for key, val := range unique {
			// 			xlsx.SetCellValue("addrgrp-Unique", fmt.Sprintf("A%d", i), k)
			// 			xlsx.SetCellValue("addrgrp-Unique", fmt.Sprintf("B%d", i), v)
			// 			i++
			// 		}
			// 	}

			// 	for i := 2; i < len(duplicates); i++ {
			// 		for key, val := range unique {
			// 			xlsx.SetCellValue("addrgrp-Duplicates", fmt.Sprintf("A%d", i), k)
			// 			xlsx.SetCellValue("addrgrp-Duplicates", fmt.Sprintf("B%d", i), v)
			// 			i++
			// 		}
			// 	}
			// }

			// Service objects
			if t == "service" || t == "all" {
				xlsx.NewSheet("service-Unique")
				xlsx.NewSheet("service-Duplicates")
				xlsx.SetCellValue("service-Unique", "A1", "Object Name")
				xlsx.SetCellValue("service-Unique", "B1", "Object Value")
				xlsx.SetCellValue("service-Unique", "C1", "Device Group/Vsys")
				xlsx.SetCellValue("service-Unique", "D1", "Total Objects")
				xlsx.SetCellValue("service-Duplicates", "A1", "Object Name")
				xlsx.SetCellValue("service-Duplicates", "B1", "Object Value")
				xlsx.SetCellValue("service-Duplicates", "C1", "Device Group/Vsys")
				xlsx.SetCellValue("service-Duplicates", "D1", "Total Objects")
				xlsx.SetColWidth("service-Unique", "A", "C", 20)
				xlsx.SetColWidth("service-Unique", "D", "D", 10)
				xlsx.SetColWidth("service-Duplicates", "A", "C", 20)
				xlsx.SetColWidth("service-Duplicates", "D", "D", 10)

				srvcs, err := c.Objects.Services.GetList(v)
				if err != nil {
					log.Printf("Failed to get the list of service objects: %s", err)
				}

				numobjs := len(srvcs)
				log.Printf("Finding duplicates on %d service objects", numobjs)

				objs := map[string]string{}

				for _, sentry := range srvcs {
					s, err := c.Objects.Services.Get(v, sentry)
					if err != nil {
						log.Printf("Failed to retrieve object data for '%s': %s", sentry, err)
					}

					objs[s.Name] = s.DestinationPort
				}

				unique, duplicates := duplicateObjects(objs)

				xlsx.SetCellValue("service-Unique", "D2", len(unique))
				xlsx.SetCellValue("service-Duplicates", "D2", len(duplicates))

				for i := 1; i <= len(unique); i++ {
					for key, val := range unique {
						xlsx.SetCellValue("service-Unique", fmt.Sprintf("A%d", i+1), key)
						xlsx.SetCellValue("service-Unique", fmt.Sprintf("B%d", i+1), val)
						xlsx.SetCellValue("service-Unique", fmt.Sprintf("C%d", i+1), v)
						i++
					}
				}

				for i := 1; i <= len(duplicates); i++ {
					for key, val := range duplicates {
						xlsx.SetCellValue("service-Duplicates", fmt.Sprintf("A%d", i+1), key)
						xlsx.SetCellValue("service-Duplicates", fmt.Sprintf("B%d", i+1), val)
						xlsx.SetCellValue("service-Duplicates", fmt.Sprintf("C%d", i+1), v)
						i++
					}
				}
			}

			// Service groups
			// if t == "servicegroup" || t == "all" {
			// 	xlsx.NewSheet("srvcgrp-Unique")
			// 	xlsx.NewSheet("srvcgrp-Duplicates")
			// 	xlsx.SetCellValue("srvcgrp-Unique", "A1", "Object Name")
			// 	xlsx.SetCellValue("srvcgrp-Unique", "B1", "Object Value")
			// 	xlsx.SetCellValue("srvcgrp-Unique", "C1", "Device Group/Vsys")
			// 	xlsx.SetCellValue("srvcgrp-Unique", "D1", "Total Objects")
			// 	xlsx.SetCellValue("srvcgrp-Duplicates", "A1", "Object Name")
			// 	xlsx.SetCellValue("srvcgrp-Duplicates", "B1", "Object Value")
			// 	xlsx.SetCellValue("srvcgrp-Duplicates", "C1", "Device Group/Vsys")
			// 	xlsx.SetCellValue("srvcgrp-Duplicates", "D1", "Total Objects")
			// 	xlsx.SetColWidth("srvcgrp-Unique", "A", "C", 20)
			// 	xlsx.SetColWidth("srvcgrp-Unique", "D", "D", 10)
			// 	xlsx.SetColWidth("srvcgrp-Duplicates", "A", "C", 20)
			// 	xlsx.SetColWidth("srvcgrp-Duplicates", "D", "D", 10)

			// 	srvcgrps, err := c.Objects.ServiceGroup.GetList(v)
			// 	if err != nil {
			// 		log.Printf("Failed to get the list of service groups: %s", err)
			// 	}

			// 	numobjs := len(srvcgrps)
			// 	log.Printf("Finding duplicates on %d service group objects - this might take a few minutes if you have a lot", numobjs)

			// 	objs := map[string]string{}

			// 	for _, sgentry := range srvcgrps {
			// 		sg, err := c.Objects.ServiceGroup.Get(v, sgentry)
			// 		if err != nil {
			// 			log.Printf("Failed to retrieve object data for '%s': %s", sgentry, err)
			// 		}

			// 		objs[sg.Name] = sliceToString(sg.Services)
			// 	}

			// 	unique, duplicates := duplicateObjects(objs)

			// 	xlsx.SetCellValue("srvcgrp-Unique", "D2", len(unique))
			// 	xlsx.SetCellValue("srvcgrp-Duplicates", "D2", len(duplicates))

			// 	for i := 2; i < len(unique); i++ {
			// 		for key, val := range unique {
			// 			xlsx.SetCellValue("srvcgrp-Unique", fmt.Sprintf("A%d", i), k)
			// 			xlsx.SetCellValue("srvcgrp-Unique", fmt.Sprintf("B%d", i), v)
			// 			i++
			// 		}
			// 	}

			// 	for i := 2; i < len(duplicates); i++ {
			// 		for key, val := range unique {
			// 			xlsx.SetCellValue("srvcgrp-Duplicates", fmt.Sprintf("A%d", i), k)
			// 			xlsx.SetCellValue("srvcgrp-Duplicates", fmt.Sprintf("B%d", i), v)
			// 			i++
			// 		}
			// 	}
			// }

			xlsx.DeleteSheet("Sheet1")
			if err := xlsx.SaveAs(fmt.Sprintf("%s.xlsx", f)); err != nil {
				fmt.Println(err.Error())
			}

			// Tags
			// if t == "tags" || t == "all" {
			// 	tc, err := easycsv.NewCSV(tagFile)
			// 	if err != nil {
			// 		log.Printf("CSV file error - %s", err)
			// 		os.Exit(1)
			// 	}

			// 	tags, err := c.Objects.Tags.GetList(v)
			// 	if err != nil {
			// 		log.Printf("Failed to get the list of tags: %s", err)
			// 		os.Remove(tagFile)
			// 	}

			// 	numobjs := len(tags)
			// 	log.Printf("Exporting %d tag objects - this might take a few minutes if you have a lot", numobjs)

			// 	tc.Write("#Name,Type,Value,Description,Tags,Device Group/Vsys\n")
			// 	for _, tag := range tags {
			// 		t, err := c.Objects.Tags.Get(v, tag)
			// 		if err != nil {
			// 			log.Printf("Failed to retrieve object data for '%s': %s", tag, err)
			// 		}

			// 		tc.Write(fmt.Sprintf("%s,tag,%s,\"%s\",,%s\n", t.Name, color2tag[t.Color], t.Comment, v))
			// 	}

			// 	tc.End()
			// }
		case *pango.Panorama:
			if dg == "" {
				dg = "shared"
			}

			xlsx := excelize.NewFile()
			f = strings.TrimSuffix(f, ".csv")
			f = strings.TrimSuffix(f, "xlsx")

			// Address objects
			if t == "address" || t == "all" {
				xlsx.NewSheet("address-Unique")
				xlsx.NewSheet("address-Duplicates")
				xlsx.SetCellValue("address-Unique", "A1", "Object Name")
				xlsx.SetCellValue("address-Unique", "B1", "Object Value")
				xlsx.SetCellValue("address-Unique", "C1", "Device Group/Vsys")
				xlsx.SetCellValue("address-Unique", "D1", "Total Objects")
				xlsx.SetCellValue("address-Duplicates", "A1", "Object Name")
				xlsx.SetCellValue("address-Duplicates", "B1", "Object Value")
				xlsx.SetCellValue("address-Duplicates", "C1", "Device Group/Vsys")
				xlsx.SetCellValue("address-Duplicates", "D1", "Total Objects")
				xlsx.SetColWidth("address-Unique", "A", "C", 20)
				xlsx.SetColWidth("address-Unique", "D", "D", 10)
				xlsx.SetColWidth("address-Duplicates", "A", "C", 20)
				xlsx.SetColWidth("address-Duplicates", "D", "D", 10)

				addrs, err := c.Objects.Address.GetList(dg)
				if err != nil {
					log.Printf("Failed to get the list of address objects: %s", err)
				}

				numobjs := len(addrs)
				log.Printf("Finding duplicates on %d address objects", numobjs)

				objs := map[string]string{}

				for _, aentry := range addrs {
					a, err := c.Objects.Address.Get(dg, aentry)
					if err != nil {
						log.Printf("Failed to retrieve object data for '%s': %s", aentry, err)
					}

					objs[a.Name] = a.Value
				}

				unique, duplicates := duplicateObjects(objs)

				xlsx.SetCellValue("address-Unique", "D2", len(unique))
				xlsx.SetCellValue("address-Duplicates", "D2", len(duplicates))

				for i := 1; i <= len(unique); i++ {
					for key, val := range unique {
						xlsx.SetCellValue("address-Unique", fmt.Sprintf("A%d", i+1), key)
						xlsx.SetCellValue("address-Unique", fmt.Sprintf("B%d", i+1), val)
						xlsx.SetCellValue("address-Unique", fmt.Sprintf("C%d", i+1), dg)
						i++
					}
				}

				for i := 1; i <= len(duplicates); i++ {
					for key, val := range duplicates {
						xlsx.SetCellValue("address-Duplicates", fmt.Sprintf("A%d", i+1), key)
						xlsx.SetCellValue("address-Duplicates", fmt.Sprintf("B%d", i+1), val)
						xlsx.SetCellValue("address-Duplicates", fmt.Sprintf("C%d", i+1), dg)
						i++
					}
				}
			}

			// Address groups
			// if t == "addressgroup" || t == "all" {
			// 	xlsx.NewSheet("addrgrp-Unique")
			// 	xlsx.NewSheet("addrgrp-Duplicates")
			// 	xlsx.SetCellValue("addrgrp-Unique", "A1", "Object Name")
			// 	xlsx.SetCellValue("addrgrp-Unique", "B1", "Object Value")
			// 	xlsx.SetCellValue("addrgrp-Unique", "C1", "Device Group/Vsys")
			// 	xlsx.SetCellValue("addrgrp-Unique", "D1", "Total Objects")
			// 	xlsx.SetCellValue("addrgrp-Duplicates", "A1", "Object Name")
			// 	xlsx.SetCellValue("addrgrp-Duplicates", "B1", "Object Value")
			// 	xlsx.SetCellValue("addrgrp-Duplicates", "C1", "Device Group/Vsys")
			// 	xlsx.SetCellValue("addrgrp-Duplicates", "D1", "Total Objects")
			// 	xlsx.SetColWidth("addrgrp-Unique", "A", "C", 20)
			// 	xlsx.SetColWidth("addrgrp-Unique", "D", "D", 10)
			// 	xlsx.SetColWidth("addrgrp-Duplicates", "A", "C", 20)
			// 	xlsx.SetColWidth("addrgrp-Duplicates", "D", "D", 10)

			// 	addrgrps, err := c.Objects.AddressGroup.GetList(dg)
			// 	if err != nil {
			// 		log.Printf("Failed to get the list of address groups: %s", err)
			// 	}

			// 	numobjs := len(addrgrps)
			// 	log.Printf("Finding duplicates on %d address group objects - this might take a few minutes if you have a lot", numobjs)

			// 	objs := map[string]string{}

			// 	for _, agentry := range addrgrps {
			// 		var val string
			// 		a, err := c.Objects.AddressGroup.Get(dg, agentry)
			// 		if err != nil {
			// 			log.Printf("Failed to retrieve object data for '%s': %s", agentry, err)
			// 		}

			// 		if len(a.StaticAddresses) <= 0 && len(a.DynamicMatch) > 0 {
			// 			// gtype = "dynamic"
			// 			val = a.DynamicMatch
			// 		}

			// 		if len(a.DynamicMatch) <= 0 && len(a.StaticAddresses) > 0 {
			// 			// gtype = "static"
			// 			val = sliceToString(a.StaticAddresses)
			// 		}

			// 		objs[a.Name] = val
			// 	}

			// 	unique, duplicates := duplicateObjects(objs)

			// 	xlsx.SetCellValue("addrgrp-Unique", "D2", len(unique))
			// 	xlsx.SetCellValue("addrgrp-Duplicates", "D2", len(duplicates))

			// 	for i := 2; i < len(unique); i++ {
			// 		for key, val := range unique {
			// 			xlsx.SetCellValue("addrgrp-Unique", fmt.Sprintf("A%d", i), k)
			// 			xlsx.SetCellValue("addrgrp-Unique", fmt.Sprintf("B%d", i), v)
			// 			i++
			// 		}
			// 	}

			// 	for i := 2; i < len(duplicates); i++ {
			// 		for key, val := range unique {
			// 			xlsx.SetCellValue("addrgrp-Duplicates", fmt.Sprintf("A%d", i), k)
			// 			xlsx.SetCellValue("addrgrp-Duplicates", fmt.Sprintf("B%d", i), v)
			// 			i++
			// 		}
			// 	}
			// }

			// Service objects
			if t == "service" || t == "all" {
				xlsx.NewSheet("service-Unique")
				xlsx.NewSheet("service-Duplicates")
				xlsx.SetCellValue("service-Unique", "A1", "Object Name")
				xlsx.SetCellValue("service-Unique", "B1", "Object Value")
				xlsx.SetCellValue("service-Unique", "C1", "Device Group/Vsys")
				xlsx.SetCellValue("service-Unique", "D1", "Total Objects")
				xlsx.SetCellValue("service-Duplicates", "A1", "Object Name")
				xlsx.SetCellValue("service-Duplicates", "B1", "Object Value")
				xlsx.SetCellValue("service-Duplicates", "C1", "Device Group/Vsys")
				xlsx.SetCellValue("service-Duplicates", "D1", "Total Objects")
				xlsx.SetColWidth("service-Unique", "A", "C", 20)
				xlsx.SetColWidth("service-Unique", "D", "D", 10)
				xlsx.SetColWidth("service-Duplicates", "A", "C", 20)
				xlsx.SetColWidth("service-Duplicates", "D", "D", 10)

				srvcs, err := c.Objects.Services.GetList(dg)
				if err != nil {
					log.Printf("Failed to get the list of service objects: %s", err)
				}

				numobjs := len(srvcs)
				log.Printf("Finding duplicates on %d service objects - this might take a few minutes if you have a lot", numobjs)

				objs := map[string]string{}

				for _, sentry := range srvcs {
					s, err := c.Objects.Services.Get(dg, sentry)
					if err != nil {
						log.Printf("Failed to retrieve object data for '%s': %s", sentry, err)
					}

					objs[s.Name] = s.DestinationPort
				}

				unique, duplicates := duplicateObjects(objs)

				xlsx.SetCellValue("service-Unique", "D2", len(unique))
				xlsx.SetCellValue("service-Duplicates", "D2", len(duplicates))

				for i := 1; i <= len(unique); i++ {
					for key, val := range unique {
						xlsx.SetCellValue("service-Unique", fmt.Sprintf("A%d", i+1), key)
						xlsx.SetCellValue("service-Unique", fmt.Sprintf("B%d", i+1), val)
						xlsx.SetCellValue("service-Unique", fmt.Sprintf("C%d", i+1), dg)
						i++
					}
				}

				for i := 1; i <= len(duplicates); i++ {
					for key, val := range duplicates {
						xlsx.SetCellValue("service-Duplicates", fmt.Sprintf("A%d", i+1), key)
						xlsx.SetCellValue("service-Duplicates", fmt.Sprintf("B%d", i+1), val)
						xlsx.SetCellValue("service-Duplicates", fmt.Sprintf("C%d", i+1), dg)
						i++
					}
				}
			}

			// Service groups
			// if t == "servicegroup" || t == "all" {
			// 	xlsx.NewSheet("srvcgrp-Unique")
			// 	xlsx.NewSheet("srvcgrp-Duplicates")
			// 	xlsx.SetCellValue("srvcgrp-Unique", "A1", "Object Name")
			// 	xlsx.SetCellValue("srvcgrp-Unique", "B1", "Object Value")
			// 	xlsx.SetCellValue("srvcgrp-Unique", "C1", "Device Group/Vsys")
			// 	xlsx.SetCellValue("srvcgrp-Unique", "D1", "Total Objects")
			// 	xlsx.SetCellValue("srvcgrp-Duplicates", "A1", "Object Name")
			// 	xlsx.SetCellValue("srvcgrp-Duplicates", "B1", "Object Value")
			// 	xlsx.SetCellValue("srvcgrp-Duplicates", "C1", "Device Group/Vsys")
			// 	xlsx.SetCellValue("srvcgrp-Duplicates", "D1", "Total Objects")
			// 	xlsx.SetColWidth("srvcgrp-Unique", "A", "C", 20)
			// 	xlsx.SetColWidth("srvcgrp-Unique", "D", "D", 10)
			// 	xlsx.SetColWidth("srvcgrp-Duplicates", "A", "C", 20)
			// 	xlsx.SetColWidth("srvcgrp-Duplicates", "D", "D", 10)

			// 	srvcgrps, err := c.Objects.ServiceGroup.GetList(dg)
			// 	if err != nil {
			// 		log.Printf("Failed to get the list of service groups: %s", err)
			// 	}

			// 	numobjs := len(srvcgrps)
			// 	log.Printf("Finding duplicates on %d service group objects - this might take a few minutes if you have a lot", numobjs)

			// 	objs := map[string]string{}

			// 	for _, sgentry := range srvcgrps {
			// 		sg, err := c.Objects.ServiceGroup.Get(dg, sgentry)
			// 		if err != nil {
			// 			log.Printf("Failed to retrieve object data for '%s': %s", sgentry, err)
			// 		}

			// 		objs[sg.Name] = sliceToString(sg.Services)
			// 	}

			// 	unique, duplicates := duplicateObjects(objs)

			// 	xlsx.SetCellValue("srvcgrp-Unique", "D2", len(unique))
			// 	xlsx.SetCellValue("srvcgrp-Duplicates", "D2", len(duplicates))

			// 	for i := 2; i < len(unique); i++ {
			// 		for key, val := range unique {
			// 			xlsx.SetCellValue("srvcgrp-Unique", fmt.Sprintf("A%d", i), k)
			// 			xlsx.SetCellValue("srvcgrp-Unique", fmt.Sprintf("B%d", i), v)
			// 			i++
			// 		}
			// 	}

			// 	for i := 2; i < len(duplicates); i++ {
			// 		for key, val := range unique {
			// 			xlsx.SetCellValue("srvcgrp-Duplicates", fmt.Sprintf("A%d", i), k)
			// 			xlsx.SetCellValue("srvcgrp-Duplicates", fmt.Sprintf("B%d", i), v)
			// 			i++
			// 		}
			// 	}
			// }

			xlsx.DeleteSheet("Sheet1")
			if err := xlsx.SaveAs(fmt.Sprintf("%s.xlsx", f)); err != nil {
				fmt.Println(err.Error())
			}

			// Tags
			// if t == "tags" || t == "all" {
			// 	tc, err := easycsv.NewCSV(tagFile)
			// 	if err != nil {
			// 		log.Printf("CSV file error - %s", err)
			// 		os.Exit(1)
			// 	}

			// 	tags, err := c.Objects.Tags.GetList(v)
			// 	if err != nil {
			// 		log.Printf("Failed to get the list of tags: %s", err)
			// 		os.Remove(tagFile)
			// 	}

			// 	numobjs := len(tags)
			// 	log.Printf("Exporting %d tag objects - this might take a few minutes if you have a lot", numobjs)

			// 	tc.Write("#Name,Type,Value,Description,Tags,Device Group/Vsys\n")
			// 	for _, tag := range tags {
			// 		t, err := c.Objects.Tags.Get(v, tag)
			// 		if err != nil {
			// 			log.Printf("Failed to retrieve object data for '%s': %s", tag, err)
			// 		}

			// 		tc.Write(fmt.Sprintf("%s,tag,%s,\"%s\",,%s\n", t.Name, color2tag[t.Color], t.Comment, v))
			// 	}

			// 	tc.End()
			// }
		}
	},
}

func init() {
	// objectsCmd.AddCommand(objectsDuplicatesCmd)

	// objectsDuplicatesCmd.Flags().StringVarP(&dupfind, "find", "s", "", "Object type to run duplicate actions against")
	// objectsDuplicatesCmd.Flags().StringVarP(&dupremove, "remove", "r", "", "Object type to remove duplicates for")
	objectsDuplicatesCmd.Flags().StringVarP(&user, "user", "u", "", "User to connect to the device as")
	// objectsDuplicatesCmd.Flags().StringVarP(&pass, "pass", "p", "", "Password for the user account specified")
	objectsDuplicatesCmd.Flags().StringVarP(&device, "device", "d", "", "Device to connect to")
	objectsDuplicatesCmd.Flags().StringVarP(&f, "file", "f", "PaloAltoDuplicates", "Name of the output file (you don't need an extension)")
	objectsDuplicatesCmd.Flags().StringVarP(&dg, "devicegroup", "g", "shared", "Device Group name when exporting from Panorama")
	objectsDuplicatesCmd.Flags().StringVarP(&v, "vsys", "v", "vsys1", "Vsys name when exporting from a firewall")
	objectsDuplicatesCmd.Flags().StringVarP(&t, "type", "t", "", "<address|service|all>")
	objectsDuplicatesCmd.MarkFlagRequired("user")
	// objectsDuplicatesCmd.MarkFlagRequired("pass")
	objectsDuplicatesCmd.MarkFlagRequired("device")
	objectsDuplicatesCmd.MarkFlagRequired("file")
}
