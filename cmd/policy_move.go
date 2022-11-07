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
	"time"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/util"
	"github.com/Songmu/prompter"
	easycsv "github.com/scottdware/go-easycsv"
	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
)

// moveCmd represents the move command
var policyMoveCmd = &cobra.Command{
	Use:   "move",
	Short: "Move multiple rules within a security, NAT or Policy-Based Forwarding policy",
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
			moveOptions := map[string]int{
				"after":  util.MoveAfter,
				"before": util.MoveBefore,
				"bottom": util.MoveBottom,
				"top":    util.MoveTop,
			}

			rules, err := easycsv.Open(f)
			if err != nil {
				log.Printf("CSV file error - %s", err)
				os.Exit(1)
			}

			numrules := len(rules)
			log.Printf("Moving %d rules", numrules)

			for _, rule := range rules {
				ruletype := rule[0]
				rulename := rule[2]
				ruledest := rule[3]
				targetrule := rule[4]
				v = rule[5]

				switch ruletype {
				case "security":
					r, err := c.Policies.Security.Get(v, rulename)
					if err != nil {
						log.Printf("Failed to retrieve Security rule: %s", err)
					}

					err = c.Policies.Security.MoveGroup(v, moveOptions[ruledest], targetrule, r)
					if err != nil {
						fmt.Println(err)
					}
				case "nat":
					r, err := c.Policies.Nat.Get(v, rulename)
					if err != nil {
						log.Printf("Failed to retrieve NAT rule: %s", err)
					}

					err = c.Policies.Nat.MoveGroup(v, moveOptions[ruledest], targetrule, r)
					if err != nil {
						fmt.Println(err)
					}
				case "pbf":
					r, err := c.Policies.PolicyBasedForwarding.Get(v, rulename)
					if err != nil {
						log.Printf("Failed to retrieve Policy-Based Forwarding rule: %s", err)
					}

					err = c.Policies.PolicyBasedForwarding.MoveGroup(v, moveOptions[ruledest], targetrule, r)
					if err != nil {
						fmt.Println(err)
					}
				}

				time.Sleep(100 * time.Millisecond)
			}
		case *pango.Panorama:
			moveOptions := map[string]int{
				"after":  util.MoveAfter,
				"before": util.MoveBefore,
				"bottom": util.MoveBottom,
				"top":    util.MoveTop,
			}

			rules, err := easycsv.Open(f)
			if err != nil {
				log.Printf("CSV file error - %s", err)
				os.Exit(1)
			}

			numrules := len(rules)
			log.Printf("Moving %d rules", numrules)

			for _, rule := range rules {
				ruletype := rule[0]
				l = rule[1]
				rulename := rule[2]
				ruledest := rule[3]
				targetrule := rule[4]
				dg := rule[5]

				switch l {
				case "pre":
					l = util.PreRulebase
				case "post":
					l = util.PostRulebase
				default:
					l = util.PostRulebase
				}

				switch ruletype {
				case "security":
					r, err := c.Policies.Security.Get(dg, l, rulename)
					if err != nil {
						log.Printf("Failed to retrieve Security rule: %s", err)
					}

					err = c.Policies.Security.MoveGroup(dg, l, moveOptions[ruledest], targetrule, r)
					if err != nil {
						fmt.Println(err)
					}
				case "nat":
					r, err := c.Policies.Nat.Get(dg, l, rulename)
					if err != nil {
						log.Printf("Failed to retrieve NAT rule: %s", err)
					}

					err = c.Policies.Nat.MoveGroup(dg, l, moveOptions[ruledest], targetrule, r)
					if err != nil {
						fmt.Println(err)
					}
				case "pbf":
					r, err := c.Policies.PolicyBasedForwarding.Get(dg, l, rulename)
					if err != nil {
						log.Printf("Failed to retrieve Policy-Based Forwarding rule: %s", err)
					}

					err = c.Policies.PolicyBasedForwarding.MoveGroup(dg, l, moveOptions[ruledest], targetrule, r)
					if err != nil {
						fmt.Println(err)
					}
				}

				time.Sleep(100 * time.Millisecond)
			}
		}
	},
}

func init() {
	policyCmd.AddCommand(policyMoveCmd)

	policyMoveCmd.Flags().StringVarP(&f, "file", "f", "", "Name of the CSV file")
	policyMoveCmd.Flags().StringVarP(&user, "user", "u", "", "User to connect to the device as")
	// policyMoveCmd.Flags().StringVarP(&pass, "pass", "p", "", "Password for the user account specified")
	policyMoveCmd.Flags().StringVarP(&device, "device", "d", "", "Firewall or Panorama device to connect to")
	policyMoveCmd.MarkFlagRequired("user")
	// policyMoveCmd.MarkFlagRequired("pass")
	policyMoveCmd.MarkFlagRequired("device")
	policyMoveCmd.MarkFlagRequired("file")
}
