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

// sessionsCmd represents the sessions command
var sessionsCmd = &cobra.Command{
	Use:   "sessions",
	Short: "Query the session table on a firewall, and export it to a CSV file",
	Long: `This command will dump the entire session table on a firewall to the CSV file
that you specify. You can optionally define a filter, and use the same criteria as you would
on the command line. The filter query must be enclosed in quotes "", and the format is:

option=value (e.g. --query "application=ssl")

Your filter can include multiple items, and each group must be separated by a comma, e.g.:

--query "application=ssl, ssl-decrypt=yes, protocol=tcp"

Depending on the number of sessions, the export could take some time.`,
	Run: func(cmd *cobra.Command, args []string) {
		var sessions *panos.SessionTable
		pass := passwd()
		creds := &panos.AuthMethod{
			Credentials: []string{user, pass},
		}

		pan, err := panos.NewSession(device, creds)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		sessions, err = pan.Sessions()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(query) > 0 {
			sessions, err = pan.Sessions(query)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		if !strings.Contains(fh, ".csv") {
			fh = fmt.Sprintf("%s.csv", fh)
		}

		csv, err := easycsv.NewCSV(fh)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		csv.Write("StartTime,From,To,SourceAddress,State,DestinationAddress,DecryptMirror,Proxy")
		csv.Write(",SourcePort,DestinationPort,Protocol,Application,SecurityRule,IngressInterface")
		csv.Write(",EgressInterface,TotalByteCount,Vsys")
		csv.Write(",Type,VsysID,NAT,SourceNAT,NATSourceAddress,NATSourcePort,DestinationNAT,NATDestinationAddress,NATDestinationPort")
		csv.Write(",ID,Flags\n")

		time.Sleep(50 * time.Millisecond)

		for _, s := range sessions.Sessions {
			csv.Write(fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%d,%d,%d,%s,%s,%s,%s,%d,%s",
				s.StartTime, s.From, s.To, s.SourceAddress, s.State, s.DestinationAddress, s.DecryptMirror, s.Proxy,
				s.SourcePort, s.DestinationPort, s.Protocol, s.Application, s.SecurityRule, s.IngressInterface,
				s.EgressInterface, s.TotalByteCount, s.Vsys))
			csv.Write(fmt.Sprintf(",%s,%d,%s,%s,%s,%d,%s,%s,%d,%d,\"%s\"\n",
				s.Type, s.VsysID, s.NAT, s.SourceNAT, s.NATSourceAddress, s.NATSourcePort, s.DestinationNAT,
				s.NATDestinationAddress, s.NATDestinationPort, s.ID, s.Flags))

			time.Sleep(10 * time.Millisecond)
		}

		csv.End()
	},
}

func init() {
	rootCmd.AddCommand(sessionsCmd)

	sessionsCmd.Flags().StringVarP(&user, "user", "u", "", "User to connect to the device as")
	sessionsCmd.Flags().StringVarP(&device, "device", "d", "", "Firewall or Panorama device to connect to")
	sessionsCmd.Flags().StringVarP(&fh, "file", "f", "", "Name of the CSV file to export the session table to")
	sessionsCmd.Flags().StringVarP(&query, "query", "q", "", "Filter string to include sessions that only matching the criteria")
	sessionsCmd.MarkFlagRequired("user")
	sessionsCmd.MarkFlagRequired("device")
	sessionsCmd.MarkFlagRequired("export")
}
