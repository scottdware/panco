// Copyright © 2019 Scott Ware <scottdware@gmail.com>
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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/PaloAltoNetworks/pango"
	"github.com/levigross/grequests"
	"github.com/spf13/cobra"
)

type contents struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Size        int    `json:"size"`
	DownloadURL string `json:"download_url"`
	Type        string `json:"type"`
	Content     string `json:"content"`
}

// provisionCmd represents the provision command
var provisionCmd = &cobra.Command{
	Use:   "provision",
	Short: "Provision a device using IronSkillet or a local or remote (HTTP) file",
	Long: `This command will allow you to configure your device with best practice templates that have
been setup by Palo Alto Networks, in their IronSkillet Github repository. You can also choose to load a 
config file locally, or from a remote source using HTTP. Only the pre-built "loadable_configs" are an 
option at this time.

When using IronSkillet as the source (--source ironskillet), you will need to specifcy a couple of options:

'--os' must be either "panos" or "panorama"
'--config' can be one of: "aws", "azure", "gcp", "dhcp" or "static"

If you do not specify the --load flag, the the configuration will only be transferred to the device, where you
will have to load it manually. If you specify the --load flag, then the configuration will be loaded automatically,
but will still need to be manually committed.

See https://github.com/scottdware/panco/Wiki for more information`,
	Run: func(cmd *cobra.Command, args []string) {
		var tmpl contents
		var err error
		pass := passwd()
		configs := map[string]string{
			"aws":    "https://api.github.com/repos/PaloAltoNetworks/iron-skillet/contents/loadable_configs/sample-cloud-AWS/%s/iron_skillet_%s_full.xml",
			"azure":  "https://api.github.com/repos/PaloAltoNetworks/iron-skillet/contents/loadable_configs/sample-cloud-Azure/%s/iron_skillet_%s_full.xml",
			"gcp":    "https://api.github.com/repos/PaloAltoNetworks/iron-skillet/contents/loadable_configs/sample-cloud-GCP/%s/iron_skillet_%s_full.xml",
			"dhcp":   "https://api.github.com/repos/PaloAltoNetworks/iron-skillet/contents/loadable_configs/sample-mgmt-dhcp/%s/iron_skillet_%s_full.xml",
			"static": "https://api.github.com/repos/PaloAltoNetworks/iron-skillet/contents/loadable_configs/sample-mgmt-static/%s/iron_skillet_%s_full.xml",
		}

		cl := pango.Client{
			Hostname: device,
			Username: user,
			Password: pass,
			Logging:  pango.LogQuiet,
		}

		con, err := pango.Connect(cl)
		if err != nil {
			log.Printf("Failed to connect: %s", err)
			os.Exit(1)
		}

		if source == "ironskillet" {
			if ostype == "" {
				log.Println("You must choose the device OS - panos or panorama")
				os.Exit(1)
			}

			if config == "" {
				log.Println("You didn't specify a config type...static will be used by default")
			}

			giturl := fmt.Sprintf(configs[config], ostype, ostype)
			resp, err := grequests.Get(giturl, nil)
			if err != nil {
				log.Printf("Error downloading IronSkillet file: %s", err)
				os.Exit(1)
			}

			if err = json.Unmarshal([]byte(resp.String()), &tmpl); err != nil {
				log.Printf("Error parsing JSON: %s", err)
				os.Exit(1)
			}

			fp := fmt.Sprintf("%s\\%s", os.TempDir(), tmpl.Name)
			dl, _ := grequests.Get(tmpl.DownloadURL, nil)
			if err := dl.DownloadToFile(fp); err != nil {
				log.Printf("Error saving file to temp dir: %s", err)
				os.Exit(1)
			}

			fdata, err := ioutil.ReadFile(fp)
			if err != nil {
				log.Printf("Error reading file: %s", err)
				os.Exit(1)
			}
			fcontent := string(fdata)

			switch c := con.(type) {
			case *pango.Firewall:
				vals := url.Values{}
				vals.Set("key", c.ApiKey)
				vals.Set("type", "import")
				vals.Set("category", "configuration")

				_, err = c.CommunicateFile(fcontent, tmpl.Name, fp, vals, nil)
				if err != nil {
					log.Printf("Failed to transfer file to device: %s", err)
				}

				if load == true {
					cmd := fmt.Sprintf("<load><config><from>%s</from></config></load>", tmpl.Name)
					vals := url.Values{}
					vals.Set("key", c.ApiKey)
					_, err = c.Op(cmd, "", vals, nil)
					if err != nil {
						log.Printf("Error loading file: %s", err)
					}
				}
			case *pango.Panorama:
				vals := url.Values{}
				vals.Set("key", c.ApiKey)
				vals.Set("type", "import")
				vals.Set("category", "configuration")

				_, err = c.CommunicateFile(fcontent, tmpl.Name, fp, vals, nil)
				if err != nil {
					log.Printf("Failed to transfer file to device: %s", err)
				}

				if load == true {
					cmd := fmt.Sprintf("<load><config><from>%s</from></config></load>", tmpl.Name)
					vals := url.Values{}
					vals.Set("key", c.ApiKey)
					_, err = c.Op(cmd, "", vals, nil)
					if err != nil {
						log.Printf("Error loading file: %s", err)
					}
				}
			}
		}

		if source == "local" {
			if fh == "" {
				log.Println("You must specify a config (XML) file to use")
				os.Exit(1)
			}

			fbase := filepath.Base(fh)
			fdata, err := ioutil.ReadFile(fh)
			if err != nil {
				log.Printf("Error reading file: %s", err)
				os.Exit(1)
			}
			fcontent := string(fdata)

			switch c := con.(type) {
			case *pango.Firewall:
				vals := url.Values{}
				vals.Set("key", c.ApiKey)
				vals.Set("type", "import")
				vals.Set("category", "configuration")

				_, err = c.CommunicateFile(fcontent, fbase, fh, vals, nil)
				if err != nil {
					log.Printf("Failed to transfer file to device: %s", err)
				}

				if load == true {
					cmd := fmt.Sprintf("<load><config><from>%s</from></config></load>", fbase)
					vals := url.Values{}
					vals.Set("key", c.ApiKey)
					_, err = c.Op(cmd, "", vals, nil)
					if err != nil {
						log.Printf("Error loading file: %s", err)
					}
				}
			case *pango.Panorama:
				vals := url.Values{}
				vals.Set("key", c.ApiKey)
				vals.Set("type", "import")
				vals.Set("category", "configuration")

				_, err = c.CommunicateFile(fcontent, fbase, fh, vals, nil)
				if err != nil {
					log.Printf("Failed to transfer file to device: %s", err)
				}

				if load == true {
					cmd := fmt.Sprintf("<load><config><from>%s</from></config></load>", fbase)
					vals := url.Values{}
					vals.Set("key", c.ApiKey)
					_, err = c.Op(cmd, "", vals, nil)
					if err != nil {
						log.Printf("Error loading file: %s", err)
					}
				}
			}
		}

		if source == "remote" {
			if fh == "" {
				log.Println("You must specify a config (XML) file to use")
				os.Exit(1)
			}

			fbase := filepath.Base(fh)
			fp := fmt.Sprintf("%s\\%s", os.TempDir(), fbase)

			dl, _ := grequests.Get(fh, nil)
			if err := dl.DownloadToFile(fp); err != nil {
				log.Printf("Error saving file to temp dir: %s", err)
				os.Exit(1)
			}

			fdata, err := ioutil.ReadFile(fp)
			if err != nil {
				log.Printf("Error reading file: %s", err)
				os.Exit(1)
			}
			fcontent := string(fdata)

			switch c := con.(type) {
			case *pango.Firewall:
				vals := url.Values{}
				vals.Set("key", c.ApiKey)
				vals.Set("type", "import")
				vals.Set("category", "configuration")

				_, err = c.CommunicateFile(fcontent, fbase, fp, vals, nil)
				if err != nil {
					log.Printf("Failed to transfer file to device: %s", err)
				}

				if load == true {
					cmd := fmt.Sprintf("<load><config><from>%s</from></config></load>", fbase)
					vals := url.Values{}
					vals.Set("key", c.ApiKey)
					_, err = c.Op(cmd, "", vals, nil)
					if err != nil {
						log.Printf("Error loading file: %s", err)
					}
				}
			case *pango.Panorama:
				vals := url.Values{}
				vals.Set("key", c.ApiKey)
				vals.Set("type", "import")
				vals.Set("category", "configuration")

				_, err = c.CommunicateFile(fcontent, fbase, fp, vals, nil)
				if err != nil {
					log.Printf("Failed to transfer file to device: %s", err)
				}

				if load == true {
					cmd := fmt.Sprintf("<load><config><from>%s</from></config></load>", fbase)
					vals := url.Values{}
					vals.Set("key", c.ApiKey)
					_, err = c.Op(cmd, "", vals, nil)
					if err != nil {
						log.Printf("Error loading file: %s", err)
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(provisionCmd)
	provisionCmd.Flags().StringVarP(&source, "source", "s", "", "Source of config - ironskillet, local or remote")
	provisionCmd.Flags().StringVarP(&ostype, "os", "o", "", "Device OS - only used when source is ironskillet; panos or panorama")
	// provisionCmd.Flags().StringVarP(&ver, "version", "v", "", "IronSkillet template version to use - 8.0, 8.1, 9.0")
	provisionCmd.Flags().StringVarP(&config, "config", "c", "static", "IronSkillet config to use - aws|azure|gcp|dhcp|static")
	provisionCmd.Flags().StringVarP(&user, "user", "u", "", "User to connect to the device as")
	provisionCmd.Flags().StringVarP(&device, "device", "d", "", "Device to connect to")
	provisionCmd.Flags().StringVarP(&fh, "file", "f", "", "Name of the XML config file to use - only used when source is local or remote")
	provisionCmd.Flags().BoolVarP(&load, "load", "l", false, "Load the file into the running-config - use with caution")
	provisionCmd.MarkFlagRequired("user")
	provisionCmd.MarkFlagRequired("device")
	provisionCmd.MarkFlagRequired("source")
}
