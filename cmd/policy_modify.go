/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/PaloAltoNetworks/pango"
	"github.com/Songmu/prompter"
	easycsv "github.com/scottdware/go-easycsv"
	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
)

// modifyCmd represents the modify command
var policyModifyCmd = &cobra.Command{
	Use:   "modify",
	Short: "Modify existing rules - add source, destination objects to them",
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
			rules, err := easycsv.Open(f)
			if err != nil {
				log.Printf("CSV file error - %s", err)
				os.Exit(1)
			}

			for _, rule := range rules {
				name := rule[0]
				ruleloc := rule[1]
				action := rule[2]
				value := stringToSlice(rule[3])
				dgroup := rule[4]

				switch action {
				case "addsource":
					var xpath, xmlBody string

					if dgroup == "shared" {
						for _, src := range value {
							xmlBody += fmt.Sprintf("<member>%s</member>", src)
						}
						xpath = fmt.Sprintf("/config/shared/%s-rulebase/security/rules/entry[@name='%s']/source", ruleloc, name)
					}

					if dgroup != "shared" {
						for _, src := range value {
							xmlBody += fmt.Sprintf("<member>%s</member>", src)
						}
						xpath = fmt.Sprintf("/config/devices/entry[@name='localhost.localdomain']/device-group/entry[@name='%s']/%s-rulebase/security/rules/entry[@name='%s']/source", dgroup, ruleloc, name)
					}

					_, err := resty.R().Get(fmt.Sprintf("https://%s/api/?type=config&action=set&xpath=%s&element=%s&key=%s", device, xpath, xmlBody, c.ApiKey))
					if err != nil {
						formatkey := keyrexp.ReplaceAllString(err.Error(), "key=********")
						log.Printf("Line %d - failed to add sources to rule %s: %s", i+1, name, formatkey)
					}
				case "adddest":
					var xpath, xmlBody string

					if dgroup == "shared" {
						for _, src := range value {
							xmlBody += fmt.Sprintf("<member>%s</member>", src)
						}
						xpath = fmt.Sprintf("/config/shared/%s-rulebase/security/rules/entry[@name='%s']/destination", ruleloc, name)
					}

					if dgroup != "shared" {
						for _, src := range value {
							xmlBody += fmt.Sprintf("<member>%s</member>", src)
						}
						xpath = fmt.Sprintf("/config/devices/entry[@name='localhost.localdomain']/device-group/entry[@name='%s']/%s-rulebase/security/rules/entry[@name='%s']/destination", dgroup, ruleloc, name)
					}

					_, err := resty.R().Get(fmt.Sprintf("https://%s/api/?type=config&action=set&xpath=%s&element=%s&key=%s", device, xpath, xmlBody, c.ApiKey))
					if err != nil {
						formatkey := keyrexp.ReplaceAllString(err.Error(), "key=********")
						log.Printf("Line %d - failed to add sources to rule %s: %s", i+1, name, formatkey)
					}
				}

				time.Sleep(100 * time.Millisecond)
			}
		case *pango.Panorama:
			rules, err := easycsv.Open(f)
			if err != nil {
				log.Printf("CSV file error - %s", err)
				os.Exit(1)
			}

			for _, rule := range rules {
				name := rule[0]
				ruleloc := rule[1]
				action := rule[2]
				value := stringToSlice(rule[3])
				dgroup := rule[4]

				switch action {
				case "addsource":
					var xpath, xmlBody string

					if dgroup == "shared" {
						for _, src := range value {
							xmlBody += fmt.Sprintf("<member>%s</member>", src)
						}
						xpath = fmt.Sprintf("/config/shared/%s-rulebase/security/rules/entry[@name='%s']/source", ruleloc, name)
					}

					if dgroup != "shared" {
						for _, src := range value {
							xmlBody += fmt.Sprintf("<member>%s</member>", src)
						}
						xpath = fmt.Sprintf("/config/devices/entry[@name='localhost.localdomain']/device-group/entry[@name='%s']/%s-rulebase/security/rules/entry[@name='%s']/source", dgroup, ruleloc, name)
					}

					_, err := resty.R().Get(fmt.Sprintf("https://%s/api/?type=config&action=set&xpath=%s&element=%s&key=%s", device, xpath, xmlBody, c.ApiKey))
					if err != nil {
						formatkey := keyrexp.ReplaceAllString(err.Error(), "key=********")
						log.Printf("Line %d - failed to add sources to rule %s: %s", i+1, name, formatkey)
					}
				case "adddest":
					var xpath, xmlBody string

					if dgroup == "shared" {
						for _, src := range value {
							xmlBody += fmt.Sprintf("<member>%s</member>", src)
						}
						xpath = fmt.Sprintf("/config/shared/%s-rulebase/security/rules/entry[@name='%s']/destination", ruleloc, name)
					}

					if dgroup != "shared" {
						for _, src := range value {
							xmlBody += fmt.Sprintf("<member>%s</member>", src)
						}
						xpath = fmt.Sprintf("/config/devices/entry[@name='localhost.localdomain']/device-group/entry[@name='%s']/%s-rulebase/security/rules/entry[@name='%s']/destination", dgroup, ruleloc, name)
					}

					_, err := resty.R().Get(fmt.Sprintf("https://%s/api/?type=config&action=set&xpath=%s&element=%s&key=%s", device, xpath, xmlBody, c.ApiKey))
					if err != nil {
						formatkey := keyrexp.ReplaceAllString(err.Error(), "key=********")
						log.Printf("Line %d - failed to add sources to rule %s: %s", i+1, name, formatkey)
					}
				}

				time.Sleep(100 * time.Millisecond)
			}
		}
	},
}

func init() {
	policyCmd.AddCommand(policyModifyCmd)

	policyModifyCmd.Flags().StringVarP(&f, "file", "f", "", "Name of the CSV file")
	policyModifyCmd.Flags().StringVarP(&user, "user", "u", "", "User to connect to the device as")
	// policyMoveCmd.Flags().StringVarP(&pass, "pass", "p", "", "Password for the user account specified")
	policyModifyCmd.Flags().StringVarP(&device, "device", "d", "", "Firewall or Panorama device to connect to")
	policyModifyCmd.MarkFlagRequired("user")
	// policyMoveCmd.MarkFlagRequired("pass")
	policyModifyCmd.MarkFlagRequired("device")
	policyModifyCmd.MarkFlagRequired("file")
}
