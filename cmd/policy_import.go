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

// importCmd represents the import command
var policyImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import (create, modify) a Security, NAT, Decryption or Policy-Based Forwarding policy",
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
			timeoutCount := 0
			timeoutData := []string{}

			if t == "security" {
				rules, err := easycsv.Open(f)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				rc := len(rules)
				log.Printf("Importing/modifying %d Security rules", rc)

				for i, rule := range rules {
					boolopt := map[string]bool{
						"TRUE":  true,
						"true":  true,
						"FALSE": false,
						"false": false,
					}

					v := rule[0]

					// ruletype := rule[1]
					// apps := rule[12]

					// if len(ruletype) <= 0 {
					// 	ruletype = "universal"
					// }

					// if len(apps) <= 0 {
					// 	apps = "any"
					// }

					if len(strings.TrimSpace(rule[2])) > 63 {
						log.Printf("Line %d - failed to create %s: name is over the max 63 characters", i+1, strings.TrimSpace(rule[0]))
					} else {
						e := security.Entry{
							Name:                            strings.TrimSpace(rule[0]),
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
							if strings.Contains(err.Error(), "Client.Timeout") {
								timeoutCount++
								timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, strings.TrimSpace(rule[0])))
							} else {
								log.Printf("Line %d - failed to create Security rule: %s", i+1, err)
							}
						}
					}

					time.Sleep(delay * time.Millisecond)
				}
			}

			if t == "nat" {
				rules, err := easycsv.Open(f)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				rc := len(rules)
				log.Printf("Importing/modifying %d NAT rules", rc)

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

					ruletype := rule[1]

					if len(ruletype) <= 0 {
						ruletype = "universal"
					}

					if len(strings.TrimSpace(rule[0])) > 63 {
						log.Printf("Line %d - failed to create %s: name is over the max 63 characters", i+1, strings.TrimSpace(rule[0]))
					} else {
						e := nat.Entry{
							Name:                           strings.TrimSpace(rule[0]),
							Type:                           ruletype,
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

						err = c.Policies.Nat.Set(v, e)
						if err != nil {
							if strings.Contains(err.Error(), "Client.Timeout") {
								timeoutCount++
								timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, strings.TrimSpace(rule[0])))
							} else {
								log.Printf("Line %d - failed to create NAT rule: %s", i+1, err)
							}
						}
					}

					time.Sleep(delay * time.Millisecond)
				}
			}

			if t == "pbf" {
				rules, err := easycsv.Open(f)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				rc := len(rules)
				log.Printf("Importing/modifying %d Policy-Based Forwarding rules", rc)

				for i, rule := range rules {
					boolopt := map[string]bool{
						"TRUE":  true,
						"true":  true,
						"FALSE": false,
						"false": false,
					}

					if len(strings.TrimSpace(rule[0])) > 63 {
						log.Printf("Line %d - failed to create %s: name is over the max 63 characters", i+1, strings.TrimSpace(rule[0]))
					} else {
						e := pbf.Entry{
							Name:                               strings.TrimSpace(rule[0]),
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

						err = c.Policies.PolicyBasedForwarding.Set(v, e)
						if err != nil {
							if strings.Contains(err.Error(), "Client.Timeout") {
								timeoutCount++
								timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, strings.TrimSpace(rule[0])))
							} else {
								log.Printf("Line %d - failed to create Policy-Based Forwarding rule: %s", i+1, err)
							}
						}
					}

					time.Sleep(delay * time.Millisecond)
				}
			}

			if t == "decrypt" {
				rules, err := easycsv.Open(f)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				rc := len(rules)
				log.Printf("Importing/modifying %d Decryption rules", rc)

				for i, rule := range rules {
					boolopt := map[string]bool{
						"TRUE":  true,
						"true":  true,
						"FALSE": false,
						"false": false,
					}

					if len(strings.TrimSpace(rule[0])) > 63 {
						log.Printf("Line %d - failed to create %s: name is over the max 63 characters", i+1, strings.TrimSpace(rule[0]))
					} else {
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

						err = c.Policies.Decryption.Set(v, e)
						if err != nil {
							if strings.Contains(err.Error(), "Client.Timeout") {
								timeoutCount++
								timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, strings.TrimSpace(rule[0])))
							} else {
								log.Printf("Line %d - failed to create Decryption rule: %s", i+1, err)
							}
						}
					}

					time.Sleep(delay * time.Millisecond)
				}
			}

			if timeoutCount > 0 {
				log.Printf("There were %d API timeout errors during import. Please verify the following have been imported/modified:\n\n", timeoutCount)
				for _, data := range timeoutData {
					info := strings.Split(data, ":")
					fmt.Printf("Line %s: Rule \"%s\"\n", info[0], info[1])
				}
			}
		case *pango.Panorama:
			// switch l {
			// case "pre":
			// 	l = util.PreRulebase
			// case "post":
			// 	l = util.PostRulebase
			// default:
			// 	l = util.PreRulebase
			// }

			timeoutCount := 0
			timeoutData := []string{}

			if t == "security" {
				rules, err := easycsv.Open(f)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				rc := len(rules)
				log.Printf("Importing/modifying %d Security rules", rc)

				for i, rule := range rules {
					boolopt := map[string]bool{
						"TRUE":  true,
						"true":  true,
						"FALSE": false,
						"false": false,
					}

					dg := rule[0]
					l := rule[1]

					switch l {
					case "pre":
						l = util.PreRulebase
					case "post":
						l = util.PostRulebase
					default:
						l = util.PreRulebase
					}

					ruletype := rule[3]
					apps := rule[14]

					if len(ruletype) <= 0 {
						ruletype = "universal"
					}

					if len(apps) <= 0 {
						apps = "any"
					}

					if len(strings.TrimSpace(rule[2])) > 63 {
						log.Printf("Line %d - failed to create %s: name is over the max 63 characters", i+1, strings.TrimSpace(rule[2]))
					} else {
						e := security.Entry{
							Name:                            strings.TrimSpace(rule[2]),
							Type:                            ruletype,
							Description:                     rule[4],
							Tags:                            stringToSlice(rule[5]),
							SourceZones:                     stringToSlice(rule[6]),
							SourceAddresses:                 stringToSlice(rule[7]),
							NegateSource:                    boolopt[rule[8]],
							SourceUsers:                     userStringToSlice(rule[9]),
							HipProfiles:                     stringToSlice(rule[10]),
							DestinationZones:                stringToSlice(rule[11]),
							DestinationAddresses:            stringToSlice(rule[12]),
							NegateDestination:               boolopt[rule[13]],
							Applications:                    stringToSlice(apps),
							Services:                        stringToSlice(rule[15]),
							Categories:                      stringToSlice(rule[16]),
							Action:                          rule[17],
							LogSetting:                      rule[18],
							LogStart:                        boolopt[rule[19]],
							LogEnd:                          boolopt[rule[20]],
							Disabled:                        boolopt[rule[21]],
							Schedule:                        rule[22],
							IcmpUnreachable:                 boolopt[rule[23]],
							DisableServerResponseInspection: boolopt[rule[24]],
							Group:                           rule[25],
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
							if strings.Contains(err.Error(), "Client.Timeout") {
								timeoutCount++
								timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, strings.TrimSpace(rule[2])))
							} else {
								log.Printf("Line %d - failed to create security rule: %s", i+1, err)
							}
						}
					}

					time.Sleep(delay * time.Millisecond)
				}
			}

			if t == "nat" {
				// #DeviceGroup,Location,Name,Type,Description,Tags,SourceZones,DestinationZone,ToInterface,Service,SourceAddresses,DestinationAddresses,
				// SatType,SatAddressType,SatTranslatedAddresses,SatInterface,SatIpAddress,SatFallbackType,SatFallbackTranslatedAddresses,SatFallbackInterface,
				// SatFallbackIpType,SatFallbackIpAddress,SatStaticTranslatedAddress,SatStaticBiDirectional,DatType,DatAddress,DatPort,DatDynamicDistribution,Disabled
				// 29 columns (0-28)
				rules, err := easycsv.Open(f)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				rc := len(rules)
				log.Printf("Importing/modifying %d NAT rules", rc)

				for i, rule := range rules {
					boolopt := map[string]bool{
						"TRUE":  true,
						"true":  true,
						"FALSE": false,
						"false": false,
					}

					dg := rule[0]
					l := rule[1]

					switch l {
					case "pre":
						l = util.PreRulebase
					case "post":
						l = util.PostRulebase
					default:
						l = util.PreRulebase
					}

					datport, _ := strconv.Atoi(rule[26])

					if rule[3] != "ipv4" {
						log.Printf("Line %d - only NAT type 'ipv4' is supported", i+1)
					}

					ruletype := rule[3]

					if len(ruletype) <= 0 {
						ruletype = "universal"
					}

					if len(strings.TrimSpace(rule[2])) > 63 {
						log.Printf("Line %d - failed to create %s: name is over the max 63 characters", i+1, strings.TrimSpace(rule[2]))
					} else {
						e := nat.Entry{
							Name:                           strings.TrimSpace(rule[2]),
							Type:                           ruletype,
							Description:                    rule[4],
							Tags:                           stringToSlice(rule[5]),
							SourceZones:                    stringToSlice(rule[6]),
							DestinationZone:                rule[7],
							ToInterface:                    rule[8],
							Service:                        rule[9],
							SourceAddresses:                stringToSlice(rule[10]),
							DestinationAddresses:           stringToSlice(rule[11]),
							SatType:                        rule[12],
							SatAddressType:                 rule[13],
							SatTranslatedAddresses:         stringToSlice(rule[14]),
							SatInterface:                   rule[15],
							SatIpAddress:                   rule[16],
							SatFallbackType:                rule[17],
							SatFallbackTranslatedAddresses: stringToSlice(rule[18]),
							SatFallbackInterface:           rule[19],
							SatFallbackIpType:              rule[20],
							SatFallbackIpAddress:           rule[21],
							SatStaticTranslatedAddress:     rule[22],
							SatStaticBiDirectional:         boolopt[rule[23]],
							DatType:                        rule[24],
							DatAddress:                     rule[25],
							DatPort:                        datport,
							DatDynamicDistribution:         rule[27],
							Disabled:                       boolopt[rule[28]],
						}

						err = c.Policies.Nat.Set(dg, l, e)
						if err != nil {
							if strings.Contains(err.Error(), "Client.Timeout") {
								timeoutCount++
								timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, strings.TrimSpace(rule[2])))
							} else {
								log.Printf("Line %d - failed to create NAT rule: %s", i+1, err)
							}
						}
					}

					time.Sleep(delay * time.Millisecond)
				}
			}

			if t == "pbf" {
				rules, err := easycsv.Open(f)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				rc := len(rules)
				log.Printf("Importing/modifying %d Policy-Based Forwarding rules", rc)

				for i, rule := range rules {
					boolopt := map[string]bool{
						"TRUE":  true,
						"true":  true,
						"FALSE": false,
						"false": false,
					}

					dg := rule[0]
					l := rule[1]

					switch l {
					case "pre":
						l = util.PreRulebase
					case "post":
						l = util.PostRulebase
					default:
						l = util.PreRulebase
					}

					if len(strings.TrimSpace(rule[2])) > 63 {
						log.Printf("Line %d - failed to create %s: name is over the max 63 characters", i+1, strings.TrimSpace(rule[2]))
					} else {
						e := pbf.Entry{
							Name:                               strings.TrimSpace(rule[2]),
							Description:                        rule[3],
							Tags:                               stringToSlice(rule[4]),
							FromType:                           rule[5],
							FromValues:                         stringToSlice(rule[6]),
							SourceAddresses:                    stringToSlice(rule[7]),
							SourceUsers:                        stringToSlice(rule[8]),
							NegateSource:                       boolopt[rule[9]],
							DestinationAddresses:               stringToSlice(rule[10]),
							NegateDestination:                  boolopt[rule[11]],
							Applications:                       stringToSlice(rule[12]),
							Services:                           stringToSlice(rule[13]),
							Schedule:                           rule[14],
							Disabled:                           boolopt[rule[15]],
							Action:                             rule[16],
							ForwardVsys:                        rule[17],
							ForwardEgressInterface:             rule[18],
							ForwardNextHopType:                 rule[19],
							ForwardNextHopValue:                rule[20],
							ForwardMonitorProfile:              rule[21],
							ForwardMonitorIpAddress:            rule[22],
							ForwardMonitorDisableIfUnreachable: boolopt[rule[23]],
							EnableEnforceSymmetricReturn:       boolopt[rule[24]],
							SymmetricReturnAddresses:           stringToSlice(rule[25]),
							ActiveActiveDeviceBinding:          rule[26],
							NegateTarget:                       boolopt[rule[27]],
							Uuid:                               rule[28],
						}

						err = c.Policies.PolicyBasedForwarding.Set(dg, l, e)
						if err != nil {
							if strings.Contains(err.Error(), "Client.Timeout") {
								timeoutCount++
								timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, strings.TrimSpace(rule[2])))
							} else {
								log.Printf("Line %d - failed to create Policy-Based Forwarding rule: %s", i+1, err)
							}
						}
					}

					time.Sleep(delay * time.Millisecond)
				}
			}

			if t == "decrypt" {
				// #DeviceGroup,Location,Name,Description,SourceZones,SourceAddresses,NegateSource,SourceUsers,DestinationZones,DestinationAddresses,
				// NegateDestination,Tags,Disabled,Services,UrlCategories,Action,DecryptionType,SslCertificate,DecryptionProfile,NegateTarget,
				// ForwardingProfile,GroupTag,SourceHips,DestinationHips,LogSuccessfulTlsHandshakes,LogFailedTlsHandshakes,LogSetting,SslCertificates
				// 28 columns (0-27)
				rules, err := easycsv.Open(f)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				rc := len(rules)
				log.Printf("Importing/modifying %d Policy-Based Forwarding rules", rc)

				for i, rule := range rules {
					boolopt := map[string]bool{
						"TRUE":  true,
						"true":  true,
						"FALSE": false,
						"false": false,
					}

					dg := rule[0]
					l := rule[1]

					switch l {
					case "pre":
						l = util.PreRulebase
					case "post":
						l = util.PostRulebase
					default:
						l = util.PreRulebase
					}

					if len(strings.TrimSpace(rule[0])) > 63 {
						log.Printf("Line %d - failed to create %s: name is over the max 63 characters", i+1, strings.TrimSpace(rule[2]))
					} else {
						e := decryption.Entry{
							Name:                       strings.TrimSpace(rule[2]),
							Description:                rule[3],
							SourceZones:                stringToSlice(rule[4]),
							SourceAddresses:            stringToSlice(rule[5]),
							NegateSource:               boolopt[rule[6]],
							SourceUsers:                userStringToSlice(rule[7]),
							DestinationZones:           stringToSlice(rule[8]),
							DestinationAddresses:       stringToSlice(rule[9]),
							NegateDestination:          boolopt[rule[10]],
							Tags:                       stringToSlice(rule[11]),
							Disabled:                   boolopt[rule[12]],
							Services:                   stringToSlice(rule[13]),
							UrlCategories:              stringToSlice(rule[14]),
							Action:                     rule[15],
							DecryptionType:             rule[16],
							SslCertificate:             rule[17],
							DecryptionProfile:          rule[18],
							NegateTarget:               boolopt[rule[19]],
							ForwardingProfile:          rule[20],
							GroupTag:                   rule[21],
							SourceHips:                 stringToSlice(rule[22]),
							DestinationHips:            stringToSlice(rule[23]),
							LogSuccessfulTlsHandshakes: boolopt[rule[24]],
							LogFailedTlsHandshakes:     boolopt[rule[25]],
							LogSetting:                 rule[26],
							SslCertificates:            stringToSlice(rule[27]),
						}

						err = c.Policies.Decryption.Set(dg, l, e)
						if err != nil {
							if strings.Contains(err.Error(), "Client.Timeout") {
								timeoutCount++
								timeoutData = append(timeoutData, fmt.Sprintf("%d:%s", i+1, strings.TrimSpace(rule[2])))
							} else {
								log.Printf("Line %d - failed to create Decryption rule: %s", i+1, err)
							}
						}
					}

					time.Sleep(delay * time.Millisecond)
				}
			}

			if timeoutCount > 0 {
				log.Printf("There were %d API timeout errors during import. Please verify the following have been imported/modified:\n\n", timeoutCount)
				for _, data := range timeoutData {
					info := strings.Split(data, ":")
					fmt.Printf("Line %s: Rule \"%s\"\n", info[0], info[1])
				}
			}
		}
	},
}

func init() {
	policyCmd.AddCommand(policyImportCmd)

	policyImportCmd.Flags().StringVarP(&user, "user", "u", "", "User to connect to the device as")
	policyImportCmd.Flags().StringVarP(&p, "delay", "p", "100", "Delay (in milliseconds) to pause between each API call")
	policyImportCmd.Flags().StringVarP(&device, "device", "d", "", "Device to connect to")
	policyImportCmd.Flags().StringVarP(&f, "file", "f", "", "Name of the CSV file to export to")
	// policyImportCmd.Flags().StringVarP(&dg, "devicegroup", "g", "shared", "Device Group name when importing to Panorama")
	// policyImportCmd.Flags().StringVarP(&v, "vsys", "v", "vsys1", "Vsys name when importing to a firewall")
	policyImportCmd.Flags().StringVarP(&t, "type", "t", "", "Type of policy to import - <security|nat|pbf|decrypt>")
	// policyImportCmd.Flags().StringVarP(&l, "location", "l", "pre", "Location of the rulebase - <pre|post>")
	policyImportCmd.MarkFlagRequired("user")
	// policyImportCmd.MarkFlagRequired("pass")
	policyImportCmd.MarkFlagRequired("device")
	policyImportCmd.MarkFlagRequired("file")
	policyImportCmd.MarkFlagRequired("type")
	// policyImportCmd.MarkFlagRequired("location")
}
