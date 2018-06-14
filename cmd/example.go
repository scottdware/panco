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

	easycsv "github.com/scottdware/go-easycsv"
	"github.com/spf13/cobra"
)

var exampleCreate = `# Lines starting with a hashtag will be ignored,,,,,
,,,,,
# Creating address objects,,,,,
host_10.1.1.1,ip,10.1.1.1,,,shared
net_10.0.0.0_8,ip,10.0.0.0/8,,,shared
host_172.16.1.1-172.16.1.10,range,172.16.1.1-172.16.1.10,,,shared
host_microsoft.com,fqdn,microsoft.com,,,Some_Device_Group
,,,,,
# Creating service objects,,,,,
tcp_port_9999,tcp,9999,,,shared
udp_ports_30000-39999,udp,30000-39999,,,Some_Device_Group
web_ports,tcp,"80, 8080, 443",,,shared
udp_port_7777,udp,7777,,,shared
,,,,,
# Creating objects with tags - tags must be pre-existing!,,,,,
host_192.168.1.0_24,ip,192.168.1.0/24,,Server-Network,shared
host_192.168.1.20,ip,192.168.1.20,,Server,shared
,,,,,
# Creating address and service groups,,,,,
# When creating address or service groups the members must already exist. The best way to do this,,,,,
# is to create the address or service objects first (earlier in the spreadsheet) then reference them in the group.,,,,,
Inside_Sources,static,host_10.1.1.1 net_10.0.0.0_8,,,shared
server_networks_dynamic,dynamic,Server-Network or Servers,,,shared
tcp_port_group,service,tcp_port_9999,,,shared`
var exampleModify = `# Lines starting with a hashtag will be ignored,,,
,,,
# Adding objects to groups,,,
address,add,host_172.16.1.1-172.16.1.10,Inside_Sources
service,add,udp_port_7777,tcp_port_group
service,add,web_ports,tcp_port_group
,,,
# Removing objects from groups,,,
service,remove,udp_port_7777,tcp_port_group
address,remove,host_172.16.1.1-172.16.1.10,Inside_Sources`

// exampleCmd represents the example command
var exampleCmd = &cobra.Command{
	Use:   "example",
	Short: "Create example CSV files for import reference",
	Long: `This command will create two sample, reference CSV files for use with the
import command. The files will be placed in the location where you are running
the command from, and are named as such:
	
example-create.csv and example-modify.csv`,
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
		}

		modifyFh.Write(exampleModify)
		modifyFh.End()
	},
}

func init() {
	rootCmd.AddCommand(exampleCmd)
}
