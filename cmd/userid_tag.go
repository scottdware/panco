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

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/userid"
	"github.com/Songmu/prompter"
	easycsv "github.com/scottdware/go-easycsv"
	"github.com/spf13/cobra"
)

// useridTagCmd represents the tag command
var useridTagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Tag IP addresses for use in dynamic address groups",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		ipmatch := regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)
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
			lines, err := easycsv.Open(f)
			if err != nil {
				log.Printf("CSV file error - %s", err)
				os.Exit(1)
			}

			lc := len(lines)
			log.Printf("Running tag options on %d IP addresses", lc)

			for _, line := range lines {
				var typeis string
				name := line[0]
				action := line[1]
				value := line[2]

				if ipmatch.Match([]byte(name)) {
					typeis = "ip"
				} else {
					typeis = "user"
				}

				switch typeis {
				case "ip":
					if action == "tag" {
						message := userid.Message{}

						entry := userid.TagIp{
							Ip:   name,
							Tags: stringToSlice(value),
						}

						message.TagIps = append(message.TagIps, entry)

						err := c.UserId.Run(&message, v)
						if err != nil {
							log.Printf("Error tagging IP address: %s", err)
							os.Exit(1)
						}
					}

					if action == "untag" {
						message := userid.Message{}

						entry := userid.UntagIp{
							Ip:   name,
							Tags: stringToSlice(value),
						}

						message.UntagIps = append(message.UntagIps, entry)

						err := c.UserId.Run(&message, v)
						if err != nil {
							log.Printf("Error untagging IP address: %s", err)
							os.Exit(1)
						}
					}
				}
			}
		case *pango.Panorama:
			log.Printf("You can only tag IP addresses on a firewall at this time")
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
	useridCmd.AddCommand(useridTagCmd)

	useridTagCmd.Flags().StringVarP(&user, "user", "u", "", "User to connect to the device as")
	// useridTagCmd.Flags().StringVarP(&pass, "pass", "p", "", "Password for the user account specified")
	useridTagCmd.Flags().StringVarP(&device, "device", "d", "", "Device to connect to")
	useridTagCmd.Flags().StringVarP(&f, "file", "f", "", "Name of the CSV file to import")
	useridTagCmd.Flags().StringVarP(&v, "vsys", "v", "vsys1", "Vsys name")
	useridTagCmd.MarkFlagRequired("user")
	// useridTagCmd.MarkFlagRequired("pass")
	useridTagCmd.MarkFlagRequired("device")
	useridTagCmd.MarkFlagRequired("file")
	useridTagCmd.MarkFlagRequired("vsys")
}
