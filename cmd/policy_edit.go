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
	"strconv"
	"strings"
	"time"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/poli/decryption"
	"github.com/PaloAltoNetworks/pango/poli/nat"
	"github.com/PaloAltoNetworks/pango/poli/pbf"
	"github.com/PaloAltoNetworks/pango/poli/security"
	"github.com/PaloAltoNetworks/pango/util"
	"github.com/Songmu/prompter"
	easycsv "github.com/scottdware/go-easycsv"
	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
)

// modifyCmd represents the import command
var policyEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a Security, NAT, Decryption or Policy-Based Forwarding policy",
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
			if t == "security" {
				rules, err := easycsv.Open(f)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				rc := len(rules)
				log.Printf("Editing %d Security rules", rc)

				for i, rule := range rules {
					boolopt := map[string]bool{
						"TRUE":  true,
						"true":  true,
						"FALSE": false,
						"false": false,
					}

					// ruletype := rule[1]
					// apps := rule[12]

					// if len(ruletype) <= 0 {
					// 	ruletype = "universal"
					// }

					// if len(apps) <= 0 {
					// 	apps = "any"
					// }

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

					err = c.Policies.Security.Edit(v, e)
					if err != nil {
						log.Printf("Line %d - failed to edit Security rule: %s", i+1, err)
					}

					time.Sleep(delay * time. Millisecond)
				}
			}

			if t == "nat" {
				rules, err := easycsv.Open(f)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				rc := len(rules)
				log.Printf("Editing %d NAT rules", rc)

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

					// ruletype := rule[1]

					// if len(ruletype) <= 0 {
					// 	ruletype = "universal"
					// }

					e := nat.Entry{
						Name:                           rule[0],
						Type:                           rule[1],
						Description:                    rule[2],
						Tags:                           stringToSlice(rule[3]),
						SourceZones:                    stringToSlice(rule[4]),
						DestinationZone:                rule[5],
						ToInterface:                    rule[6],
						Service:                        rule[7],
						SourceAddresses:                stringToSlice(rule[8]),
						DestinationAddresses:           stringToSlice(rule[9]),
						SatType:                        rule[10],
						SatAddressType:                 rule[11],
						SatTranslatedAddresses:         stringToSlice(rule[12]),
						SatInterface:                   rule[13],
						SatIpAddress:                   rule[14],
						SatFallbackType:                rule[15],
						SatFallbackTranslatedAddresses: stringToSlice(rule[16]),
						SatFallbackInterface:           rule[17],
						SatFallbackIpType:              rule[18],
						SatFallbackIpAddress:           rule[19],
						SatStaticTranslatedAddress:     rule[20],
						SatStaticBiDirectional:         boolopt[rule[21]],
						DatType:                        rule[22],
						DatAddress:                     rule[23],
						DatPort:                        datport,
						DatDynamicDistribution:         rule[25],
						Disabled:                       boolopt[rule[26]],
					}

					err = c.Policies.Nat.Edit(v, e)
					if err != nil {
						log.Printf("Line %d - failed to edit NAT rule: %s", i+1, err)
					}

					time.Sleep(delay * time. Millisecond)
				}
			}

			if t == "pbf" {
				rules, err := easycsv.Open(f)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				rc := len(rules)
				log.Printf("Editing %d Policy-Based Forwarding rules", rc)

				for i, rule := range rules {
					boolopt := map[string]bool{
						"TRUE":  true,
						"true":  true,
						"FALSE": false,
						"false": false,
					}

					e := pbf.Entry{
						Name:                               rule[0],
						Description:                        rule[1],
						Tags:                               stringToSlice(rule[2]),
						FromType:                           rule[3],
						FromValues:                         stringToSlice(rule[4]),
						SourceAddresses:                    stringToSlice(rule[5]),
						SourceUsers:                        stringToSlice(rule[6]),
						NegateSource:                       boolopt[rule[7]],
						DestinationAddresses:               stringToSlice(rule[8]),
						NegateDestination:                  boolopt[rule[9]],
						Applications:                       stringToSlice(rule[10]),
						Services:                           stringToSlice(rule[11]),
						Schedule:                           rule[12],
						Disabled:                           boolopt[rule[13]],
						Action:                             rule[14],
						ForwardVsys:                        rule[15],
						ForwardEgressInterface:             rule[16],
						ForwardNextHopType:                 rule[17],
						ForwardNextHopValue:                rule[18],
						ForwardMonitorProfile:              rule[19],
						ForwardMonitorIpAddress:            rule[20],
						ForwardMonitorDisableIfUnreachable: boolopt[rule[21]],
						EnableEnforceSymmetricReturn:       boolopt[rule[22]],
						SymmetricReturnAddresses:           stringToSlice(rule[23]),
						ActiveActiveDeviceBinding:          rule[24],
						NegateTarget:                       boolopt[rule[25]],
						Uuid:                               rule[26],
					}

					err = c.Policies.PolicyBasedForwarding.Edit(v, e)
					if err != nil {
						log.Printf("Line %d - failed to edit Policy-Based Forwarding rule: %s", i+1, err)
					}

					time.Sleep(delay * time. Millisecond)
				}
			}

			if t == "decrypt" {
				rules, err := easycsv.Open(f)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				rc := len(rules)
				log.Printf("Editing %d Decryption rules", rc)

				for i, rule := range rules {
					boolopt := map[string]bool{
						"TRUE":  true,
						"true":  true,
						"FALSE": false,
						"false": false,
					}

					e := decryption.Entry{
						Name:                       strings.TrimSpace(rule[0]),
						Description:                rule[1],
						SourceZones:                stringToSlice(rule[2]),
						SourceAddresses:            stringToSlice(rule[3]),
						NegateSource:               boolopt[rule[4]],
						SourceUsers:                userStringToSlice(rule[5]),
						DestinationZones:           stringToSlice(rule[6]),
						DestinationAddresses:       stringToSlice(rule[7]),
						NegateDestination:          boolopt[rule[8]],
						Tags:                       stringToSlice(rule[9]),
						Disabled:                   boolopt[rule[10]],
						Services:                   stringToSlice(rule[11]),
						UrlCategories:              stringToSlice(rule[12]),
						Action:                     rule[13],
						DecryptionType:             rule[14],
						SslCertificate:             rule[15],
						DecryptionProfile:          rule[16],
						NegateTarget:               boolopt[rule[17]],
						ForwardingProfile:          rule[18],
						GroupTag:                   rule[19],
						SourceHips:                 stringToSlice(rule[20]),
						DestinationHips:            stringToSlice(rule[21]),
						LogSuccessfulTlsHandshakes: boolopt[rule[22]],
						LogFailedTlsHandshakes:     boolopt[rule[23]],
						LogSetting:                 rule[24],
						SslCertificates:            stringToSlice(rule[25]),
					}

					err = c.Policies.Decryption.Edit(v, e)
					if err != nil {
						log.Printf("Line %d - failed to edit Decryption rule: %s", i+1, err)
					}

					time.Sleep(delay * time. Millisecond)
				}
			}
		case *pango.Panorama:
			switch l {
			case "pre":
				l = util.PreRulebase
			case "post":
				l = util.PostRulebase
			default:
				l = util.PreRulebase
			}

			if t == "security" {
				rules, err := easycsv.Open(f)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				rc := len(rules)
				log.Printf("Editing %d Security rules", rc)

				for i, rule := range rules {
					boolopt := map[string]bool{
						"TRUE":  true,
						"true":  true,
						"FALSE": false,
						"false": false,
					}

					// ruletype := rule[1]
					// apps := rule[12]

					// if len(ruletype) <= 0 {
					// 	ruletype = "universal"
					// }

					// if len(apps) <= 0 {
					// 	apps = "any"
					// }

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

					err = c.Policies.Security.Edit(dg, l, e)
					if err != nil {
						log.Printf("Line %d - failed to edit security rule: %s", i+1, err)
					}

					time.Sleep(delay * time. Millisecond)
				}
			}

			if t == "nat" {
				rules, err := easycsv.Open(f)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				rc := len(rules)
				log.Printf("Editing %d NAT rules", rc)

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

					// ruletype := rule[1]

					// if len(ruletype) <= 0 {
					// 	ruletype = "universal"
					// }

					e := nat.Entry{
						Name:                           rule[0],
						Type:                           rule[1],
						Description:                    rule[2],
						Tags:                           stringToSlice(rule[3]),
						SourceZones:                    stringToSlice(rule[4]),
						DestinationZone:                rule[5],
						ToInterface:                    rule[6],
						Service:                        rule[7],
						SourceAddresses:                stringToSlice(rule[8]),
						DestinationAddresses:           stringToSlice(rule[9]),
						SatType:                        rule[10],
						SatAddressType:                 rule[11],
						SatTranslatedAddresses:         stringToSlice(rule[12]),
						SatInterface:                   rule[13],
						SatIpAddress:                   rule[14],
						SatFallbackType:                rule[15],
						SatFallbackTranslatedAddresses: stringToSlice(rule[16]),
						SatFallbackInterface:           rule[17],
						SatFallbackIpType:              rule[18],
						SatFallbackIpAddress:           rule[19],
						SatStaticTranslatedAddress:     rule[20],
						SatStaticBiDirectional:         boolopt[rule[21]],
						DatType:                        rule[22],
						DatAddress:                     rule[23],
						DatPort:                        datport,
						DatDynamicDistribution:         rule[25],
						Disabled:                       boolopt[rule[26]],
					}

					err = c.Policies.Nat.Edit(dg, l, e)
					if err != nil {
						log.Printf("Line %d - failed to edit NAT rule: %s", i+1, err)
					}

					time.Sleep(delay * time. Millisecond)
				}
			}

			if t == "pbf" {
				rules, err := easycsv.Open(f)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				rc := len(rules)
				log.Printf("Editing %d Policy-Based Forwarding rules", rc)

				for i, rule := range rules {
					boolopt := map[string]bool{
						"TRUE":  true,
						"true":  true,
						"FALSE": false,
						"false": false,
					}

					e := pbf.Entry{
						Name:                               rule[0],
						Description:                        rule[1],
						Tags:                               stringToSlice(rule[2]),
						FromType:                           rule[3],
						FromValues:                         stringToSlice(rule[4]),
						SourceAddresses:                    stringToSlice(rule[5]),
						SourceUsers:                        stringToSlice(rule[6]),
						NegateSource:                       boolopt[rule[7]],
						DestinationAddresses:               stringToSlice(rule[8]),
						NegateDestination:                  boolopt[rule[9]],
						Applications:                       stringToSlice(rule[10]),
						Services:                           stringToSlice(rule[11]),
						Schedule:                           rule[12],
						Disabled:                           boolopt[rule[13]],
						Action:                             rule[14],
						ForwardVsys:                        rule[15],
						ForwardEgressInterface:             rule[16],
						ForwardNextHopType:                 rule[17],
						ForwardNextHopValue:                rule[18],
						ForwardMonitorProfile:              rule[19],
						ForwardMonitorIpAddress:            rule[20],
						ForwardMonitorDisableIfUnreachable: boolopt[rule[21]],
						EnableEnforceSymmetricReturn:       boolopt[rule[22]],
						SymmetricReturnAddresses:           stringToSlice(rule[23]),
						ActiveActiveDeviceBinding:          rule[24],
						NegateTarget:                       boolopt[rule[25]],
						Uuid:                               rule[26],
					}

					err = c.Policies.PolicyBasedForwarding.Edit(dg, l, e)
					if err != nil {
						log.Printf("Line %d - failed to edit Policy-Based Forwarding rule: %s", i+1, err)
					}

					time.Sleep(delay * time. Millisecond)
				}
			}

			if t == "decrypt" {
				rules, err := easycsv.Open(f)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				rc := len(rules)
				log.Printf("Editing %d Decryption rules", rc)

				for i, rule := range rules {
					boolopt := map[string]bool{
						"TRUE":  true,
						"true":  true,
						"FALSE": false,
						"false": false,
					}

					e := decryption.Entry{
						Name:                       strings.TrimSpace(rule[0]),
						Description:                rule[1],
						SourceZones:                stringToSlice(rule[2]),
						SourceAddresses:            stringToSlice(rule[3]),
						NegateSource:               boolopt[rule[4]],
						SourceUsers:                userStringToSlice(rule[5]),
						DestinationZones:           stringToSlice(rule[6]),
						DestinationAddresses:       stringToSlice(rule[7]),
						NegateDestination:          boolopt[rule[8]],
						Tags:                       stringToSlice(rule[9]),
						Disabled:                   boolopt[rule[10]],
						Services:                   stringToSlice(rule[11]),
						UrlCategories:              stringToSlice(rule[12]),
						Action:                     rule[13],
						DecryptionType:             rule[14],
						SslCertificate:             rule[15],
						DecryptionProfile:          rule[16],
						NegateTarget:               boolopt[rule[17]],
						ForwardingProfile:          rule[18],
						GroupTag:                   rule[19],
						SourceHips:                 stringToSlice(rule[20]),
						DestinationHips:            stringToSlice(rule[21]),
						LogSuccessfulTlsHandshakes: boolopt[rule[22]],
						LogFailedTlsHandshakes:     boolopt[rule[23]],
						LogSetting:                 rule[24],
						SslCertificates:            stringToSlice(rule[25]),
					}

					err = c.Policies.Decryption.Edit(dg, l, e)
					if err != nil {
						log.Printf("Line %d - failed to edit Decryption rule: %s", i+1, err)
					}

					time.Sleep(delay * time. Millisecond)
				}
			}
		}
	},
}

func init() {
	policyCmd.AddCommand(policyEditCmd)

	policyEditCmd.Flags().StringVarP(&user, "user", "u", "", "User to connect to the device as")
    policyEditCmd.Flags().StringVarP(&p, "delay", "p", "100", "Delay (in milliseconds) to pause between each API call")
	// policyModifyCmd.Flags().StringVarP(&pass, "pass", "p", "", "Password for the user account specified")
	policyEditCmd.Flags().StringVarP(&device, "device", "d", "", "Device to connect to")
	policyEditCmd.Flags().StringVarP(&f, "file", "f", "", "Name of the CSV file to export to")
	policyEditCmd.Flags().StringVarP(&dg, "devicegroup", "g", "shared", "Device Group name when importing to Panorama")
	policyEditCmd.Flags().StringVarP(&v, "vsys", "v", "vsys1", "Vsys name when importing to a firewall")
	policyEditCmd.Flags().StringVarP(&t, "type", "t", "", "Type of policy to import - <security|nat|decrypt|pbf>")
	policyEditCmd.Flags().StringVarP(&l, "location", "l", "pre", "Location of the rulebase - <pre|post>")
	policyEditCmd.MarkFlagRequired("user")
	// policyModifyCmd.MarkFlagRequired("pass")
	policyEditCmd.MarkFlagRequired("device")
	policyEditCmd.MarkFlagRequired("file")
	policyEditCmd.MarkFlagRequired("type")
	// policyModifyCmd.MarkFlagRequired("location")
}
