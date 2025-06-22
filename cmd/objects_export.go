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
	"time"

	"github.com/PaloAltoNetworks/pango"
	"github.com/Songmu/prompter"
	easycsv "github.com/scottdware/go-easycsv"
	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
)

// exportCmd represents the export command
var objectsExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export address, service and tag objects",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var delay time.Duration
		resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
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
			if v == "" {
				v = "vsys1"
			}

			f = strings.TrimSuffix(f, ".csv")

			// Address objects
			if t == "address" {
				getFwAddr(c, fmt.Sprintf("%s.csv", f))
			}

			if t == "all" {
				getFwAddr(c, fmt.Sprintf("%s_Addrs.csv", f))
			}

			// Address groups
			if t == "addressgroup" {
				getFwAddrGrp(c, fmt.Sprintf("%s.csv", f))
			}

			if t == "all" {
				getFwAddrGrp(c, fmt.Sprintf("%s_AddrGroups.csv", f))
			}

			// Service objects
			if t == "service" {
				getFwSrvc(c, fmt.Sprintf("%s.csv", f))
			}

			if t == "all" {
				getFwSrvc(c, fmt.Sprintf("%s_Svcs.csv", f))
			}

			// Service groups
			if t == "servicegroup" {
				getFwSrvcGrp(c, fmt.Sprintf("%s.csv", f))
			}

			if t == "all" {
				getFwSrvcGrp(c, fmt.Sprintf("%s_SvcGroups.csv", f))
			}

			// Tags
			if t == "tags" {
				getFwTags(c, fmt.Sprintf("%s.csv", f))
			}

			if t == "all" {
				getFwTags(c, fmt.Sprintf("%s_Tags.csv", f))
			}
		case *pango.Panorama:
			if dg == "" {
				dg = "shared"
			}

			f = strings.TrimSuffix(f, ".csv")

			// Address objects
			if t == "address" {
				getPanoAddr(c, fmt.Sprintf("%s.csv", f))
			}

			if t == "all" {
				getPanoAddr(c, fmt.Sprintf("%s_Addrs.csv", f))
			}

			// Address groups
			if t == "addressgroup" {
				getPanoAddrGrp(c, fmt.Sprintf("%s.csv", f))
			}

			if t == "all" {
				getPanoAddrGrp(c, fmt.Sprintf("%s_AddrGroups.csv", f))
			}

			// Service objects
			if t == "service" {
				getPanoSrvc(c, fmt.Sprintf("%s.csv", f))
			}

			if t == "all" {
				getPanoSrvc(c, fmt.Sprintf("%s_Svcs.csv", f))
			}

			// Service groups
			if t == "servicegroup" {
				getPanoSrvcGrp(c, fmt.Sprintf("%s.csv", f))
			}

			if t == "all" {
				getPanoSrvcGrp(c, fmt.Sprintf("%s_SvcGroups.csv", f))
			}

			// Tags
			if t == "tags" {
				getPanoTags(c, fmt.Sprintf("%s.csv", f))
			}

			if t == "all" {
				getPanoTags(c, fmt.Sprintf("%s_Tags.csv", f))
			}
		}
	},
}

func init() {
	objectsCmd.AddCommand(objectsExportCmd)

	objectsExportCmd.Flags().StringVarP(&user, "user", "u", "", "User to connect to the device as")
	objectsExportCmd.Flags().StringVarP(&p, "delay", "p", "100", "Delay (in milliseconds) to pause between each API call")
	// objectsExportCmd.Flags().StringVarP(&pass, "pass", "p", "", "Password for the user account specified")
	objectsExportCmd.Flags().StringVarP(&device, "device", "d", "", "Device to connect to")
	objectsExportCmd.Flags().StringVarP(&f, "file", "f", "PaloAltoObjects", "Name of the CSV file you'd like to export to")
	objectsExportCmd.Flags().StringVarP(&dg, "devicegroup", "g", "shared", "Device Group name")
	objectsExportCmd.Flags().StringVarP(&v, "vsys", "v", "vsys1", "Vsys name")
	objectsExportCmd.Flags().StringVarP(&t, "type", "t", "", "<address|addressgroup|service|servicegroup|tags|all>")
	objectsExportCmd.MarkFlagRequired("user")
	// objectsExportCmd.MarkFlagRequired("pass")
	objectsExportCmd.MarkFlagRequired("device")
	objectsExportCmd.MarkFlagRequired("file")
	objectsExportCmd.MarkFlagRequired("type")
}

// getFwAddr is used to export the address objects from a firewall
func getFwAddr(c *pango.Firewall, file string) {
	addrs, err := c.Objects.Address.GetList(v)
	if err != nil {
		log.Printf("Failed to get the list of address objects: %s", err)
		return
	}

	numobjs := len(addrs)
	if numobjs <= 0 {
		log.Printf("No address objects to export\n")
		return
	}

	ac, err := easycsv.NewCSV(file)
	if err != nil {
		log.Printf("CSV file error - %s", err)
		return
	}

	log.Printf("Exporting %d address objects", numobjs)

	ac.Write("#Name,Type,Value,Description,Tags,Device Group/Vsys\n")
	for _, aentry := range addrs {
		a, err := c.Objects.Address.Get(v, aentry)
		if err != nil {
			log.Printf("Failed to retrieve object data for '%s': %s", aentry, err)
		}

		ac.Write(fmt.Sprintf("%s,%s,\"%s\",\"%s\",\"%s\",%s\n", a.Name, a.Type, a.Value, a.Description, sliceToString(a.Tags), v))

		time.Sleep(delay * time.Millisecond)
	}

	ac.End()
}

// getFwAddrGrp is used to export the address group objects from a firewall
func getFwAddrGrp(c *pango.Firewall, file string) {
	addrgrps, err := c.Objects.AddressGroup.GetList(v)
	if err != nil {
		log.Printf("Failed to get the list of address groups: %s", err)
		return
	}

	numobjs := len(addrgrps)
	if numobjs <= 0 {
		log.Printf("No address group objects to export")
		return
	}

	agc, err := easycsv.NewCSV(file)
	if err != nil {
		log.Printf("CSV file error - %s", err)
		return
	}

	log.Printf("Exporting %d address group objects", numobjs)

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

		time.Sleep(delay * time.Millisecond)
	}

	agc.End()
}

// getFwSrvc is used to export the service objects from a firewall
func getFwSrvc(c *pango.Firewall, file string) {
	srvcs, err := c.Objects.Services.GetList(v)
	if err != nil {
		log.Printf("Failed to get the list of service objects: %s", err)
		return
	}

	numobjs := len(srvcs)
	if numobjs <= 0 {
		log.Printf("No service objects to export")
		return
	}

	sc, err := easycsv.NewCSV(file)
	if err != nil {
		log.Printf("CSV file error - %s", err)
		return
	}

	log.Printf("Exporting %d service objects", numobjs)

	sc.Write("#Name,Type,Value,Description,Tags,Device Group/Vsys\n")
	for _, sentry := range srvcs {
		s, err := c.Objects.Services.Get(v, sentry)
		if err != nil {
			log.Printf("Failed to retrieve object data for '%s': %s", sentry, err)
		}

		sc.Write(fmt.Sprintf("%s,%s,\"%s\",\"%s\",\"%s\",%s\n", s.Name, s.Protocol, s.DestinationPort, s.Description, sliceToString(s.Tags), v))

		time.Sleep(delay * time.Millisecond)
	}

	sc.End()
}

// getFwSrvcGrp is used to export the service group objects from a firewall
func getFwSrvcGrp(c *pango.Firewall, file string) {
	srvcgrps, err := c.Objects.ServiceGroup.GetList(v)
	if err != nil {
		log.Printf("Failed to get the list of service groups: %s", err)
		return
	}

	numobjs := len(srvcgrps)
	if numobjs <= 0 {
		log.Printf("No service group objects to export")
		return
	}

	sgc, err := easycsv.NewCSV(file)
	if err != nil {
		log.Printf("CSV file error - %s", err)
		return
	}

	log.Printf("Exporting %d service group objects", numobjs)

	sgc.Write("#Name,Type,Value,Description,Tags,Device Group/Vsys\n")
	for _, sgentry := range srvcgrps {
		sg, err := c.Objects.ServiceGroup.Get(v, sgentry)
		if err != nil {
			log.Printf("Failed to retrieve object data for '%s': %s", sgentry, err)
		}

		sgc.Write(fmt.Sprintf("%s,service,\"%s\",,\"%s\",%s\n", sg.Name, sliceToString(sg.Services), sliceToString(sg.Tags), v))

		time.Sleep(delay * time.Millisecond)
	}

	sgc.End()
}

// getFwTags is used to export the tag objects from a firewall
func getFwTags(c *pango.Firewall, file string) {
	tags, err := c.Objects.Tags.GetList(v)
	if err != nil {
		log.Printf("Failed to get the list of tags: %s", err)
		return
	}

	numobjs := len(tags)
	if numobjs <= 0 {
		log.Printf("No tag objects to export")
		return
	}

	tc, err := easycsv.NewCSV(file)
	if err != nil {
		log.Printf("CSV file error - %s", err)
		return
	}

	log.Printf("Exporting %d tag objects", numobjs)

	tc.Write("#Name,Type,Value,Description,Tags,Device Group/Vsys\n")
	for _, tag := range tags {
		t, err := c.Objects.Tags.Get(v, tag)
		if err != nil {
			log.Printf("Failed to retrieve object data for '%s': %s", tag, err)
		}

		tc.Write(fmt.Sprintf("%s,tag,%s,\"%s\",,%s\n", t.Name, color2tag[t.Color], t.Comment, v))

		time.Sleep(delay * time.Millisecond)
	}

	tc.End()
}

// getPanoAddr is used to export the address objects from Panorama
func getPanoAddr(c *pango.Panorama, file string) {
	addrs, err := c.Objects.Address.GetList(dg)
	if err != nil {
		log.Printf("Failed to get the list of address objects: %s", err)
		return
	}

	numobjs := len(addrs)
	if numobjs <= 0 {
		log.Printf("No address objects to export")
		return
	}

	ac, err := easycsv.NewCSV(file)
	if err != nil {
		log.Printf("CSV file error - %s", err)
		return
	}

	log.Printf("Exporting %d address objects", numobjs)

	ac.Write("#Name,Type,Value,Description,Tags,Device Group/Vsys\n")
	for _, aentry := range addrs {
		a, err := c.Objects.Address.Get(dg, aentry)
		if err != nil {
			log.Printf("Failed to retrieve object data for '%s': %s", aentry, err)
		}

		ac.Write(fmt.Sprintf("%s,%s,\"%s\",\"%s\",\"%s\",%s\n", a.Name, a.Type, a.Value, a.Description, sliceToString(a.Tags), dg))

		time.Sleep(delay * time.Millisecond)
	}

	ac.End()
}

// getPanoAddrGrp is used to export the address group objects from Panorama
func getPanoAddrGrp(c *pango.Panorama, file string) {
	addrgrps, err := c.Objects.AddressGroup.GetList(dg)
	if err != nil {
		log.Printf("Failed to get the list of address groups: %s", err)
		return
	}

	numobjs := len(addrgrps)
	if numobjs <= 0 {
		log.Printf("No address group objects to export")
		return
	}

	agc, err := easycsv.NewCSV(file)
	if err != nil {
		log.Printf("CSV file error - %s", err)
		return
	}

	log.Printf("Exporting %d address group objects", numobjs)

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

		time.Sleep(delay * time.Millisecond)
	}

	agc.End()
}

// getPanoSrvc is used to export the service objects from Panorama
func getPanoSrvc(c *pango.Panorama, file string) {
	srvcs, err := c.Objects.Services.GetList(dg)
	if err != nil {
		log.Printf("Failed to get the list of service objects: %s", err)
		return
	}

	numobjs := len(srvcs)
	if numobjs <= 0 {
		log.Printf("No service objects to export")
		return
	}

	sc, err := easycsv.NewCSV(file)
	if err != nil {
		log.Printf("CSV file error - %s", err)
		return
	}

	log.Printf("Exporting %d service objects", numobjs)

	sc.Write("#Name,Type,Value,Description,Tags,Device Group/Vsys\n")
	for _, sentry := range srvcs {
		s, err := c.Objects.Services.Get(dg, sentry)
		if err != nil {
			log.Printf("Failed to retrieve object data for '%s': %s", sentry, err)
		}

		sc.Write(fmt.Sprintf("%s,%s,\"%s\",\"%s\",\"%s\",%s\n", s.Name, s.Protocol, s.DestinationPort, s.Description, sliceToString(s.Tags), dg))

		time.Sleep(delay * time.Millisecond)
	}

	sc.End()
}

// getPanoSrvGrp is used to export the service group objects from Panorama
func getPanoSrvcGrp(c *pango.Panorama, file string) {
	srvcgrps, err := c.Objects.ServiceGroup.GetList(dg)
	if err != nil {
		log.Printf("Failed to get the list of service groups: %s", err)
		return
	}

	numobjs := len(srvcgrps)
	if numobjs <= 0 {
		log.Printf("No service group objects to export")
		return
	}

	sgc, err := easycsv.NewCSV(file)
	if err != nil {
		log.Printf("CSV file error - %s", err)
		return
	}

	log.Printf("Exporting %d service group objects", numobjs)

	sgc.Write("#Name,Type,Value,Description,Tags,Device Group/Vsys\n")
	for _, sgentry := range srvcgrps {
		sg, err := c.Objects.ServiceGroup.Get(dg, sgentry)
		if err != nil {
			log.Printf("Failed to retrieve object data for '%s': %s", sgentry, err)
		}

		sgc.Write(fmt.Sprintf("%s,service,\"%s\",,\"%s\",%s\n", sg.Name, sliceToString(sg.Services), sliceToString(sg.Tags), dg))

		time.Sleep(delay * time.Millisecond)
	}

	sgc.End()
}

// getPanoTags is used to export the tag objects from Panorama
func getPanoTags(c *pango.Panorama, file string) {
	tags, err := c.Objects.Tags.GetList(dg)
	if err != nil {
		log.Printf("Failed to get the list of tags: %s", err)
		return
	}

	numobjs := len(tags)
	if numobjs <= 0 {
		log.Printf("No tag objects to export")
		return
	}

	tc, err := easycsv.NewCSV(file)
	if err != nil {
		log.Printf("CSV file error - %s", err)
		return
	}

	log.Printf("Exporting %d tag objects", numobjs)

	tc.Write("#Name,Type,Value,Description,Tags,Device Group/Vsys\n")
	for _, tag := range tags {
		t, err := c.Objects.Tags.Get(dg, tag)
		if err != nil {
			log.Printf("Failed to retrieve object data for '%s': %s", tag, err)
		}

		tc.Write(fmt.Sprintf("%s,tag,%s,\"%s\",,%s\n", t.Name, color2tag[t.Color], t.Comment, dg))

		time.Sleep(delay * time.Millisecond)
	}

	tc.End()
}
