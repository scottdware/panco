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
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
)

type release struct {
	URL    string  `json:"html_url"`
	Tag    string  `json:"tag_name"`
	Assets []asset `json:"assets"`
}

type asset struct {
	Name    string `json:"name"`
	URL     string `json:"browser_download_url"`
	Created string `json:"created_at"`
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version information for panco",
	Long: `Shows the current version of panco you are running, and checks to see if there is a
newer version available.`,
	Run: func(cmd *cobra.Command, args []string) {
		var releases []release
		curver := "v2022.11"

		resp, err := resty.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Accept", "application/vnd.github.v3+json").
			Get("https://api.github.com/repos/scottdware/panco/releases")

		if err != nil {
			fmt.Printf("unable to connect to Github - %s", err)
		}

		if err := json.Unmarshal([]byte(resp.String()), &releases); err != nil {
			fmt.Printf("JSON parse error on release info - %s", err)
		}

		latestver := releases[0].Tag
		download := releases[0].URL

		if curver == latestver {
			fmt.Printf("You are running the latest version of panco - %s\n\nSee https://panco.dev for complete documentation", curver)
		} else {
			fmt.Printf("New version available - %s! Download here: %s\n\nSee https://panco.dev for complete documentation", latestver, download)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
