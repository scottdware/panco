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
	"regexp"
	"time"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/util"
	"github.com/Songmu/prompter"
	easycsv "github.com/scottdware/go-easycsv"
	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
)

// groupCmd represents the group command
var policyGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "Group Security, NAT or Policy-Based Forwarding rules by tags",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
		keyrexp := regexp.MustCompile(`key=([0-9A-Za-z\=]+).*`)
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
			if t == "security" {
				rules, err := easycsv.Open(f)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				rc := len(rules)
				log.Printf("Grouping %d Security rules by tags", rc)

				for i, rule := range rules {
					xpath := fmt.Sprintf("/config/devices/entry[@name='localhost.localdomain']/vsys/entry[@name='%s']/rulebase/security/rules/entry[@name='%s']", v, rule[0])
					ele := fmt.Sprintf("<group-tag>%s</group-tag>", rule[1])

					_, err := resty.R().Post(fmt.Sprintf("https://%s/api/?type=config&action=set&xpath=%s&key=%s&element=%s", device, xpath, c.ApiKey, ele))
					if err != nil {
						formatkey := keyrexp.ReplaceAllString(err.Error(), "key=********")
						log.Printf("Line %d - failed to group rule by tag %s: %s", i+1, rule[0], formatkey)
					}

					time.Sleep(100 * time.Millisecond)
				}
			}

			if t == "nat" {
				rules, err := easycsv.Open(f)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				rc := len(rules)
				log.Printf("Grouping %d NAT rules by tags", rc)

				for i, rule := range rules {
					xpath := fmt.Sprintf("/config/devices/entry[@name='localhost.localdomain']/vsys/entry[@name='%s']/rulebase/nat/rules/entry[@name='%s']", v, rule[0])
					ele := fmt.Sprintf("<group-tag>%s</group-tag>", rule[1])

					_, err := resty.R().Post(fmt.Sprintf("https://%s/api/?type=config&action=set&xpath=%s&key=%s&element=%s", device, xpath, c.ApiKey, ele))
					if err != nil {
						formatkey := keyrexp.ReplaceAllString(err.Error(), "key=********")
						log.Printf("Line %d - failed to group rule by tag %s: %s", i+1, rule[0], formatkey)
					}

					time.Sleep(100 * time.Millisecond)
				}
			}

			if t == "pbf" {
				rules, err := easycsv.Open(f)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				rc := len(rules)
				log.Printf("Grouping %d Policy-Based Forwarding rules by tags", rc)

				for i, rule := range rules {
					xpath := fmt.Sprintf("/config/devices/entry[@name='localhost.localdomain']/vsys/entry[@name='%s']/rulebase/pbf/rules/entry[@name='%s']", v, rule[0])
					ele := fmt.Sprintf("<group-tag>%s</group-tag>", rule[1])

					_, err := resty.R().Post(fmt.Sprintf("https://%s/api/?type=config&action=set&xpath=%s&key=%s&element=%s", device, xpath, c.ApiKey, ele))
					if err != nil {
						formatkey := keyrexp.ReplaceAllString(err.Error(), "key=********")
						log.Printf("Line %d - failed to group rule by tag %s: %s", i+1, rule[0], formatkey)
					}

					time.Sleep(100 * time.Millisecond)
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

			if t == "security" {
				// if dg == "" {
				// 	log.Printf("You must specify the device-group (use the -g or --devicegroup flag)")
				// 	os.Exit(1)
				// }

				rules, err := easycsv.Open(f)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				rc := len(rules)
				log.Printf("Grouping %d Security rules by tags", rc)

				for i, rule := range rules {
					var xpath string

					if dg == "shared" {
						xpath = fmt.Sprintf("/config/shared/%s/security/rules/entry[@name='%s']", l, rule[0])
					}

					if dg != "shared" {
						xpath = fmt.Sprintf("/config/devices/entry[@name='localhost.localdomain']/device-group/entry[@name='%s']/%s/security/rules/entry[@name='%s']", dg, l, rule[0])
					}

					ele := fmt.Sprintf("<group-tag>%s</group-tag>", rule[1])

					_, err := resty.R().Post(fmt.Sprintf("https://%s/api/?type=config&action=set&xpath=%s&key=%s&element=%s", device, xpath, c.ApiKey, ele))
					if err != nil {
						formatkey := keyrexp.ReplaceAllString(err.Error(), "key=********")
						log.Printf("Line %d - failed to group rule by tag %s: %s", i+1, rule[0], formatkey)
					}

					time.Sleep(100 * time.Millisecond)
				}
			}

			if t == "nat" {
				// if dg == "" {
				// 	log.Printf("You must specify the device-group (use the -g or --devicegroup flag)")
				// 	os.Exit(1)
				// }

				rules, err := easycsv.Open(f)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				rc := len(rules)
				log.Printf("Grouping %d NAT rules by tags", rc)

				for i, rule := range rules {
					var xpath string

					if dg == "shared" {
						xpath = fmt.Sprintf("/config/shared/%s/nat/rules/entry[@name='%s']", l, rule[0])
					}

					if dg != "shared" {
						xpath = fmt.Sprintf("/config/devices/entry[@name='localhost.localdomain']/device-group/entry[@name='%s']/%s/nat/rules/entry[@name='%s']", dg, l, rule[0])
					}

					ele := fmt.Sprintf("<group-tag>%s</group-tag>", rule[1])

					_, err := resty.R().Post(fmt.Sprintf("https://%s/api/?type=config&action=set&xpath=%s&key=%s&element=%s", device, xpath, c.ApiKey, ele))
					if err != nil {
						formatkey := keyrexp.ReplaceAllString(err.Error(), "key=********")
						log.Printf("Line %d - failed to group rule by tag %s: %s", i+1, rule[0], formatkey)
					}

					time.Sleep(100 * time.Millisecond)
				}
			}

			if t == "pbf" {
				// if dg == "" {
				// 	log.Printf("You must specify the device-group (use the -g or --devicegroup flag)")
				// 	os.Exit(1)
				// }

				rules, err := easycsv.Open(f)
				if err != nil {
					log.Printf("CSV file error - %s", err)
					os.Exit(1)
				}

				rc := len(rules)
				log.Printf("Grouping %d Policy-Based Forwarding rules by tags", rc)

				for i, rule := range rules {
					var xpath string

					if dg == "shared" {
						xpath = fmt.Sprintf("/config/shared/%s/pbf/rules/entry[@name='%s']", l, rule[0])
					}

					if dg != "shared" {
						xpath = fmt.Sprintf("/config/devices/entry[@name='localhost.localdomain']/device-group/entry[@name='%s']/%s/pbf/rules/entry[@name='%s']", dg, l, rule[0])
					}

					ele := fmt.Sprintf("<group-tag>%s</group-tag>", rule[1])

					_, err := resty.R().Post(fmt.Sprintf("https://%s/api/?type=config&action=set&xpath=%s&key=%s&element=%s", device, xpath, c.ApiKey, ele))
					if err != nil {
						formatkey := keyrexp.ReplaceAllString(err.Error(), "key=********")
						log.Printf("Line %d - failed to group rule by tag %s: %s", i+1, rule[0], formatkey)
					}

					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	},
}

func init() {
	policyCmd.AddCommand(policyGroupCmd)

	policyGroupCmd.Flags().StringVarP(&user, "user", "u", "", "User to connect to the device as")
	// policyGroupCmd.Flags().StringVarP(&pass, "pass", "p", "", "Password for the user account specified")
	policyGroupCmd.Flags().StringVarP(&device, "device", "d", "", "Device to connect to")
	policyGroupCmd.Flags().StringVarP(&f, "file", "f", "", "Name of the CSV file")
	policyGroupCmd.Flags().StringVarP(&dg, "devicegroup", "g", "shared", "Device Group name")
	policyGroupCmd.Flags().StringVarP(&v, "vsys", "v", "vsys1", "Vsys name")
	policyGroupCmd.Flags().StringVarP(&t, "type", "t", "", "<security|nat|pbf>")
	policyGroupCmd.MarkFlagRequired("user")
	// policyGroupCmd.MarkFlagRequired("pass")
	policyGroupCmd.MarkFlagRequired("device")
	policyGroupCmd.MarkFlagRequired("file")
	policyGroupCmd.MarkFlagRequired("type")
}
