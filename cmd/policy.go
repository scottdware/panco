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
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/poli/security"
	"github.com/PaloAltoNetworks/pango/util"
	easycsv "github.com/scottdware/go-easycsv"
	"github.com/spf13/cobra"
)

// policyCmd represents the policy command
var policyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Import and export a security policy",
	Long: `This command will allow you to import and export an entire security policy. If
you are running this against a Panorama device, it can be really helpful if you want to clone
an entire policy, as you can export it from one device-group, modify it if needed, then import
the poilcy into a different device-group (or firewall).`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		pass := passwd()

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
				if !strings.Contains(fh, ".csv") {
					fh = fmt.Sprintf("%s.csv", fh)
				}

				cfh, err := easycsv.NewCSV(fh)
				if err != nil {
					log.Printf("Failed to create CSV file %s: %s", fh, err)
					os.Exit(1)
				}

				rules, err := c.Policies.Security.GetList(v)
				if err != nil {
					log.Printf("Failed to retrieve the list of rules: %s", err)
				}

				rc := len(rules)
				log.Printf("Exporting %d rules - this might take a few of minutes if your rule base is large", rc)

				cfh.Write("#Name,Type,Description,Tags,SourceZones,SourceAddresses,NegateSource,SourceUsers,HipProfiles,")
				cfh.Write("DestinationZones,DestinationAddresses,NegateDestination,Applications,Services,Categories,")
				cfh.Write("Action,LogSetting,LogStart,LogEnd,Disabled,Schedule,IcmpUnreachable,DisableServerResponseInspection,")
				cfh.Write("Group,Targets,NegateTarget,Virus,Spyware,Vulnerability,UrlFiltering,FileBlocking,WildFireAnalysis,DataFiltering\n")
				for _, rule := range rules {
					var rtype string
					r, err := c.Policies.Security.Get(v, rule)
					if err != nil {
						log.Printf("Failed to retrieve rule data for '%s': %s", rule, err)
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
					cfh.Write(fmt.Sprintf("%s,\"%s\",%t,%s,%s,%s,%s,%s,%s,%s\n", r.Group, r.Targets, r.NegateTarget, r.Virus, r.Spyware,
						r.Vulnerability, r.UrlFiltering, r.FileBlocking, r.WildFireAnalysis, r.DataFiltering))
				}

				cfh.End()
			}

			if action == "import" {
				rules, err := easycsv.Open(fh)
				if err != nil {
					log.Printf("Failed to open CSV file: %s", err)
					os.Exit(1)
				}

				rc := len(rules)
				log.Printf("Importing %d rules - this might take a few of minutes if you have a lot of rules", rc)

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
						NegateTarget:                    boolopt[rule[25]],
						Virus:                           rule[26],
						Spyware:                         rule[27],
						Vulnerability:                   rule[28],
						UrlFiltering:                    rule[29],
						FileBlocking:                    rule[30],
						WildFireAnalysis:                rule[31],
						DataFiltering:                   rule[32],
					}

					err = c.Policies.Security.Set(v, e)
					if err != nil {
						log.Printf("Line %d - failed to create rule: %s", i+1, err)
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

			if action == "export" {
				if !strings.Contains(fh, ".csv") {
					fh = fmt.Sprintf("%s.csv", fh)
				}

				cfh, err := easycsv.NewCSV(fh)
				if err != nil {
					log.Printf("Failed to create CSV file %s: %s", fh, err)
					os.Exit(1)
				}

				rules, err := c.Policies.Security.GetList(dg, l)
				if err != nil {
					log.Printf("Failed to retrieve the list of rules: %s", err)
				}

				rc := len(rules)
				log.Printf("Exporting %d rules - this might take a few of minutes if your rule base is large", rc)

				cfh.Write("#Name,Type,Description,Tags,SourceZones,SourceAddresses,NegateSource,SourceUsers,HipProfiles,")
				cfh.Write("DestinationZones,DestinationAddresses,NegateDestination,Applications,Services,Categories,")
				cfh.Write("Action,LogSetting,LogStart,LogEnd,Disabled,Schedule,IcmpUnreachable,DisableServerResponseInspection,")
				cfh.Write("Group,Targets,NegateTarget,Virus,Spyware,Vulnerability,UrlFiltering,FileBlocking,WildFireAnalysis,DataFiltering\n")
				for _, rule := range rules {
					var rtype string
					r, err := c.Policies.Security.Get(dg, l, rule)
					if err != nil {
						log.Printf("Failed to retrieve rule data for '%s': %s", rule, err)
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
					cfh.Write(fmt.Sprintf("%s,\"%s\",%t,%s,%s,%s,%s,%s,%s,%s\n", r.Group, r.Targets, r.NegateTarget, r.Virus, r.Spyware,
						r.Vulnerability, r.UrlFiltering, r.FileBlocking, r.WildFireAnalysis, r.DataFiltering))
				}

				cfh.End()
			}

			if action == "import" {
				rules, err := easycsv.Open(fh)
				if err != nil {
					log.Printf("Failed to open CSV file: %s", err)
					os.Exit(1)
				}

				rc := len(rules)
				log.Printf("Importing %d rules - this might take a few of minutes if you have a lot of rules", rc)

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
						NegateTarget:                    boolopt[rule[25]],
						Virus:                           rule[26],
						Spyware:                         rule[27],
						Vulnerability:                   rule[28],
						UrlFiltering:                    rule[29],
						FileBlocking:                    rule[30],
						WildFireAnalysis:                rule[31],
						DataFiltering:                   rule[32],
					}

					err = c.Policies.Security.Set(dg, l, e)
					if err != nil {
						log.Printf("Line %d - failed to create rule: %s", i+1, err)
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(policyCmd)

	policyCmd.Flags().StringVarP(&action, "action", "a", "", "Action to perform; import or export")
	policyCmd.Flags().StringVarP(&fh, "file", "f", "", "Name of the CSV file to import/export to")
	policyCmd.Flags().StringVarP(&dg, "devicegroup", "g", "", "Device Group name; only needed when ran against Panorama")
	policyCmd.Flags().StringVarP(&user, "user", "u", "", "User to connect to the device as")
	policyCmd.Flags().StringVarP(&device, "device", "d", "", "Firewall or Panorama device to connect to")
	policyCmd.Flags().StringVarP(&l, "location", "l", "post", "Rule location; pre or post when ran against Panorama")
	policyCmd.Flags().StringVarP(&v, "vsys", "v", "vsys1", "Vsys name when ran against a firewall")
	policyCmd.MarkFlagRequired("user")
	policyCmd.MarkFlagRequired("device")
	policyCmd.MarkFlagRequired("action")
	policyCmd.MarkFlagRequired("file")
}
