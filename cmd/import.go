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

	panos "github.com/scottdware/go-panos"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import CSV files that will create and/or modify objects",
	Long: `This command, given the spcific flag, will create or modify address and/or
service objects based on the information you have provided in your CSV file(s).

Run 'panco example' to create a sample CSV file for each action to use as reference.`,
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

		if create != "" {
			err = pan.CreateObjectsFromCsv(create)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		if modify != "" {
			err = pan.ModifyGroupsFromCsv(modify)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().StringVarP(&create, "create", "c", "", "Name of the CSV file to create objects with")
	importCmd.Flags().StringVarP(&modify, "modify", "m", "", "Name of the CSV file to modify groups with")
	importCmd.Flags().StringVarP(&user, "user", "u", "", "User to connect to the device as")
	importCmd.Flags().StringVarP(&device, "device", "d", "", "Firewall or Panorama device to connect to")
	importCmd.MarkFlagRequired("user")
	importCmd.MarkFlagRequired("device")
}
