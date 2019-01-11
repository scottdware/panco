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

	easycsv "github.com/scottdware/go-easycsv"
	"github.com/spf13/cobra"
)

var exampleCreate = `# Lines starting with a hashtag will be ignored,,,,,
,,,,,
# CREATE ADDRESS OBJECTS,,,,,
#Name,Type,Value,Description,Tag,Device-group
Stormtrooper,ip,10.0.0.50,,,shared
Boba_Fett,fqdn,mandalorianmercs.org,,,shared
Greedo,ip,10.0.0.53,,,shared
net_10.0.0.0_24,ip,10.0.0.0/24,,,shared
DHCP_Range,range,10.0.0.200-10.0.0.250,,,shared
fqdn_starwars.com,fqdn,starwars.com,,,shared
,,,,,
# CREATE SERVICE OBJECTS,,,,,
#Name,Type,Value,Description,Tag,Device-group
tcp_port_9999,tcp,9999,,,shared
udp_ports_30000-39999,udp,30000-39999,,,shared
web_ports,tcp,"80, 8080, 443",,,shared
udp_port_7777,udp,7777,,,shared
,,,,,
# CREATE OBJECTS WITH TAGS - tags must be pre-existing! Multiple tags separated with a comma (enclosed in quotes),,,,,
#Name,Type,Value,Description,Tag,Device-group
Millenium_Falcon,ip,10.0.0.51,,"Han, Chewy",shared
Death Trooper,ip,10.0.0.52,,Death Star,shared
,,,,,
# CREATE ADDRESS AND SERVICE GROUPS,,,,,
# When creating address or service groups the members must already exist. The best way to do this,,,,,
# is to create the address or service objects first (earlier in the spreadsheet) then reference them in the group.,,,,,
#Name,Type,Value,Description,Tag,Device-group
Cantina,static,"Millenium_Falcon, net_10.0.0.0_24, DHCP_Range",,,shared
The_Dark_Side,dynamic,"'Death Star' or 'Sith'",,,shared
tcp_port_group,service,"tcp_port_9999, udp_port_7777",,,shared`
var exampleModify = `# Lines starting with a hashtag will be ignored,,,
,,,
# ADDING OBJECTS TO GROUPS,,,
#Group Type,Action,Object,Group,Device-group
address,add,Greedo,Cantina
service,add,web_ports,tcp_port_group
,,,
# REMOVING OBJECTS FROM GROUPS,,,
#Group Type,Action,Object,Group,Device-group
service,remove,udp_port_7777,tcp_port_group
address,remove,Boba_Fett,Cantina`
var exampleRename = `# Lines starting with a hashtag will be ignored,,
,,
# RENAMING OBJECTS,,
#Old name,New name,device-group
Boba_Fett,Bounty-Hunter,
Stormtrooper,Cannon_Fodder,
Greedo,Han_Shot_Me_First,
,,
# If renaming non-shared objects on Panorama - specify the device-group in the 3rd column,,
Apprentice,Sith-Lord,Star_Wars
Padawan,Jedi_Master,Star_Wars`

// exampleCmd represents the example command
var exampleCmd = &cobra.Command{
	Use:   "example",
	Short: "Create example CSV files for import reference",
	Long: `This command will create three sample, reference CSV files for use with the
import command. The files will be placed in the location where you are running
the command from, and are named as such:
	
example-create.csv
example-modify.csv
example-rename.csv`,
	Run: func(cmd *cobra.Command, args []string) {
		createFh, err := easycsv.NewCSV("example-create.csv")
		if err != nil {
			fmt.Println(err)
		}

		createFh.Write(exampleCreate)
		createFh.End()

		modifyFh, err := easycsv.NewCSV("example-modify.csv")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		modifyFh.Write(exampleModify)
		modifyFh.End()

		renameFh, err := easycsv.NewCSV("example-rename.csv")
		if err != nil {
			fmt.Println(err)
		}

		renameFh.Write(exampleRename)
		renameFh.End()
	},
}

func init() {
	rootCmd.AddCommand(exampleCmd)
}
