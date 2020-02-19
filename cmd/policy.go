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
	"strconv"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/poli/nat"
	"github.com/PaloAltoNetworks/pango/poli/security"
	"github.com/PaloAltoNetworks/pango/util"
	easycsv "github.com/scottdware/go-easycsv"
	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
)

// policyCmd represents the policy command
var policyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Import/export a security policy, move rules",
	Long: `This command will allow you to import and export an entire security policy, along
with moving rules within the policy. When importing, this allows you to create new rules, 
or modify existing values in rules.

When moving rules, if you are only doing one at a time, you do not need to specify a CSV file
or the '--movemultiple' flag. However, if you are wanting to move multiple rules around, then
you will want to use a CSV file, and it must include the '--movemultiple' flag.

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
			if action == "export" && !xlate {
				if !strings.Contains(fh, ".csv") {
					fh = fmt.Sprintf("%s.csv", fh)
				}

				cfh, err := easycsv.NewCSV(fh)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				rules, err := c.Policies.Security.GetList(v)
				if err != nil {
					log.Printf("Failed to retrieve the list of security rules: %s", err)
					os.Remove(fh)
					os.Exit(1)
				}

				rc := len(rules)
				if rc <= 0 {
					log.Printf("There are 0 security rules for '%s' - no policy was exported.", v)
					os.Remove(fh)
					os.Exit(1)
				}

				log.Printf("Exporting %d security rules - this might take a few of minutes if your rule base is large", rc)

				cfh.Write("#Name,Type,Description,Tags,SourceZones,SourceAddresses,NegateSource,SourceUsers,HipProfiles,")
				cfh.Write("DestinationZones,DestinationAddresses,NegateDestination,Applications,Services,Categories,")
				cfh.Write("Action,LogSetting,LogStart,LogEnd,Disabled,Schedule,IcmpUnreachable,DisableServerResponseInspection,")
				cfh.Write("Group,Virus,Spyware,Vulnerability,UrlFiltering,FileBlocking,WildFireAnalysis,DataFiltering\n")
				for _, rule := range rules {
					var rtype string
					r, err := c.Policies.Security.Get(v, rule)
					if err != nil {
						log.Printf("Failed to retrieve security rule data: %s", err)
					}

					switch r.Type {
					case "universal":
						rtype = "universal"
					case "intrazone":
						rtype = "intrazone"
					case "interzone":
						rtype = "interzone"
					default:
						rtype = "universal"
					}

					cfh.Write(fmt.Sprintf("%s,%s,\"%s\",\"%s\",\"%s\",\"%s\",%t,\"%s\",\"%s\",", r.Name, rtype, r.Description, sliceToString(r.Tags), sliceToString(r.SourceZones),
						sliceToString(r.SourceAddresses), r.NegateSource, userSliceToString(r.SourceUsers), sliceToString(r.HipProfiles)))
					cfh.Write(fmt.Sprintf("\"%s\",\"%s\",%t,\"%s\",\"%s\",\"%s\",", sliceToString(r.DestinationZones), sliceToString(r.DestinationAddresses), r.NegateDestination,
						sliceToString(r.Applications), sliceToString(r.Services), sliceToString(r.Categories)))
					cfh.Write(fmt.Sprintf("%s,%s,%t,%t,%t,%s,%t,%t,", r.Action, r.LogSetting, r.LogStart, r.LogEnd, r.Disabled, r.Schedule,
						r.IcmpUnreachable, r.DisableServerResponseInspection))
					cfh.Write(fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s\n", r.Group, r.Virus, r.Spyware,
						r.Vulnerability, r.UrlFiltering, r.FileBlocking, r.WildFireAnalysis, r.DataFiltering))
				}

				cfh.End()
			}

			if action == "export" && xlate {
				if !strings.Contains(fh, ".csv") {
					fh = fmt.Sprintf("%s.csv", fh)
				}

				cfh, err := easycsv.NewCSV(fh)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				rules, err := c.Policies.Nat.GetList(v)
				if err != nil {
					log.Printf("Failed to retrieve the list of NAT rules: %s", err)
					os.Remove(fh)
					os.Exit(1)
				}

				rc := len(rules)
				if rc <= 0 {
					log.Printf("There are 0 NAT rules for '%s' - no policy was exported.", v)
					os.Remove(fh)
					os.Exit(1)
				}

				log.Printf("Exporting %d NAT rules - this might take a few of minutes if your rule base is large", rc)

				cfh.Write("#Name,Description,Type,SourceZones,DestinationZone,ToInterface,Service,SourceAddresses,DestinationAddresses,")
				cfh.Write("SatType,SatAddressType,SatTranslatedAddresses,SatInterface,SatIpAddress,SatFallbackType,SatFallbackTranslatedAddresses,")
				cfh.Write("SatFallbackInterface,SatFallbackIpType,SatFallbackIpAddress,SatStaticTranslatedAddress,SatStaticBiDirectional,DatType,")
				cfh.Write("DatAddress,DatPort,DatDynamicDistribution,Disabled,Tags\n")
				for _, rule := range rules {
					var toint string
					r, err := c.Policies.Nat.Get(v, rule)
					if err != nil {
						log.Printf("Failed to retrieve NAT rule data: %s", err)
					}

					toint = r.ToInterface
					if len(r.ToInterface) <= 0 {
						toint = "any"
					}

					cfh.Write(fmt.Sprintf("%s,\"%s\",%s,\"%s\",%s,%s,%s,\"%s\",\"%s\",", r.Name, r.Description, "ipv4", sliceToString(r.SourceZones),
						r.DestinationZone, toint, r.Service, sliceToString(r.SourceAddresses), sliceToString(r.DestinationAddresses)))
					cfh.Write(fmt.Sprintf("%s,%s,\"%s\",%s,%s,%s,\"%s\",%s,", r.SatType, r.SatAddressType, sliceToString(r.SatTranslatedAddresses), r.SatInterface,
						r.SatIpAddress, r.SatFallbackType, sliceToString(r.SatFallbackTranslatedAddresses), r.SatFallbackInterface))
					cfh.Write(fmt.Sprintf("%s,%s,%s,%t,%s,%s,%d,%s,%t,\"%s\"\n", r.SatFallbackIpType, r.SatFallbackIpAddress, r.SatStaticTranslatedAddress,
						r.SatStaticBiDirectional, r.DatType, r.DatAddress, r.DatPort, r.DatDynamicDistribution, r.Disabled, sliceToString(r.Tags)))
				}

				cfh.End()
			}

			if action == "import" && !xlate {
				rules, err := easycsv.Open(fh)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				rc := len(rules)
				log.Printf("Importing %d security rules - this might take a few of minutes if you have a lot of rules", rc)

				for i, rule := range rules {
					boolopt := map[string]bool{
						"TRUE":  true,
						"true":  true,
						"FALSE": false,
						"false": false,
					}

					e := security.Entry{
						Name:                            rule[0],
						Type:                            rule[1],
						Description:                     rule[2],
						Tags:                            stringToSlice(rule[3]),
						SourceZones:                     stringToSlice(rule[4]),
						SourceAddresses:                 stringToSlice(rule[5]),
						NegateSource:                    boolopt[rule[6]],
						SourceUsers:                     userStringToSlice(rule[7]),
						HipProfiles:                     stringToSlice(rule[8]),
						DestinationZones:                stringToSlice(rule[9]),
						DestinationAddresses:            stringToSlice(rule[10]),
						NegateDestination:               boolopt[rule[11]],
						Applications:                    stringToSlice(rule[12]),
						Services:                        stringToSlice(rule[13]),
						Categories:                      stringToSlice(rule[14]),
						Action:                          rule[15],
						LogSetting:                      rule[16],
						LogStart:                        boolopt[rule[17]],
						LogEnd:                          boolopt[rule[18]],
						Disabled:                        boolopt[rule[19]],
						Schedule:                        rule[20],
						IcmpUnreachable:                 boolopt[rule[21]],
						DisableServerResponseInspection: boolopt[rule[22]],
						Group:                           rule[23],
						Virus:                           rule[24],
						Spyware:                         rule[25],
						Vulnerability:                   rule[26],
						UrlFiltering:                    rule[27],
						FileBlocking:                    rule[28],
						WildFireAnalysis:                rule[29],
						DataFiltering:                   rule[30],
					}

					err = c.Policies.Security.Set(v, e)
					if err != nil {
						log.Printf("Line %d - failed to create security rule: %s", i+1, err)
					}
				}
			}

			if action == "import" && xlate {
				rules, err := easycsv.Open(fh)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				rc := len(rules)
				log.Printf("Importing %d NAT rules - this might take a few of minutes if you have a lot of rules", rc)

				for i, rule := range rules {
					boolopt := map[string]bool{
						"TRUE":  true,
						"true":  true,
						"FALSE": false,
						"false": false,
					}

					datport, _ := strconv.Atoi(rule[23])

					if rule[2] != "ipv4" {
						log.Printf("Line %d - only NAT type 'ipv4' is supported", i+1)
					}

					e := nat.Entry{
						Name:                           rule[0],
						Description:                    rule[1],
						Type:                           rule[2],
						SourceZones:                    stringToSlice(rule[3]),
						DestinationZone:                rule[4],
						ToInterface:                    rule[5],
						Service:                        rule[6],
						SourceAddresses:                stringToSlice(rule[7]),
						DestinationAddresses:           stringToSlice(rule[8]),
						SatType:                        rule[9],
						SatAddressType:                 rule[10],
						SatTranslatedAddresses:         stringToSlice(rule[11]),
						SatInterface:                   rule[12],
						SatIpAddress:                   rule[13],
						SatFallbackType:                rule[14],
						SatFallbackTranslatedAddresses: stringToSlice(rule[15]),
						SatFallbackInterface:           rule[16],
						SatFallbackIpType:              rule[17],
						SatFallbackIpAddress:           rule[18],
						SatStaticTranslatedAddress:     rule[19],
						SatStaticBiDirectional:         boolopt[rule[20]],
						DatType:                        rule[21],
						DatAddress:                     rule[22],
						DatPort:                        datport,
						DatDynamicDistribution:         rule[24],
						Disabled:                       boolopt[rule[25]],
						Tags:                           stringToSlice(rule[26]),
					}

					err = c.Policies.Nat.Set(v, e)
					if err != nil {
						log.Printf("Line %d - failed to create NAT rule: %s", i+1, err)
					}
				}
			}

			if action == "move" {
				moveOptions := map[string]int{
					"after":  util.MoveAfter,
					"before": util.MoveBefore,
					"bottom": util.MoveBottom,
					"top":    util.MoveTop,
				}

				if fh != "" && movemultiple {
					rules, err := easycsv.Open(fh)
					if err != nil {
						log.Printf("CSV file error - %s", err)
						os.Exit(1)
					}

					numrules := len(rules)
					log.Printf("Moving %d security rules", numrules)

					for _, rule := range rules {
						rulename := rule[0]
						ruledest := rule[1]
						targetrule := rule[2]
						loc := rule[5]

						r, err := c.Policies.Security.Get(loc, rulename)
						if err != nil {
							log.Printf("Failed to retrieve security rule: %s", err)
						}

						err = c.Policies.Security.MoveGroup(loc, moveOptions[ruledest], targetrule, r)
						if err != nil {
							fmt.Println(err)
						}
					}
				} else {
					rule, err := c.Policies.Security.Get(v, rulename)
					if err != nil {
						log.Printf("Failed to retrieve security rule: %s", err)
					}

					err = c.Policies.Security.MoveGroup(v, moveOptions[ruledest], targetrule, rule)
					if err != nil {
						fmt.Println(err)
					}
				}
			}

			if action == "grouptag" {
				// /config/devices/entry[@name='localhost.localdomain']/vsys/entry[@name='vsys1']/rulebase/security/rules/entry[@name='Block_Gaming']/group-tag
				// /config/devices/entry[@name='localhost.localdomain']/vsys/entry[@name='vsys1']/rulebase/nat/rules/entry[@name='Default_Outbound']/group-tag
				rules, err := easycsv.Open(fh)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				for i, rule := range rules {
					xpath := fmt.Sprintf("/config/devices/entry[@name='localhost.localdomain']/vsys/entry[@name='%s']/service-group/entry[@name='%s']", v, rule[0])
					ele := fmt.Sprintf("<group-tag>%s</group-tag>", rule[1])

					_, err := resty.R().Get(fmt.Sprintf("https://%s/api/?type=config&action=set&xpath=%s&key=%s&element=%s", device, xpath, c.ApiKey, ele))
					if err != nil {
						log.Printf("Line %d - failed to group rule by tag %s: %s", i+1, rule[0], err)
					}
				}
			}
		case *pango.Panorama:
			switch l {
			case "pre":
				l = util.PreRulebase
			case "post":
				l = util.PostRulebase
			default:
				l = util.PostRulebase
			}

			if action == "export" && !xlate {
				if !strings.Contains(fh, ".csv") {
					fh = fmt.Sprintf("%s.csv", fh)
				}

				cfh, err := easycsv.NewCSV(fh)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				rules, err := c.Policies.Security.GetList(dg, l)
				if err != nil {
					log.Printf("Failed to retrieve the list of security rules: %s", err)
					os.Remove(fh)
					os.Exit(1)
				}

				rc := len(rules)
				if rc <= 0 {
					log.Printf("There are 0 security rules for the '%s' device group - no policy was exported.", dg)
					os.Remove(fh)
					os.Exit(1)
				}

				log.Printf("Exporting %d security rules - this might take a few of minutes if your rule base is large", rc)

				cfh.Write("#Name,Type,Description,Tags,SourceZones,SourceAddresses,NegateSource,SourceUsers,HipProfiles,")
				cfh.Write("DestinationZones,DestinationAddresses,NegateDestination,Applications,Services,Categories,")
				cfh.Write("Action,LogSetting,LogStart,LogEnd,Disabled,Schedule,IcmpUnreachable,DisableServerResponseInspection,")
				cfh.Write("Group,Virus,Spyware,Vulnerability,UrlFiltering,FileBlocking,WildFireAnalysis,DataFiltering\n")
				for _, rule := range rules {
					var rtype string
					r, err := c.Policies.Security.Get(dg, l, rule)
					if err != nil {
						log.Printf("Failed to retrieve security rule data: %s", err)
					}

					switch r.Type {
					case "universal":
						rtype = "universal"
					case "intrazone":
						rtype = "intrazone"
					case "interzone":
						rtype = "interzone"
					default:
						rtype = "universal"
					}

					cfh.Write(fmt.Sprintf("%s,%s,\"%s\",\"%s\",\"%s\",\"%s\",%t,\"%s\",\"%s\",", r.Name, rtype, r.Description, sliceToString(r.Tags), sliceToString(r.SourceZones),
						sliceToString(r.SourceAddresses), r.NegateSource, userSliceToString(r.SourceUsers), sliceToString(r.HipProfiles)))
					cfh.Write(fmt.Sprintf("\"%s\",\"%s\",%t,\"%s\",\"%s\",\"%s\",", sliceToString(r.DestinationZones), sliceToString(r.DestinationAddresses), r.NegateDestination,
						sliceToString(r.Applications), sliceToString(r.Services), sliceToString(r.Categories)))
					cfh.Write(fmt.Sprintf("%s,%s,%t,%t,%t,%s,%t,%t,", r.Action, r.LogSetting, r.LogStart, r.LogEnd, r.Disabled, r.Schedule,
						r.IcmpUnreachable, r.DisableServerResponseInspection))
					cfh.Write(fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s\n", r.Group, r.Virus, r.Spyware,
						r.Vulnerability, r.UrlFiltering, r.FileBlocking, r.WildFireAnalysis, r.DataFiltering))
				}

				cfh.End()
			}

			if action == "export" && xlate {
				if !strings.Contains(fh, ".csv") {
					fh = fmt.Sprintf("%s.csv", fh)
				}

				cfh, err := easycsv.NewCSV(fh)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				rules, err := c.Policies.Nat.GetList(dg, l)
				if err != nil {
					log.Printf("Failed to retrieve the list of NAT rules: %s", err)
					os.Remove(fh)
					os.Exit(1)
				}

				rc := len(rules)
				if rc <= 0 {
					log.Printf("There are 0 NAT rules for the '%s' device group - no policy was exported.", dg)
					os.Remove(fh)
					os.Exit(1)
				}

				log.Printf("Exporting %d NAT rules - this might take a few of minutes if your rule base is large", rc)

				cfh.Write("#Name,Description,Type,SourceZones,DestinationZone,ToInterface,Service,SourceAddresses,DestinationAddresses,")
				cfh.Write("SatType,SatAddressType,SatTranslatedAddresses,SatInterface,SatIpAddress,SatFallbackType,SatFallbackTranslatedAddresses,")
				cfh.Write("SatFallbackInterface,SatFallbackIpType,SatFallbackIpAddress,SatStaticTranslatedAddress,SatStaticBiDirectional,DatType,")
				cfh.Write("DatAddress,DatPort,DatDynamicDistribution,Disabled,Tags\n")
				for _, rule := range rules {
					var toint string
					r, err := c.Policies.Nat.Get(dg, l, rule)
					if err != nil {
						log.Printf("Failed to retrieve NAT rule data: %s", err)
					}

					toint = r.ToInterface
					if len(r.ToInterface) <= 0 {
						toint = "any"
					}

					cfh.Write(fmt.Sprintf("%s,\"%s\",%s,\"%s\",%s,%s,%s,\"%s\",\"%s\",", r.Name, r.Description, "ipv4", sliceToString(r.SourceZones),
						r.DestinationZone, toint, r.Service, sliceToString(r.SourceAddresses), sliceToString(r.DestinationAddresses)))
					cfh.Write(fmt.Sprintf("%s,%s,\"%s\",%s,%s,%s,\"%s\",%s,", r.SatType, r.SatAddressType, sliceToString(r.SatTranslatedAddresses), r.SatInterface,
						r.SatIpAddress, r.SatFallbackType, sliceToString(r.SatFallbackTranslatedAddresses), r.SatFallbackInterface))
					cfh.Write(fmt.Sprintf("%s,%s,%s,%t,%s,%s,%d,%s,%t,\"%s\"\n", r.SatFallbackIpType, r.SatFallbackIpAddress, r.SatStaticTranslatedAddress,
						r.SatStaticBiDirectional, r.DatType, r.DatAddress, r.DatPort, r.DatDynamicDistribution, r.Disabled, sliceToString(r.Tags)))
				}

				cfh.End()
			}

			if action == "import" && !xlate {
				rules, err := easycsv.Open(fh)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				rc := len(rules)
				log.Printf("Importing %d security rules - this might take a few of minutes if you have a lot of rules", rc)

				for i, rule := range rules {
					boolopt := map[string]bool{
						"TRUE":  true,
						"true":  true,
						"FALSE": false,
						"false": false,
					}

					e := security.Entry{
						Name:                            rule[0],
						Type:                            rule[1],
						Description:                     rule[2],
						Tags:                            stringToSlice(rule[3]),
						SourceZones:                     stringToSlice(rule[4]),
						SourceAddresses:                 stringToSlice(rule[5]),
						NegateSource:                    boolopt[rule[6]],
						SourceUsers:                     userStringToSlice(rule[7]),
						HipProfiles:                     stringToSlice(rule[8]),
						DestinationZones:                stringToSlice(rule[9]),
						DestinationAddresses:            stringToSlice(rule[10]),
						NegateDestination:               boolopt[rule[11]],
						Applications:                    stringToSlice(rule[12]),
						Services:                        stringToSlice(rule[13]),
						Categories:                      stringToSlice(rule[14]),
						Action:                          rule[15],
						LogSetting:                      rule[16],
						LogStart:                        boolopt[rule[17]],
						LogEnd:                          boolopt[rule[18]],
						Disabled:                        boolopt[rule[19]],
						Schedule:                        rule[20],
						IcmpUnreachable:                 boolopt[rule[21]],
						DisableServerResponseInspection: boolopt[rule[22]],
						Group:                           rule[23],
						Virus:                           rule[24],
						Spyware:                         rule[25],
						Vulnerability:                   rule[26],
						UrlFiltering:                    rule[27],
						FileBlocking:                    rule[28],
						WildFireAnalysis:                rule[29],
						DataFiltering:                   rule[30],
					}

					err = c.Policies.Security.Set(dg, l, e)
					if err != nil {
						log.Printf("Line %d - failed to create security rule: %s", i+1, err)
					}
				}
			}

			if action == "import" && xlate {
				rules, err := easycsv.Open(fh)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				rc := len(rules)
				log.Printf("Importing %d NAT rules - this might take a few of minutes if you have a lot of rules", rc)

				for i, rule := range rules {
					boolopt := map[string]bool{
						"TRUE":  true,
						"true":  true,
						"FALSE": false,
						"false": false,
					}

					datport, _ := strconv.Atoi(rule[23])

					if rule[2] != "ipv4" {
						log.Printf("Line %d - only NAT type 'ipv4' is supported", i+1)
					}

					e := nat.Entry{
						Name:                           rule[0],
						Description:                    rule[1],
						Type:                           rule[2],
						SourceZones:                    stringToSlice(rule[3]),
						DestinationZone:                rule[4],
						ToInterface:                    rule[5],
						Service:                        rule[6],
						SourceAddresses:                stringToSlice(rule[7]),
						DestinationAddresses:           stringToSlice(rule[8]),
						SatType:                        rule[9],
						SatAddressType:                 rule[10],
						SatTranslatedAddresses:         stringToSlice(rule[11]),
						SatInterface:                   rule[12],
						SatIpAddress:                   rule[13],
						SatFallbackType:                rule[14],
						SatFallbackTranslatedAddresses: stringToSlice(rule[15]),
						SatFallbackInterface:           rule[16],
						SatFallbackIpType:              rule[17],
						SatFallbackIpAddress:           rule[18],
						SatStaticTranslatedAddress:     rule[19],
						SatStaticBiDirectional:         boolopt[rule[20]],
						DatType:                        rule[21],
						DatAddress:                     rule[22],
						DatPort:                        datport,
						DatDynamicDistribution:         rule[24],
						Disabled:                       boolopt[rule[25]],
						Tags:                           stringToSlice(rule[26]),
					}

					err = c.Policies.Nat.Set(dg, l, e)
					if err != nil {
						log.Printf("Line %d - failed to create NAT rule: %s", i+1, err)
					}
				}
			}

			if action == "move" {
				moveOptions := map[string]int{
					"after":  util.MoveAfter,
					"before": util.MoveBefore,
					"bottom": util.MoveBottom,
					"top":    util.MoveTop,
				}

				if fh != "" && movemultiple {
					rules, err := easycsv.Open(fh)
					if err != nil {
						log.Printf("CSV file error - %s", err)
						os.Exit(1)
					}

					numrules := len(rules)
					log.Printf("Moving %d security rules", numrules)

					for _, rule := range rules {
						rulename := rule[0]
						ruledest := rule[1]
						targetrule := rule[2]
						dg := rule[5]

						r, err := c.Policies.Security.Get(dg, l, rulename)
						if err != nil {
							log.Printf("Failed to retrieve security rule: %s", err)
						}

						err = c.Policies.Security.MoveGroup(dg, l, moveOptions[ruledest], targetrule, r)
						if err != nil {
							fmt.Println(err)
						}
					}
				} else {
					rule, err := c.Policies.Security.Get(dg, l, rulename)
					if err != nil {
						log.Printf("Failed to retrieve security rule: %s", err)
					}

					err = c.Policies.Security.MoveGroup(dg, l, moveOptions[ruledest], targetrule, rule)
					if err != nil {
						fmt.Println(err)
					}
				}
			}

			// if action == "grouptag" {
			// /config/devices/entry[@name='localhost.localdomain']/vsys/entry[@name='vsys1']/rulebase/security/rules/entry[@name='Block_Gaming']/group-tag
			// /config/devices/entry[@name='localhost.localdomain']/vsys/entry[@name='vsys1']/rulebase/nat/rules/entry[@name='Default_Outbound']/group-tag
			// }
		}
	},
}

func init() {
	rootCmd.AddCommand(policyCmd)

	policyCmd.Flags().StringVarP(&action, "action", "a", "", "Action to perform - import, export, or move (only for security policy)")
	policyCmd.Flags().StringVarP(&fh, "file", "f", "", "Name of the CSV file to import/export to")
	policyCmd.Flags().StringVarP(&dg, "devicegroup", "g", "", "Device Group name; only needed when ran against Panorama")
	policyCmd.Flags().StringVarP(&user, "user", "u", "", "User to connect to the device as")
	policyCmd.Flags().StringVarP(&pass, "pass", "p", "", "Password for the user account specified")
	policyCmd.Flags().StringVarP(&device, "device", "d", "", "Firewall or Panorama device to connect to")
	policyCmd.Flags().StringVarP(&l, "location", "l", "post", "Rule location; pre or post when ran against Panorama")
	policyCmd.Flags().BoolVarP(&xlate, "nat", "x", false, "Run the given action on the NAT policy")
	policyCmd.Flags().StringVarP(&v, "vsys", "v", "vsys1", "Vsys name when ran against a firewall")
	policyCmd.Flags().StringVarP(&rulename, "rulename", "n", "", "Name of the security rule you wish to move")
	policyCmd.Flags().StringVarP(&ruledest, "ruledest", "w", "", "Where to move the rule - after, before, top, or bottom")
	policyCmd.Flags().StringVarP(&targetrule, "targetrule", "t", "", "Name of the rule 'ruledest' is referencing")
	policyCmd.Flags().BoolVarP(&movemultiple, "movemultiple", "m", true, "Specifies you wish to move multiple security rules; use only with --file")
	policyCmd.MarkFlagRequired("user")
	policyCmd.MarkFlagRequired("pass")
	policyCmd.MarkFlagRequired("device")
	policyCmd.MarkFlagRequired("action")
}
