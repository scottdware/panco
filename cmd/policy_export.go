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
	"time"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/util"
	"github.com/Songmu/prompter"
	easycsv "github.com/scottdware/go-easycsv"
	"github.com/spf13/cobra"
)

// policyExportCmd represents the policy export command
var policyExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export a Security, NAT, Decryption or Policy-Based Forwarding policy",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
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

		f = strings.TrimSuffix(f, ".csv")

		switch c := con.(type) {
		case *pango.Firewall:
			// Security policy
			if t == "security" {
				getFwSecPol(c, f, hit)
			}

			if t == "all" {
				getFwSecPol(c, fmt.Sprintf("%s_Security", f), hit)
			}

			// NAT policy
			if t == "nat" {
				getFwNatPol(c, f, hit)
			}

			if t == "all" {
				getFwNatPol(c, fmt.Sprintf("%s_NAT", f), hit)
			}

			// Policy-Based Forwarding policy
			if t == "pbf" {
				getFwPbfPol(c, f, hit)
			}

			if t == "all" {
				getFwPbfPol(c, fmt.Sprintf("%s_PBF", f), hit)
			}

			// Decryption policy
			if t == "decrypt" {
				getFwDecryptPol(c, f, hit)
			}

			if t == "all" {
				getFwDecryptPol(c, fmt.Sprintf("%s_Decrypt", f), hit)
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

			// Security policy
			if t == "security" {
				getPanoSecPol(c, f, hit)
			}

			if t == "all" {
				getPanoSecPol(c, fmt.Sprintf("%s_Security", f), hit)
			}

			// NAT policy
			if t == "nat" {
				getPanoNatPol(c, f, hit)
			}

			if t == "all" {
				getPanoNatPol(c, fmt.Sprintf("%s_NAT", f), hit)
			}

			// Policy-Based Forwarding policy
			if t == "pbf" {
				getPanoPbfPol(c, f, hit)
			}

			if t == "all" {
				getPanoPbfPol(c, fmt.Sprintf("%s_PBF", f), hit)
			}

			// Decryption policy
			if t == "decrypt" {
				getPanoDecryptPol(c, f, hit)
			}

			if t == "all" {
				getPanoDecryptPol(c, fmt.Sprintf("%s_Decrypt", f), hit)
			}
		}
	},
}

func init() {
	policyCmd.AddCommand(policyExportCmd)

	policyExportCmd.Flags().StringVarP(&user, "user", "u", "", "User to connect to the device as")
	// policyExportCmd.Flags().StringVarP(&pass, "pass", "p", "", "Password for the user account specified")
	policyExportCmd.Flags().StringVarP(&device, "device", "d", "", "Device to connect to")
	policyExportCmd.Flags().StringVarP(&f, "file", "f", "PaloAltoPolicy", "Name of the CSV file you'd like to export to")
	policyExportCmd.Flags().StringVarP(&dg, "devicegroup", "g", "shared", "Device Group name when exporting from Panorama")
	policyExportCmd.Flags().StringVarP(&v, "vsys", "v", "vsys1", "Vsys name when exporting from a firewall")
	policyExportCmd.Flags().StringVarP(&t, "type", "t", "", "Type of policy to export - <security|nat|pbf|decrypt|all>")
	policyExportCmd.Flags().StringVarP(&l, "location", "l", "pre", "Location of the rulebase - <pre|post>")
	policyExportCmd.Flags().StringVarP(&onlyrules, "rules", "r", "", "[OPTIONAL] Only export these specific rules in referenced text file")
	policyExportCmd.MarkFlagRequired("user")
	// policyExportCmd.MarkFlagRequired("pass")
	policyExportCmd.MarkFlagRequired("device")
	policyExportCmd.MarkFlagRequired("file")
	policyExportCmd.MarkFlagRequired("type")
	// policyExportCmd.MarkFlagRequired("location")
}

// getFwSecPol is used to export the Security policy on a firewall
func getFwSecPol(c *pango.Firewall, file string, hitcount bool) {
	rules, err := c.Policies.Security.GetList(v)
	if err != nil {
		log.Printf("Failed to retrieve the list of Security rules: %s", err)
		return
	}

	rc := len(rules)
	if rc <= 0 {
		log.Printf("There are 0 Security rules for '%s' - no policy was exported", v)
		return
	}

	secfile := fmt.Sprintf("%s.csv", file)

	cfh, err := easycsv.NewCSV(secfile)
	if err != nil {
		log.Printf("CSV file error - %s", err)
		return
	}

	cfh.Write("#Name,Type,Description,Tags,SourceZones,SourceAddresses,NegateSource,SourceUsers,HipProfiles,")
	cfh.Write("DestinationZones,DestinationAddresses,NegateDestination,Applications,Services,Categories,")
	cfh.Write("Action,LogSetting,LogStart,LogEnd,Disabled,Schedule,IcmpUnreachable,DisableServerResponseInspection,")
	cfh.Write("Group,Virus,Spyware,Vulnerability,UrlFiltering,FileBlocking,WildFireAnalysis,DataFiltering\n")

	if len(onlyrules) > 0 {
		inclrules, err := txtToSlice(onlyrules)
		if err != nil {
			log.Printf("%s", err)
		}

		log.Printf("Exporting %d Security rules", len(inclrules))

		for _, rule := range rules {
			for _, orule := range inclrules {
				if rule == orule {
					var rtype string
					r, err := c.Policies.Security.Get(v, orule)
					if err != nil {
						log.Printf("Failed to retrieve Security rule data: %s", err)
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

					cfh.Write(fmt.Sprintf("%s,%s,\"%s\",\"%s\",\"%s\",\"%s\",%t,\"%s\",\"%s\",", r.Name, rtype, formatDesc(r.Description), sliceToString(r.Tags), sliceToString(r.SourceZones),
						sliceToString(r.SourceAddresses), r.NegateSource, userSliceToString(r.SourceUsers), sliceToString(r.HipProfiles)))
					cfh.Write(fmt.Sprintf("\"%s\",\"%s\",%t,\"%s\",\"%s\",\"%s\",", sliceToString(r.DestinationZones), sliceToString(r.DestinationAddresses), r.NegateDestination,
						sliceToString(r.Applications), sliceToString(r.Services), sliceToString(r.Categories)))
					cfh.Write(fmt.Sprintf("%s,%s,%t,%t,%t,%s,%t,%t,", r.Action, r.LogSetting, r.LogStart, r.LogEnd, r.Disabled, r.Schedule,
						r.IcmpUnreachable, r.DisableServerResponseInspection))
					cfh.Write(fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s\n", r.Group, r.Virus, r.Spyware,
						r.Vulnerability, r.UrlFiltering, r.FileBlocking, r.WildFireAnalysis, r.DataFiltering))

					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	} else {
		log.Printf("Exporting %d Security rules", rc)

		for _, rule := range rules {
			var rtype string
			r, err := c.Policies.Security.Get(v, rule)
			if err != nil {
				log.Printf("Failed to retrieve Security rule data: %s", err)
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

			cfh.Write(fmt.Sprintf("%s,%s,\"%s\",\"%s\",\"%s\",\"%s\",%t,\"%s\",\"%s\",", r.Name, rtype, formatDesc(r.Description), sliceToString(r.Tags), sliceToString(r.SourceZones),
				sliceToString(r.SourceAddresses), r.NegateSource, userSliceToString(r.SourceUsers), sliceToString(r.HipProfiles)))
			cfh.Write(fmt.Sprintf("\"%s\",\"%s\",%t,\"%s\",\"%s\",\"%s\",", sliceToString(r.DestinationZones), sliceToString(r.DestinationAddresses), r.NegateDestination,
				sliceToString(r.Applications), sliceToString(r.Services), sliceToString(r.Categories)))
			cfh.Write(fmt.Sprintf("%s,%s,%t,%t,%t,%s,%t,%t,", r.Action, r.LogSetting, r.LogStart, r.LogEnd, r.Disabled, r.Schedule,
				r.IcmpUnreachable, r.DisableServerResponseInspection))
			cfh.Write(fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s\n", r.Group, r.Virus, r.Spyware,
				r.Vulnerability, r.UrlFiltering, r.FileBlocking, r.WildFireAnalysis, r.DataFiltering))

			time.Sleep(100 * time.Millisecond)
		}
	}

	cfh.End()
}

// getFwNatPol is used to export the NAT policy on a firewall
func getFwNatPol(c *pango.Firewall, file string, hitcount bool) {
	rules, err := c.Policies.Nat.GetList(v)
	if err != nil {
		log.Printf("Failed to retrieve the list of NAT rules: %s", err)
		return
	}

	rc := len(rules)
	if rc <= 0 {
		log.Printf("There are 0 NAT rules for '%s' - no policy was exported", v)
		return
	}

	natfile := fmt.Sprintf("%s.csv", file)

	cfh, err := easycsv.NewCSV(natfile)
	if err != nil {
		log.Printf("CSV file error - %s", err)
		return
	}

	cfh.Write("#Name,Type,Description,Tags,SourceZones,DestinationZone,ToInterface,Service,SourceAddresses,DestinationAddresses,")
	cfh.Write("SatType,SatAddressType,SatTranslatedAddresses,SatInterface,SatIpAddress,SatFallbackType,SatFallbackTranslatedAddresses,")
	cfh.Write("SatFallbackInterface,SatFallbackIpType,SatFallbackIpAddress,SatStaticTranslatedAddress,SatStaticBiDirectional,DatType,")
	cfh.Write("DatAddress,DatPort,DatDynamicDistribution,Disabled\n")

	if len(onlyrules) > 0 {
		inclrules, err := txtToSlice(onlyrules)
		if err != nil {
			log.Printf("%s", err)
		}

		log.Printf("Exporting %d NAT rules", len(inclrules))

		for _, rule := range rules {
			for _, orule := range inclrules {
				if rule == orule {
					var toint string
					r, err := c.Policies.Nat.Get(v, orule)
					if err != nil {
						log.Printf("Failed to retrieve NAT rule data: %s", err)
					}

					toint = r.ToInterface
					if len(r.ToInterface) <= 0 {
						toint = "any"
					}

					cfh.Write(fmt.Sprintf("%s,%s,\"%s\",\"%s\",\"%s\",%s,%s,%s,\"%s\",\"%s\",", r.Name, "ipv4", formatDesc(r.Description), sliceToString(r.Tags), sliceToString(r.SourceZones),
						r.DestinationZone, toint, r.Service, sliceToString(r.SourceAddresses), sliceToString(r.DestinationAddresses)))
					cfh.Write(fmt.Sprintf("%s,%s,\"%s\",%s,%s,%s,\"%s\",%s,", r.SatType, r.SatAddressType, sliceToString(r.SatTranslatedAddresses), r.SatInterface,
						r.SatIpAddress, r.SatFallbackType, sliceToString(r.SatFallbackTranslatedAddresses), r.SatFallbackInterface))
					cfh.Write(fmt.Sprintf("%s,%s,%s,%t,%s,%s,%d,%s,%t\n", r.SatFallbackIpType, r.SatFallbackIpAddress, r.SatStaticTranslatedAddress,
						r.SatStaticBiDirectional, r.DatType, r.DatAddress, r.DatPort, r.DatDynamicDistribution, r.Disabled))

					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	} else {
		log.Printf("Exporting %d NAT rules", rc)

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

			cfh.Write(fmt.Sprintf("%s,%s,\"%s\",\"%s\",\"%s\",%s,%s,%s,\"%s\",\"%s\",", r.Name, "ipv4", formatDesc(r.Description), sliceToString(r.Tags), sliceToString(r.SourceZones),
				r.DestinationZone, toint, r.Service, sliceToString(r.SourceAddresses), sliceToString(r.DestinationAddresses)))
			cfh.Write(fmt.Sprintf("%s,%s,\"%s\",%s,%s,%s,\"%s\",%s,", r.SatType, r.SatAddressType, sliceToString(r.SatTranslatedAddresses), r.SatInterface,
				r.SatIpAddress, r.SatFallbackType, sliceToString(r.SatFallbackTranslatedAddresses), r.SatFallbackInterface))
			cfh.Write(fmt.Sprintf("%s,%s,%s,%t,%s,%s,%d,%s,%t\n", r.SatFallbackIpType, r.SatFallbackIpAddress, r.SatStaticTranslatedAddress,
				r.SatStaticBiDirectional, r.DatType, r.DatAddress, r.DatPort, r.DatDynamicDistribution, r.Disabled))

			time.Sleep(100 * time.Millisecond)
		}
	}

	cfh.End()
}

// getFwPbfPol is used to export the Policy-Based Forwarding policy on a firewall
func getFwPbfPol(c *pango.Firewall, file string, hitcount bool) {
	rules, err := c.Policies.PolicyBasedForwarding.GetList(v)
	if err != nil {
		log.Printf("Failed to retrieve the list of Policy-Based Forwarding rules: %s", err)
		return
	}

	rc := len(rules)
	if rc <= 0 {
		log.Printf("There are 0 Policy-Based Forwarding rules for '%s' - no policy was exported", v)
		return
	}

	pbffile := fmt.Sprintf("%s.csv", file)

	cfh, err := easycsv.NewCSV(pbffile)
	if err != nil {
		log.Printf("CSV file error - %s", err)
		return
	}

	cfh.Write("#Name,Description,Tags,FromType,FromValues,SourceAddresses,SourceUsers,NegateSource,DestinationAddresses,")
	cfh.Write("NegateDestination,Applications,Services,Schedule,Disabled,Action,ForwardVsys,ForwardEgressInterface,")
	cfh.Write("ForwardNextHopType,ForwardNextHopValue,ForwardMonitorProfile,ForwardMonitorIpAddress,ForwardMonitorDisableIfUnreachable,")
	cfh.Write("EnableEnforceSymmetricReturn,SymmetricReturnAddresses,ActiveActiveDeviceBinding,NegateTarget,Uuid\n")

	if len(onlyrules) > 0 {
		inclrules, err := txtToSlice(onlyrules)
		if err != nil {
			log.Printf("%s", err)
		}

		log.Printf("Exporting %d Policy-Based Forwarding rules", len(inclrules))

		for _, rule := range rules {
			for _, orule := range inclrules {
				if rule == orule {
					r, err := c.Policies.PolicyBasedForwarding.Get(v, rule)
					if err != nil {
						log.Printf("Failed to retrieve Policy-Based Forwarding rule data: %s", err)
					}

					cfh.Write(fmt.Sprintf("%s,\"%s\",\"%s\",%s,\"%s\",\"%s\",\"%s\",%t,\"%s\",", r.Name, formatDesc(r.Description), sliceToString(r.Tags), r.FromType,
						sliceToString(r.FromValues), sliceToString(r.SourceAddresses), userSliceToString(r.SourceUsers), r.NegateSource, sliceToString(r.DestinationAddresses)))
					cfh.Write(fmt.Sprintf("%t,\"%s\",\"%s\",%s,%t,%s,%s,%s,", r.NegateDestination, sliceToString(r.Applications), sliceToString(r.Services), r.Schedule,
						r.Disabled, r.Action, r.ForwardVsys, r.ForwardEgressInterface))
					cfh.Write(fmt.Sprintf("%s,%s,%s,%s,%t,", r.ForwardNextHopType, r.ForwardNextHopValue, r.ForwardMonitorProfile, r.ForwardMonitorIpAddress,
						r.ForwardMonitorDisableIfUnreachable))
					cfh.Write(fmt.Sprintf("%t,\"%s\",%s,%t,%s\n", r.EnableEnforceSymmetricReturn, sliceToString(r.SymmetricReturnAddresses), r.ActiveActiveDeviceBinding,
						r.NegateTarget, r.Uuid))

					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	} else {
		log.Printf("Exporting %d Policy-Based Forwarding rules", rc)

		for _, rule := range rules {
			r, err := c.Policies.PolicyBasedForwarding.Get(v, rule)
			if err != nil {
				log.Printf("Failed to retrieve Policy-Based Forwarding rule data: %s", err)
			}

			cfh.Write(fmt.Sprintf("%s,\"%s\",\"%s\",%s,\"%s\",\"%s\",\"%s\",%t,\"%s\",", r.Name, formatDesc(r.Description), sliceToString(r.Tags), r.FromType,
				sliceToString(r.FromValues), sliceToString(r.SourceAddresses), userSliceToString(r.SourceUsers), r.NegateSource, sliceToString(r.DestinationAddresses)))
			cfh.Write(fmt.Sprintf("%t,\"%s\",\"%s\",%s,%t,%s,%s,%s,", r.NegateDestination, sliceToString(r.Applications), sliceToString(r.Services), r.Schedule,
				r.Disabled, r.Action, r.ForwardVsys, r.ForwardEgressInterface))
			cfh.Write(fmt.Sprintf("%s,%s,%s,%s,%t,", r.ForwardNextHopType, r.ForwardNextHopValue, r.ForwardMonitorProfile, r.ForwardMonitorIpAddress,
				r.ForwardMonitorDisableIfUnreachable))
			cfh.Write(fmt.Sprintf("%t,\"%s\",%s,%t,%s\n", r.EnableEnforceSymmetricReturn, sliceToString(r.SymmetricReturnAddresses), r.ActiveActiveDeviceBinding,
				r.NegateTarget, r.Uuid))

			time.Sleep(100 * time.Millisecond)
		}
	}

	cfh.End()
}

// getFwDecryptPol is used to export the Decryption policy on a firewall
func getFwDecryptPol(c *pango.Firewall, file string, hitcount bool) {
	rules, err := c.Policies.Decryption.GetList(v)
	if err != nil {
		log.Printf("Failed to retrieve the list of Decryption rules: %s", err)
		return
	}

	rc := len(rules)
	if rc <= 0 {
		log.Printf("There are 0 Decryption rules for '%s' - no policy was exported", v)
		return
	}

	decryptfile := fmt.Sprintf("%s.csv", file)

	cfh, err := easycsv.NewCSV(decryptfile)
	if err != nil {
		log.Printf("CSV file error - %s", err)
		return
	}

	cfh.Write("#Name,Description,SourceZones,SourceAddresses,NegateSource,SourceUsers,DestinationZones,DestinationAddresses,NegateDestination,")
	cfh.Write("Tags,Disabled,Services,UrlCategories,Action,DecryptionType,SslCertificate,DecryptionProfile,NegateTarget,")
	cfh.Write("ForwardingProfile,Uuid,GroupTag,SourceHips,DestinationHips,LogSuccessfulTlsHandshakes,LogFailedTlsHandshakes,LogSetting,SslCertificates\n")

	if len(onlyrules) > 0 {
		inclrules, err := txtToSlice(onlyrules)
		if err != nil {
			log.Printf("%s", err)
		}

		log.Printf("Exporting %d Decryption rules", len(inclrules))

		for _, rule := range rules {
			for _, orule := range inclrules {
				if rule == orule {
					r, err := c.Policies.Decryption.Get(v, rule)
					if err != nil {
						log.Printf("Failed to retrieve Decryption rule data: %s", err)
					}

					cfh.Write(fmt.Sprintf("%s,%s,\"%s\",\"%s\",%t,\"%s\",\"%s\",\"%s\",%t,", r.Name, formatDesc(r.Description),
						sliceToString(r.SourceZones), sliceToString(r.SourceAddresses), r.NegateSource, userSliceToString(r.SourceUsers), sliceToString(r.DestinationZones), sliceToString(r.DestinationAddresses), r.NegateDestination))
					cfh.Write(fmt.Sprintf("\"%s\",%t,\"%s\",\"%s\",%s,%s,%s,%s,%t,", sliceToString(r.Tags), r.Disabled,
						sliceToString(r.Services), sliceToString(r.UrlCategories), r.Action, r.DecryptionType, r.SslCertificate, r.DecryptionProfile, r.NegateTarget))
					cfh.Write(fmt.Sprintf("%s,%s,%s,\"%s\",\"%s\",%t,%t,%s,\"%s\"\n", r.ForwardingProfile, r.Uuid, r.GroupTag,
						sliceToString(r.SourceHips), sliceToString(r.DestinationHips), r.LogSuccessfulTlsHandshakes, r.LogFailedTlsHandshakes, r.LogSetting, sliceToString(r.SslCertificates)))

					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	} else {
		log.Printf("Exporting %d Decryption rules", rc)

		for _, rule := range rules {
			r, err := c.Policies.Decryption.Get(v, rule)
			if err != nil {
				log.Printf("Failed to retrieve Decryption rule data: %s", err)
			}

			cfh.Write(fmt.Sprintf("%s,%s,\"%s\",\"%s\",%t,\"%s\",\"%s\",\"%s\",%t,", r.Name, formatDesc(r.Description),
				sliceToString(r.SourceZones), sliceToString(r.SourceAddresses), r.NegateSource, userSliceToString(r.SourceUsers), sliceToString(r.DestinationZones), sliceToString(r.DestinationAddresses), r.NegateDestination))
			cfh.Write(fmt.Sprintf("\"%s\",%t,\"%s\",\"%s\",%s,%s,%s,%s,%t,", sliceToString(r.Tags), r.Disabled,
				sliceToString(r.Services), sliceToString(r.UrlCategories), r.Action, r.DecryptionType, r.SslCertificate, r.DecryptionProfile, r.NegateTarget))
			cfh.Write(fmt.Sprintf("%s,%s,%s,\"%s\",\"%s\",%t,%t,%s,\"%s\"\n", r.ForwardingProfile, r.Uuid, r.GroupTag,
				sliceToString(r.SourceHips), sliceToString(r.DestinationHips), r.LogSuccessfulTlsHandshakes, r.LogFailedTlsHandshakes, r.LogSetting, sliceToString(r.SslCertificates)))

			time.Sleep(100 * time.Millisecond)
		}
	}

	cfh.End()
}

// getPanoSecPol is used to export the Security policy from Panorama
func getPanoSecPol(c *pango.Panorama, file string, hitcount bool) {
	rules, err := c.Policies.Security.GetList(dg, l)
	if err != nil {
		log.Printf("Failed to retrieve the list of Security rules: %s", err)
		return
	}

	rc := len(rules)
	if rc <= 0 {
		log.Printf("There are 0 Security rules for the '%s' device group - no policy was exported", dg)
		return
	}

	secfile := fmt.Sprintf("%s.csv", file)

	cfh, err := easycsv.NewCSV(secfile)
	if err != nil {
		log.Printf("CSV file error - %s", err)
		return
	}

	cfh.Write("#Name,Type,Description,Tags,SourceZones,SourceAddresses,NegateSource,SourceUsers,HipProfiles,")
	cfh.Write("DestinationZones,DestinationAddresses,NegateDestination,Applications,Services,Categories,")
	cfh.Write("Action,LogSetting,LogStart,LogEnd,Disabled,Schedule,IcmpUnreachable,DisableServerResponseInspection,")
	cfh.Write("Group,Virus,Spyware,Vulnerability,UrlFiltering,FileBlocking,WildFireAnalysis,DataFiltering\n")

	if len(onlyrules) > 0 {
		inclrules, err := txtToSlice(onlyrules)
		if err != nil {
			log.Printf("%s", err)
		}

		log.Printf("Exporting %d Security rules", len(inclrules))

		for _, rule := range rules {
			for _, orule := range inclrules {
				if rule == orule {
					var rtype string
					r, err := c.Policies.Security.Get(dg, l, orule)
					if err != nil {
						log.Printf("Failed to retrieve Security rule data: %s", err)
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

					cfh.Write(fmt.Sprintf("%s,%s,\"%s\",\"%s\",\"%s\",\"%s\",%t,\"%s\",\"%s\",", r.Name, rtype, formatDesc(r.Description), sliceToString(r.Tags), sliceToString(r.SourceZones),
						sliceToString(r.SourceAddresses), r.NegateSource, userSliceToString(r.SourceUsers), sliceToString(r.HipProfiles)))
					cfh.Write(fmt.Sprintf("\"%s\",\"%s\",%t,\"%s\",\"%s\",\"%s\",", sliceToString(r.DestinationZones), sliceToString(r.DestinationAddresses), r.NegateDestination,
						sliceToString(r.Applications), sliceToString(r.Services), sliceToString(r.Categories)))
					cfh.Write(fmt.Sprintf("%s,%s,%t,%t,%t,%s,%t,%t,", r.Action, r.LogSetting, r.LogStart, r.LogEnd, r.Disabled, r.Schedule,
						r.IcmpUnreachable, r.DisableServerResponseInspection))
					cfh.Write(fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s\n", r.Group, r.Virus, r.Spyware,
						r.Vulnerability, r.UrlFiltering, r.FileBlocking, r.WildFireAnalysis, r.DataFiltering))

					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	} else {
		log.Printf("Exporting %d Security rules", rc)

		for _, rule := range rules {
			var rtype string
			r, err := c.Policies.Security.Get(dg, l, rule)
			if err != nil {
				log.Printf("Failed to retrieve Security rule data: %s", err)
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

			cfh.Write(fmt.Sprintf("%s,%s,\"%s\",\"%s\",\"%s\",\"%s\",%t,\"%s\",\"%s\",", r.Name, rtype, formatDesc(r.Description), sliceToString(r.Tags), sliceToString(r.SourceZones),
				sliceToString(r.SourceAddresses), r.NegateSource, userSliceToString(r.SourceUsers), sliceToString(r.HipProfiles)))
			cfh.Write(fmt.Sprintf("\"%s\",\"%s\",%t,\"%s\",\"%s\",\"%s\",", sliceToString(r.DestinationZones), sliceToString(r.DestinationAddresses), r.NegateDestination,
				sliceToString(r.Applications), sliceToString(r.Services), sliceToString(r.Categories)))
			cfh.Write(fmt.Sprintf("%s,%s,%t,%t,%t,%s,%t,%t,", r.Action, r.LogSetting, r.LogStart, r.LogEnd, r.Disabled, r.Schedule,
				r.IcmpUnreachable, r.DisableServerResponseInspection))
			cfh.Write(fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s\n", r.Group, r.Virus, r.Spyware,
				r.Vulnerability, r.UrlFiltering, r.FileBlocking, r.WildFireAnalysis, r.DataFiltering))

			time.Sleep(100 * time.Millisecond)
		}
	}

	cfh.End()
}

// getPanoNatPol is used to export the NAT policy from Panorama
func getPanoNatPol(c *pango.Panorama, file string, hitcount bool) {
	rules, err := c.Policies.Nat.GetList(dg, l)
	if err != nil {
		log.Printf("Failed to retrieve the list of NAT rules: %s", err)
		return
	}

	rc := len(rules)
	if rc <= 0 {
		log.Printf("There are 0 NAT rules for the '%s' device group - no policy was exported", dg)
		return
	}

	natfile := fmt.Sprintf("%s.csv", file)

	cfh, err := easycsv.NewCSV(natfile)
	if err != nil {
		log.Printf("CSV file error - %s", err)
		return
	}

	cfh.Write("#Name,Type,Description,Tags,SourceZones,DestinationZone,ToInterface,Service,SourceAddresses,DestinationAddresses,")
	cfh.Write("SatType,SatAddressType,SatTranslatedAddresses,SatInterface,SatIpAddress,SatFallbackType,SatFallbackTranslatedAddresses,")
	cfh.Write("SatFallbackInterface,SatFallbackIpType,SatFallbackIpAddress,SatStaticTranslatedAddress,SatStaticBiDirectional,DatType,")
	cfh.Write("DatAddress,DatPort,DatDynamicDistribution,Disabled\n")

	if len(onlyrules) > 0 {
		inclrules, err := txtToSlice(onlyrules)
		if err != nil {
			log.Printf("%s", err)
		}

		log.Printf("Exporting %d NAT rules", len(inclrules))

		for _, rule := range rules {
			for _, orule := range inclrules {
				if rule == orule {
					var toint string
					r, err := c.Policies.Nat.Get(dg, l, orule)
					if err != nil {
						log.Printf("Failed to retrieve NAT rule data: %s", err)
					}

					toint = r.ToInterface
					if len(r.ToInterface) <= 0 {
						toint = "any"
					}

					cfh.Write(fmt.Sprintf("%s,%s,\"%s\",\"%s\",\"%s\",%s,%s,%s,\"%s\",\"%s\",", r.Name, "ipv4", formatDesc(r.Description), sliceToString(r.Tags), sliceToString(r.SourceZones),
						r.DestinationZone, toint, r.Service, sliceToString(r.SourceAddresses), sliceToString(r.DestinationAddresses)))
					cfh.Write(fmt.Sprintf("%s,%s,\"%s\",%s,%s,%s,\"%s\",%s,", r.SatType, r.SatAddressType, sliceToString(r.SatTranslatedAddresses), r.SatInterface,
						r.SatIpAddress, r.SatFallbackType, sliceToString(r.SatFallbackTranslatedAddresses), r.SatFallbackInterface))
					cfh.Write(fmt.Sprintf("%s,%s,%s,%t,%s,%s,%d,%s,%t\n", r.SatFallbackIpType, r.SatFallbackIpAddress, r.SatStaticTranslatedAddress,
						r.SatStaticBiDirectional, r.DatType, r.DatAddress, r.DatPort, r.DatDynamicDistribution, r.Disabled))

					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	} else {
		log.Printf("Exporting %d NAT rules", rc)

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

			cfh.Write(fmt.Sprintf("%s,%s,\"%s\",\"%s\",\"%s\",%s,%s,%s,\"%s\",\"%s\",", r.Name, "ipv4", formatDesc(r.Description), sliceToString(r.Tags), sliceToString(r.SourceZones),
				r.DestinationZone, toint, r.Service, sliceToString(r.SourceAddresses), sliceToString(r.DestinationAddresses)))
			cfh.Write(fmt.Sprintf("%s,%s,\"%s\",%s,%s,%s,\"%s\",%s,", r.SatType, r.SatAddressType, sliceToString(r.SatTranslatedAddresses), r.SatInterface,
				r.SatIpAddress, r.SatFallbackType, sliceToString(r.SatFallbackTranslatedAddresses), r.SatFallbackInterface))
			cfh.Write(fmt.Sprintf("%s,%s,%s,%t,%s,%s,%d,%s,%t\n", r.SatFallbackIpType, r.SatFallbackIpAddress, r.SatStaticTranslatedAddress,
				r.SatStaticBiDirectional, r.DatType, r.DatAddress, r.DatPort, r.DatDynamicDistribution, r.Disabled))

			time.Sleep(100 * time.Millisecond)
		}
	}

	cfh.End()
}

// getPanoPbfPol is used to export the Policy-Based Forwarding policy from Panorama
func getPanoPbfPol(c *pango.Panorama, file string, hitcount bool) {
	rules, err := c.Policies.PolicyBasedForwarding.GetList(dg, l)
	if err != nil {
		log.Printf("Failed to retrieve the list of Policy-Based Forwarding rules: %s", err)
		return
	}

	rc := len(rules)
	if rc <= 0 {
		log.Printf("There are 0 Policy-Based Forwarding rules for '%s' - no policy was exported", dg)
		return
	}

	pbffile := fmt.Sprintf("%s.csv", file)

	cfh, err := easycsv.NewCSV(pbffile)
	if err != nil {
		log.Printf("CSV file error - %s", err)
		return
	}

	cfh.Write("#Name,Description,Tags,FromType,FromValues,SourceAddresses,SourceUsers,NegateSource,DestinationAddresses,")
	cfh.Write("NegateDestination,Applications,Services,Schedule,Disabled,Action,ForwardVsys,ForwardEgressInterface,")
	cfh.Write("ForwardNextHopType,ForwardNextHopValue,ForwardMonitorProfile,ForwardMonitorIpAddress,ForwardMonitorDisableIfUnreachable,")
	cfh.Write("EnableEnforceSymmetricReturn,SymmetricReturnAddresses,ActiveActiveDeviceBinding,NegateTarget,Uuid\n")

	if len(onlyrules) > 0 {
		inclrules, err := txtToSlice(onlyrules)
		if err != nil {
			log.Printf("%s", err)
		}

		log.Printf("Exporting %d Policy-Based Forwarding rules", len(inclrules))

		for _, rule := range rules {
			for _, orule := range inclrules {
				if rule == orule {
					r, err := c.Policies.PolicyBasedForwarding.Get(dg, l, orule)
					if err != nil {
						log.Printf("Failed to retrieve Policy-Based Forwarding rule data: %s", err)
					}

					cfh.Write(fmt.Sprintf("%s,\"%s\",\"%s\",%s,\"%s\",\"%s\",\"%s\",%t,\"%s\",", r.Name, formatDesc(r.Description), sliceToString(r.Tags), r.FromType,
						sliceToString(r.FromValues), sliceToString(r.SourceAddresses), userSliceToString(r.SourceUsers), r.NegateSource, sliceToString(r.DestinationAddresses)))
					cfh.Write(fmt.Sprintf("%t,\"%s\",\"%s\",%s,%t,%s,%s,%s,", r.NegateDestination, sliceToString(r.Applications), sliceToString(r.Services), r.Schedule,
						r.Disabled, r.Action, r.ForwardVsys, r.ForwardEgressInterface))
					cfh.Write(fmt.Sprintf("%s,%s,%s,%s,%t,", r.ForwardNextHopType, r.ForwardNextHopValue, r.ForwardMonitorProfile, r.ForwardMonitorIpAddress,
						r.ForwardMonitorDisableIfUnreachable))
					cfh.Write(fmt.Sprintf("%t,\"%s\",%s,%t,%s\n", r.EnableEnforceSymmetricReturn, sliceToString(r.SymmetricReturnAddresses), r.ActiveActiveDeviceBinding,
						r.NegateTarget, r.Uuid))

					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	} else {
		log.Printf("Exporting %d Policy-Based Forwarding rules", rc)

		for _, rule := range rules {
			r, err := c.Policies.PolicyBasedForwarding.Get(dg, l, rule)
			if err != nil {
				log.Printf("Failed to retrieve Policy-Based Forwarding rule data: %s", err)
			}

			cfh.Write(fmt.Sprintf("%s,\"%s\",\"%s\",%s,\"%s\",\"%s\",\"%s\",%t,\"%s\",", r.Name, formatDesc(r.Description), sliceToString(r.Tags), r.FromType,
				sliceToString(r.FromValues), sliceToString(r.SourceAddresses), userSliceToString(r.SourceUsers), r.NegateSource, sliceToString(r.DestinationAddresses)))
			cfh.Write(fmt.Sprintf("%t,\"%s\",\"%s\",%s,%t,%s,%s,%s,", r.NegateDestination, sliceToString(r.Applications), sliceToString(r.Services), r.Schedule,
				r.Disabled, r.Action, r.ForwardVsys, r.ForwardEgressInterface))
			cfh.Write(fmt.Sprintf("%s,%s,%s,%s,%t,", r.ForwardNextHopType, r.ForwardNextHopValue, r.ForwardMonitorProfile, r.ForwardMonitorIpAddress,
				r.ForwardMonitorDisableIfUnreachable))
			cfh.Write(fmt.Sprintf("%t,\"%s\",%s,%t,%s\n", r.EnableEnforceSymmetricReturn, sliceToString(r.SymmetricReturnAddresses), r.ActiveActiveDeviceBinding,
				r.NegateTarget, r.Uuid))

			time.Sleep(100 * time.Millisecond)
		}
	}

	cfh.End()
}

// getPanoDecryptPol is used to export the Decryption policy from Panorama
func getPanoDecryptPol(c *pango.Panorama, file string, hitcount bool) {
	rules, err := c.Policies.Decryption.GetList(dg, l)
	if err != nil {
		log.Printf("Failed to retrieve the list of Decryption rules: %s", err)
		return
	}

	rc := len(rules)
	if rc <= 0 {
		log.Printf("There are 0 Decryption rules for '%s' - no policy was exported", dg)
		return
	}

	decryptfile := fmt.Sprintf("%s.csv", file)

	cfh, err := easycsv.NewCSV(decryptfile)
	if err != nil {
		log.Printf("CSV file error - %s", err)
		return
	}

	cfh.Write("#Name,Description,SourceZones,SourceAddresses,NegateSource,SourceUsers,DestinationZones,DestinationAddresses,NegateDestination,")
	cfh.Write("Tags,Disabled,Services,UrlCategories,Action,DecryptionType,SslCertificate,DecryptionProfile,NegateTarget,")
	cfh.Write("ForwardingProfile,Uuid,GroupTag,SourceHips,DestinationHips,LogSuccessfulTlsHandshakes,LogFailedTlsHandshakes,LogSetting,SslCertificates\n")

	if len(onlyrules) > 0 {
		inclrules, err := txtToSlice(onlyrules)
		if err != nil {
			log.Printf("%s", err)
		}

		log.Printf("Exporting %d Decryption rules", len(inclrules))

		for _, rule := range rules {
			for _, orule := range inclrules {
				if rule == orule {
					r, err := c.Policies.Decryption.Get(dg, l, orule)
					if err != nil {
						log.Printf("Failed to retrieve Decryption rule data: %s", err)
					}

					cfh.Write(fmt.Sprintf("%s,%s,\"%s\",\"%s\",%t,\"%s\",\"%s\",\"%s\",%t,", r.Name, formatDesc(r.Description),
						sliceToString(r.SourceZones), sliceToString(r.SourceAddresses), r.NegateSource, userSliceToString(r.SourceUsers), sliceToString(r.DestinationZones), sliceToString(r.DestinationAddresses), r.NegateDestination))
					cfh.Write(fmt.Sprintf("\"%s\",%t,\"%s\",\"%s\",%s,%s,%s,%s,%t,", sliceToString(r.Tags), r.Disabled,
						sliceToString(r.Services), sliceToString(r.UrlCategories), r.Action, r.DecryptionType, r.SslCertificate, r.DecryptionProfile, r.NegateTarget))
					cfh.Write(fmt.Sprintf("%s,%s,%s,\"%s\",\"%s\",%t,%t,%s,\"%s\"\n", r.ForwardingProfile, r.Uuid, r.GroupTag,
						sliceToString(r.SourceHips), sliceToString(r.DestinationHips), r.LogSuccessfulTlsHandshakes, r.LogFailedTlsHandshakes, r.LogSetting, sliceToString(r.SslCertificates)))

					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	} else {
		log.Printf("Exporting %d Decryption rules", rc)

		for _, rule := range rules {
			r, err := c.Policies.Decryption.Get(dg, l, rule)
			if err != nil {
				log.Printf("Failed to retrieve Decryption rule data: %s", err)
			}

			cfh.Write(fmt.Sprintf("%s,%s,\"%s\",\"%s\",%t,\"%s\",\"%s\",\"%s\",%t,", r.Name, formatDesc(r.Description),
				sliceToString(r.SourceZones), sliceToString(r.SourceAddresses), r.NegateSource, userSliceToString(r.SourceUsers), sliceToString(r.DestinationZones), sliceToString(r.DestinationAddresses), r.NegateDestination))
			cfh.Write(fmt.Sprintf("\"%s\",%t,\"%s\",\"%s\",%s,%s,%s,%s,%t,", sliceToString(r.Tags), r.Disabled,
				sliceToString(r.Services), sliceToString(r.UrlCategories), r.Action, r.DecryptionType, r.SslCertificate, r.DecryptionProfile, r.NegateTarget))
			cfh.Write(fmt.Sprintf("%s,%s,%s,\"%s\",\"%s\",%t,%t,%s,\"%s\"\n", r.ForwardingProfile, r.Uuid, r.GroupTag,
				sliceToString(r.SourceHips), sliceToString(r.DestinationHips), r.LogSuccessfulTlsHandshakes, r.LogFailedTlsHandshakes, r.LogSetting, sliceToString(r.SslCertificates)))

			time.Sleep(100 * time.Millisecond)
		}
	}

	cfh.End()
}
