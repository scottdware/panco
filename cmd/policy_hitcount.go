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
	"regexp"
	"strings"
	"time"

	"github.com/PaloAltoNetworks/pango"
	"github.com/Songmu/prompter"
	easycsv "github.com/scottdware/go-easycsv"
	"github.com/spf13/cobra"
)

type rule struct {
	Name           string `xml:"name,attr"`
	HitCount       int64  `xml:"hit-count"`
	LastHitTime    int64  `xml:"last-hit-timestamp"`
	LastResetTime  int64  `xml:"last-reset-timestamp"`
	FirstHitTime   int64  `xml:"first-hit-timestamp"`
	RuleCreateTime int64  `xml:"rule-creation-timestamp"`
	RuleModifyTime int64  `xml:"rule-modification-timestamp"`
}

type response struct {
	Rules []rule `xml:"result>rule-hit-count>vsys>entry>rule-base>entry>rules>entry"`
}

var (
	keyrexp = regexp.MustCompile(`key=([0-9A-Za-z\=]+).*`)
)

// hitcountCmd represents the hitcount command
var policyHitCountCmd = &cobra.Command{
	Use:   "hitcount",
	Short: "Get the hit count data on a Security, NAT or Policy-Based Forwarding policy - FIREWALL ONLY",
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
			if t == "security" || t == "all" {
				getFwSecHits(c, f)
			}

			// NAT policy
			if t == "nat" || t == "all" {
				getFwNatHits(c, f)
			}

			// Policy-Based Forwarding policy
			if t == "pbf" || t == "all" {
				getFwPbfHits(c, f)
			}
		case *pango.Panorama:
			log.Printf("You can only retreive hit count data from a firewall at this time")
			os.Exit(1)
			// switch l {
			// case "pre":
			// 	l = util.PreRulebase
			// case "post":
			// 	l = util.PostRulebase
			// default:
			// 	l = util.PostRulebase
			// }

			// Security policy
			// if t == "security" || t == "all" {
			// 	getPanoSecHits(c, f)
			// }

			// NAT policy
			// if t == "nat" || t == "all" {
			// 	getPanoNatHits(c, f)
			// }

			// PBF policy
			// if t == "pbf" || t == "all" {
			// 	getPanoPbfHits(c, f)
			// }
		}
	},
}

func init() {
	policyCmd.AddCommand(policyHitCountCmd)

	policyHitCountCmd.Flags().StringVarP(&user, "user", "u", "", "User to connect to the device as")
	// policyHitCountCmd.Flags().StringVarP(&pass, "pass", "p", "", "Password for the user account specified")
	policyHitCountCmd.Flags().StringVarP(&device, "device", "d", "", "Device to connect to")
	policyHitCountCmd.Flags().StringVarP(&f, "file", "f", "PaloAltoPolicy", "Name of the CSV file you'd like to export to")
	// policyHitCountCmd.Flags().StringVarP(&dg, "devicegroup", "g", "shared", "Device Group name when exporting from Panorama")
	policyHitCountCmd.Flags().StringVarP(&v, "vsys", "v", "vsys1", "Vsys name when exporting from a firewall")
	policyHitCountCmd.Flags().StringVarP(&t, "type", "t", "", "Type of policy to gather hit count on - <security|nat|pbf|all>")
	policyHitCountCmd.MarkFlagRequired("user")
	// policyHitCountCmd.MarkFlagRequired("pass")
	policyHitCountCmd.MarkFlagRequired("device")
	policyHitCountCmd.MarkFlagRequired("file")
	policyHitCountCmd.MarkFlagRequired("type")
}

// getFwSecHits is used to export the Security policy hit count data on a firewall
func getFwSecHits(c *pango.Firewall, file string) {
	var resp response
	command := fmt.Sprintf("<show><rule-hit-count><vsys><vsys-name><entry name='%s'><rule-base><entry name='security'><rules><all></all></rules></entry></rule-base></entry></vsys-name></vsys></rule-hit-count></show>", v)
	secHitFile := fmt.Sprintf("%s-Security_HitCount.csv", file)

	_, err := c.Op(command, "", nil, &resp)
	if err != nil {
		formatkey := keyrexp.ReplaceAllString(err.Error(), "key=********")
		log.Printf("Failed to get hit count on Security rules: %s", formatkey)
		return
	}

	if len(resp.Rules) <= 0 {
		log.Printf("There are no Security rules to export hit count data on")
		return
	}

	cfh, err := easycsv.NewCSV(secHitFile)
	if err != nil {
		log.Printf("CSV file error - %s", err)
		return
	}

	log.Printf("Exporting hit count data on %d Security rules", len(resp.Rules))

	cfh.Write("#Name,Hit Count,First Hit,Last Hit,Last Reset,Rule Created,Rule Modified\n")
	for _, r := range resp.Rules {
		cfh.Write(fmt.Sprintf("%s,%d,%v,%v,%v,%v,%v\n", r.Name, r.HitCount, time.Unix(r.FirstHitTime, 0).Format(time.RFC3339),
			time.Unix(r.LastHitTime, 0).Format(time.RFC3339), time.Unix(r.LastResetTime, 0).Format(time.RFC3339), time.Unix(r.RuleCreateTime, 0).Format(time.RFC3339),
			time.Unix(r.RuleModifyTime, 0).Format(time.RFC3339)))
	}

	cfh.End()
}

// getFwNatHits is used to export the NAT policy hit count data on a firewall
func getFwNatHits(c *pango.Firewall, file string) {
	var resp response
	command := fmt.Sprintf("<show><rule-hit-count><vsys><vsys-name><entry name='%s'><rule-base><entry name='nat'><rules><all></all></rules></entry></rule-base></entry></vsys-name></vsys></rule-hit-count></show>", v)
	secHitFile := fmt.Sprintf("%s-NAT_HitCount.csv", file)

	_, err := c.Op(command, "", nil, &resp)
	if err != nil {
		formatkey := keyrexp.ReplaceAllString(err.Error(), "key=********")
		log.Printf("Failed to get hit count on NAT rules: %s", formatkey)
		return
	}

	if len(resp.Rules) <= 0 {
		log.Printf("There are no NAT rules to export hit count data on")
		return
	}

	cfh, err := easycsv.NewCSV(secHitFile)
	if err != nil {
		log.Printf("CSV file error - %s", err)
		return
	}

	log.Printf("Exporting hit count data on %d NAT rules", len(resp.Rules))

	cfh.Write("#Name,Hit Count,First Hit,Last Hit,Last Reset,Rule Created,Rule Modified\n")
	for _, r := range resp.Rules {
		cfh.Write(fmt.Sprintf("%s,%d,%v,%v,%v,%v,%v\n", r.Name, r.HitCount, time.Unix(r.FirstHitTime, 0).Format(time.RFC3339),
			time.Unix(r.LastHitTime, 0).Format(time.RFC3339), time.Unix(r.LastResetTime, 0).Format(time.RFC3339), time.Unix(r.RuleCreateTime, 0).Format(time.RFC3339),
			time.Unix(r.RuleModifyTime, 0).Format(time.RFC3339)))
	}

	cfh.End()
}

// getFwPbfHits is used to export the Policy-Based Forwarding policy hit count data on a firewall
func getFwPbfHits(c *pango.Firewall, file string) {
	var resp response
	command := fmt.Sprintf("<show><rule-hit-count><vsys><vsys-name><entry name='%s'><rule-base><entry name='pbf'><rules><all></all></rules></entry></rule-base></entry></vsys-name></vsys></rule-hit-count></show>", v)
	secHitFile := fmt.Sprintf("%s-PBF_HitCount.csv", file)

	_, err := c.Op(command, "", nil, &resp)
	if err != nil {
		formatkey := keyrexp.ReplaceAllString(err.Error(), "key=********")
		log.Printf("Failed to get hit count on Policy-Based Forwarding rules: %s", formatkey)
		return
	}

	if len(resp.Rules) <= 0 {
		log.Printf("There are no Policy-Based Forwarding rules to export hit count data on")
		return
	}

	cfh, err := easycsv.NewCSV(secHitFile)
	if err != nil {
		log.Printf("CSV file error - %s", err)
		return
	}

	log.Printf("Exporting hit count data on %d Policy-Based Forwarding rules", len(resp.Rules))

	cfh.Write("#Name,Hit Count,First Hit,Last Hit,Last Reset,Rule Created,Rule Modified\n")
	for _, r := range resp.Rules {
		cfh.Write(fmt.Sprintf("%s,%d,%v,%v,%v,%v,%v\n", r.Name, r.HitCount, time.Unix(r.FirstHitTime, 0).Format(time.RFC3339),
			time.Unix(r.LastHitTime, 0).Format(time.RFC3339), time.Unix(r.LastResetTime, 0).Format(time.RFC3339), time.Unix(r.RuleCreateTime, 0).Format(time.RFC3339),
			time.Unix(r.RuleModifyTime, 0).Format(time.RFC3339)))
	}

	cfh.End()
}

// getPanoSecHits is used to export the security policy hit count data from Panorama
// func getPanoSecHits(c *pango.Panorama, file string) {
// 	var resp response
// 	command := fmt.Sprintf("<show><rule-hit-count><vsys><vsys-name><entry name='%s'><rule-base><entry name='security'><rules><all></all></rules></entry></rule-base></entry></vsys-name></vsys></rule-hit-count></show>", v)
// 	secHitFile := fmt.Sprintf("%s-Security_HitCount.csv", file)

// 	_, err := c.Op(command, "", nil, &resp)
// 	if err != nil {
// 		formatkey := keyrexp.ReplaceAllString(err.Error(), "key=********")
// 		log.Printf("Failed to get hit count on security rules: %s", formatkey)
// 		return
// 	}

// 	cfh, err := easycsv.NewCSV(secHitFile)
// 	if err != nil {
// 		log.Printf("CSV file error - %s", err)
// 		return
// 	}

// 	cfh.Write("#Name,Hit Count,First Hit,Last Hit,Last Reset,Rule Created,Rule Modified\n")
// 	for _, r := range resp.Rules {
// 		cfh.Write(fmt.Sprintf("%s,%d,%v,%v,%v,%v,%v\n", r.Name, r.HitCount, time.Unix(r.FirstHitTime, 0).Format(time.RFC3339),
// 			time.Unix(r.LastHitTime, 0).Format(time.RFC3339), time.Unix(r.LastResetTime, 0).Format(time.RFC3339), time.Unix(r.RuleCreateTime, 0).Format(time.RFC3339),
// 			time.Unix(r.RuleModifyTime, 0).Format(time.RFC3339)))
// 	}

// 	cfh.End()
// }

// getPanoNatHits is used to export the NAT policy hit count data from Panorama
// func getPanoNatHits(c *pango.Panorama, file string) {
// 	var resp response
// 	command := fmt.Sprintf("<show><rule-hit-count><vsys><vsys-name><entry name='%s'><rule-base><entry name='nat'><rules><all></all></rules></entry></rule-base></entry></vsys-name></vsys></rule-hit-count></show>", v)
// 	secHitFile := fmt.Sprintf("%s-NAT_HitCount.csv", file)

// 	_, err := c.Op(command, "", nil, &resp)
// 	if err != nil {
// 		formatkey := keyrexp.ReplaceAllString(err.Error(), "key=********")
// 		log.Printf("Failed to get hit count on NAT rules: %s", formatkey)
// 		return
// 	}

// 	cfh, err := easycsv.NewCSV(secHitFile)
// 	if err != nil {
// 		log.Printf("CSV file error - %s", err)
// 		return
// 	}

// 	cfh.Write("#Name,Hit Count,First Hit,Last Hit,Last Reset,Rule Created,Rule Modified\n")
// 	for _, r := range resp.Rules {
// 		cfh.Write(fmt.Sprintf("%s,%d,%v,%v,%v,%v,%v\n", r.Name, r.HitCount, time.Unix(r.FirstHitTime, 0).Format(time.RFC3339),
// 			time.Unix(r.LastHitTime, 0).Format(time.RFC3339), time.Unix(r.LastResetTime, 0).Format(time.RFC3339), time.Unix(r.RuleCreateTime, 0).Format(time.RFC3339),
// 			time.Unix(r.RuleModifyTime, 0).Format(time.RFC3339)))
// 	}

// 	cfh.End()
// }

// getPanoPbfHits is used to export the Policy-Based Forwarding policy hit count from Panorama
// func getPanoPbfHits(c *pango.Panorama, file string) {
// 	var resp response
// 	command := fmt.Sprintf("<show><rule-hit-count><vsys><vsys-name><entry name='%s'><rule-base><entry name='pbf'><rules><all></all></rules></entry></rule-base></entry></vsys-name></vsys></rule-hit-count></show>", v)
// 	secHitFile := fmt.Sprintf("%s-PBF_HitCount.csv", file)

// 	_, err := c.Op(command, "", nil, &resp)
// 	if err != nil {
// 		formatkey := keyrexp.ReplaceAllString(err.Error(), "key=********")
// 		log.Printf("Failed to get hit count on Policy-Based Forwarding rules: %s", formatkey)
// 		return
// 	}

// 	cfh, err := easycsv.NewCSV(secHitFile)
// 	if err != nil {
// 		log.Printf("CSV file error - %s", err)
// 		return
// 	}

// 	cfh.Write("#Name,Hit Count,First Hit,Last Hit,Last Reset,Rule Created,Rule Modified\n")
// 	for _, r := range resp.Rules {
// 		cfh.Write(fmt.Sprintf("%s,%d,%v,%v,%v,%v,%v\n", r.Name, r.HitCount, time.Unix(r.FirstHitTime, 0).Format(time.RFC3339),
// 			time.Unix(r.LastHitTime, 0).Format(time.RFC3339), time.Unix(r.LastResetTime, 0).Format(time.RFC3339), time.Unix(r.RuleCreateTime, 0).Format(time.RFC3339),
// 			time.Unix(r.RuleModifyTime, 0).Format(time.RFC3339)))
// 	}

// 	cfh.End()
// }
