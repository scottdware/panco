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
	"log"
	"os"

	easycsv "github.com/scottdware/go-easycsv"
	"github.com/spf13/cobra"
)

var exampleCreate = `# Lines starting with a hashtag will be ignored,,,,,
,,,,,
# CREATE ADDRESS OBJECTS,,,,,
#Name,Type,Value,Description,Tag,Device Group/Vsys
IP-10.1.1.1,ip,10.1.1.1,,,vsys1
DNS-Object,fqdn,paloaltonetworks.com,,,shared
Net-10.1.1.0_24,ip,10.1.1.0/24,,,Corporate_HQ
DHCP_Range,range,10.1.1.200-10.1.1.250,,,Branch-Office
,,,,,
# CREATE SERVICE OBJECTS,,,,,
#Name,Type,Value,Description,Tag,Device Group/Vsys
TCP_9999,tcp,9999,,,vsys1
UDP_Range-30000-39999,udp,30000-39999,,,shared
Web-Ports,tcp,"80, 443, 8080",,,Corporate_HQ
UDP_7777,udp,7777,,,Branch-Office
,,,,,
# CREATE OBJECTS WITH TAGS - tags MUST be pre-existing! Multiple tags separated with a comma or semicolon (enclosed in quotes),,,,,
#Name,Type,Value,Description,Tag,Device Group/Vsys
File-Server,ip,10.1.1.2,,"Server, DMZ",shared
DMZ_Web-Server,ip,10.2.2.2,,Public,vsys1
,,,,,
# CREATE ADDRESS OR SERVICE GROUPS,,,,,
# When creating address or service groups the members must already exist. The best way to do this,,,,,
# is to create the address or service objects first (earlier in the spreadsheet) then reference them in the group.,,,,,
#Name,Type,Value,Description,Tag,Device Group/Vsys
Corporate_Devices,static,"Net-10.1.1.0_24, File-Server, DHCP_Range",,,shared
Dynamic-Servers,dynamic,"'Server' or 'DMZ'",,,vsys1
Allowed_Services,service,"TCP_9999, Web-Ports, UDP_7777",,,Corporate_HQ
,,,,,
# ADD OBJECTS TO A GROUP,,,,,
# If you have pre-existing groups and just want to add objects to them (similar to creating groups above),,,,,
#Name,Type,Value,Description,Tag,Device Group/Vsys
Group_to_AddTo,static,"IP-10.2.2.2, IP-192.168.0.1",,,shared`

// exampleCmd represents the example command
var exampleCmd = &cobra.Command{
	Use:   "example",
	Short: "Create example CSV files for import reference",
	Long: `This command will create a sample CSV file for use with the
import command. The files will be placed in the location where you are running
the command from, and are named as such:
	
panco-example-import.csv`,
	Run: func(cmd *cobra.Command, args []string) {
		createFh, err := easycsv.NewCSV("panco-example-import.csv")
		if err != nil {
			log.Printf("Failed to create example CSV file: %s", err)
			os.Exit(1)
		}

		createFh.Write(exampleCreate)
		createFh.End()
	},
}

func init() {
	rootCmd.AddCommand(exampleCmd)
}
