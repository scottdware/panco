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
	"os"
	"strings"
	"time"

	easycsv "github.com/scottdware/go-easycsv"
	panos "github.com/scottdware/go-panos"
	"github.com/spf13/cobra"
)

// policyCmd represents the policy command
var policyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Export/import a security policy",
	Long: `This command will allow you to export and import an entire security policy. If
you are running this against a Panorama device, it can be really helpful if you want to clone
an entire policy, as you can export it from one device-group, modify it if needed, then import
the poilcy into a different device-group.

For an example CSV format of how a policy import should look, use the --action export flag to
export a policy. The following columns in the CSV file must not be blank, and at the very minimum
have the value of "any" if you wish to allow that:

From, To, Source, Destination, SourceUser, Application, Service, HIPProfiles, Category

You must always specify the action you want to take via the --action flag. Actions are either
export or import.`,
	Run: func(cmd *cobra.Command, args []string) {
		pass := passwd()
		creds := &panos.AuthMethod{
			Credentials: []string{user, pass},
		}

		pan, err := panos.NewSession(device, creds)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if action == "export" {
			policies, err := pan.Policy(dg)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if !strings.Contains(fh, ".csv") {
				fh = fmt.Sprintf("%s.csv", fh)
			}

			csv, err := easycsv.NewCSV(fh)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			csv.Write("#DeviceGroup,Type,Name,Tag,From,To,Source,Destination,SourceUser,Application,Service,HIPProfiles,Category")
			csv.Write(",Action,LogStart,LogEnd,LogSetting,Disabled,URLFilteringProfile,FileBlockingProfile")
			csv.Write(",AntiVirusProfile,AntiSpywareProfile,VulnerabilityProfile,WildfireProfile,SecurityProfileGroup,Description\n")

			time.Sleep(50 * time.Millisecond)

			if len(policies.Local) > 0 {
				for _, p := range policies.Local {
					var tag string

					if len(p.Tag) > 0 {
						tag = sliceToString(p.Tag)
					}

					from := sliceToString(p.From)
					to := sliceToString(p.To)
					source := sliceToString(p.Source)
					dest := sliceToString(p.Destination)
					srcuser := sliceToString(p.SourceUser)
					app := sliceToString(p.Application)
					service := sliceToString(p.Service)
					hip := sliceToString(p.HIPProfiles)
					category := sliceToString(p.Category)

					csv.Write(fmt.Sprintf("\"%s\",%s,\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"",
						"", "local", p.Name, tag, from, to, source, dest, srcuser, app, service, hip, category))
					csv.Write(fmt.Sprintf(",%s,%s,%s,%s,%s,%s,%s",
						p.Action, p.LogStart, p.LogEnd, p.LogSetting, p.Disabled, p.URLFilteringProfile, p.FileBlockingProfile))
					csv.Write(fmt.Sprintf(",%s,%s,%s,%s,%s,\"%s\"\n",
						p.AntiVirusProfile, p.AntiSpywareProfile, p.VulnerabilityProfile, p.WildfireProfile, p.SecurityProfileGroup, p.Description))

					time.Sleep(10 * time.Millisecond)
				}
			}

			if len(policies.Pre) > 0 {
				for _, p := range policies.Pre {
					var tag string

					if len(p.Tag) > 0 {
						tag = sliceToString(p.Tag)
					}

					from := sliceToString(p.From)
					to := sliceToString(p.To)
					source := sliceToString(p.Source)
					dest := sliceToString(p.Destination)
					srcuser := sliceToString(p.SourceUser)
					app := sliceToString(p.Application)
					service := sliceToString(p.Service)
					hip := sliceToString(p.HIPProfiles)
					category := sliceToString(p.Category)

					csv.Write(fmt.Sprintf("\"%s\",%s,\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"",
						dg, "pre", p.Name, tag, from, to, source, dest, srcuser, app, service, hip, category))
					csv.Write(fmt.Sprintf(",%s,%s,%s,%s,%s,%s,%s",
						p.Action, p.LogStart, p.LogEnd, p.LogSetting, p.Disabled, p.URLFilteringProfile, p.FileBlockingProfile))
					csv.Write(fmt.Sprintf(",%s,%s,%s,%s,%s,\"%s\"\n",
						p.AntiVirusProfile, p.AntiSpywareProfile, p.VulnerabilityProfile, p.WildfireProfile, p.SecurityProfileGroup, p.Description))

					time.Sleep(10 * time.Millisecond)
				}
			}

			if len(policies.Post) > 0 {
				for _, p := range policies.Post {
					var tag string

					if len(p.Tag) > 0 {
						tag = sliceToString(p.Tag)
					}

					from := sliceToString(p.From)
					to := sliceToString(p.To)
					source := sliceToString(p.Source)
					dest := sliceToString(p.Destination)
					srcuser := sliceToString(p.SourceUser)
					app := sliceToString(p.Application)
					service := sliceToString(p.Service)
					hip := sliceToString(p.HIPProfiles)
					category := sliceToString(p.Category)

					csv.Write(fmt.Sprintf("\"%s\",%s,\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"",
						dg, "post", p.Name, tag, from, to, source, dest, srcuser, app, service, hip, category))
					csv.Write(fmt.Sprintf(",%s,%s,%s,%s,%s,%s,%s",
						p.Action, p.LogStart, p.LogEnd, p.LogSetting, p.Disabled, p.URLFilteringProfile, p.FileBlockingProfile))
					csv.Write(fmt.Sprintf(",%s,%s,%s,%s,%s,\"%s\"\n",
						p.AntiVirusProfile, p.AntiSpywareProfile, p.VulnerabilityProfile, p.WildfireProfile, p.SecurityProfileGroup, p.Description))

					time.Sleep(10 * time.Millisecond)
				}
			}

			csv.End()
		}

		if action == "import" {
			rules, err := easycsv.Open(fh)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			for _, rule := range rules {
				var tag []string

				if len(rule[3]) > 0 {
					tag = stringToSlice(rule[3])
				}

				from := stringToSlice(rule[4])
				to := stringToSlice(rule[5])
				source := stringToSlice(rule[6])
				dest := stringToSlice(rule[7])
				srcuser := stringToSlice(rule[8])
				app := stringToSlice(rule[9])
				service := stringToSlice(rule[10])
				hip := stringToSlice(rule[11])
				category := stringToSlice(rule[12])

				content := &panos.RuleContent{
					Name:                 rule[2],
					Tag:                  tag,
					From:                 from,
					To:                   to,
					Source:               source,
					Destination:          dest,
					SourceUser:           srcuser,
					Application:          app,
					Service:              service,
					HIPProfiles:          hip,
					Category:             category,
					Action:               rule[13],
					LogStart:             rule[14],
					LogEnd:               rule[15],
					LogSetting:           rule[16],
					Disabled:             rule[17],
					URLFilteringProfile:  rule[18],
					FileBlockingProfile:  rule[19],
					AntiVirusProfile:     rule[20],
					AntiSpywareProfile:   rule[21],
					VulnerabilityProfile: rule[22],
					WildfireProfile:      rule[23],
					SecurityProfileGroup: rule[24],
					Description:          rule[25],
				}

				err = pan.CreateRule(rule[2], rule[1], content)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				time.Sleep(100 * time.Millisecond)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(policyCmd)

	policyCmd.Flags().StringVarP(&action, "action", "a", "", "Action to perform - export or import")
	policyCmd.Flags().StringVarP(&fh, "file", "f", "", "Name of the CSV file to export/import")
	policyCmd.Flags().StringVarP(&dg, "devicegroup", "g", "", "Device group - only needed when ran against Panorama")
	policyCmd.Flags().StringVarP(&user, "user", "u", "", "User to connect to the device as")
	policyCmd.Flags().StringVarP(&device, "device", "d", "", "Firewall or Panorama device to connect to")
	policyCmd.MarkFlagRequired("user")
	policyCmd.MarkFlagRequired("device")
	policyCmd.MarkFlagRequired("action")
	policyCmd.MarkFlagRequired("file")
}
