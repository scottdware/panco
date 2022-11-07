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

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/userid"
	"github.com/Songmu/prompter"
	easycsv "github.com/scottdware/go-easycsv"
	"github.com/spf13/cobra"
)

// useridLoginCmd represents the login command
var useridLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Manually register/unregister (map) users to an IP address (login/logout)",
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

		switch c := con.(type) {
		case *pango.Firewall:

			lines, err := easycsv.Open(f)
			if err != nil {
				log.Printf("CSV file error - %s", err)
				os.Exit(1)
			}

			lc := len(lines)
			log.Printf("Running login options on %d users", lc)

			for _, line := range lines {
				name := line[0]
				action := line[1]
				value := line[2]

				if action == "login" {
					message := userid.Message{}

					entry := userid.Login{
						User: name,
						Ip:   value,
					}

					message.Logins = append(message.Logins, entry)

					err := c.UserId.Run(&message, v)
					if err != nil {
						log.Printf("Error registering/logging in user: %s", err)
						os.Exit(1)
					}
				}

				if action == "logout" {
					message := userid.Message{}

					entry := userid.Logout{
						User: name,
						Ip:   value,
					}

					message.Logouts = append(message.Logouts, entry)

					err := c.UserId.Run(&message, v)
					if err != nil {
						log.Printf("Error unregistering/logging out user: %s", err)
						os.Exit(1)
					}
				}
			}
		case *pango.Panorama:
			log.Printf("You can only register/unregister users on a firewall at this time")
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
	useridCmd.AddCommand(useridLoginCmd)

	useridLoginCmd.Flags().StringVarP(&user, "user", "u", "", "User to connect to the device as")
	// useridLoginCmd.Flags().StringVarP(&pass, "pass", "p", "", "Password for the user account specified")
	useridLoginCmd.Flags().StringVarP(&device, "device", "d", "", "Device to connect to")
	useridLoginCmd.Flags().StringVarP(&f, "file", "f", "", "Name of the CSV file to import")
	useridLoginCmd.Flags().StringVarP(&v, "vsys", "v", "vsys1", "Vsys name")
	useridLoginCmd.MarkFlagRequired("user")
	// useridLoginCmd.MarkFlagRequired("pass")
	useridLoginCmd.MarkFlagRequired("device")
	useridLoginCmd.MarkFlagRequired("file")
	useridLoginCmd.MarkFlagRequired("vsys")
}
