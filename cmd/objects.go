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

// objectsCmd represents the objects command
var objectsCmd = &cobra.Command{
	Use:   "objects",
	Short: "Import/export address and service objects, rename objects, and modify groups",
	Long: `This command allows you to perform the following actions on address and service
objects: export, import, rename, and modify groups. When you select the export option (--action export),
there are two files that will be created. One will hold all of the address objects, and the other
will hold all of the service objects.

When exporting and run against a Panorama device without specifying the --devicegroup flag, all objects will be
exported, including shared ones. Importing objects into Panorama without specifying the --devicegroup flag does
not matter.

The rename action allows you to rename address, service and tag objects.

Using the modify action, allows you to add or remove objects from groups, based on the data you have
within your CSV file.

Please see "panco example" for sample CSV files to use as a reference.`,
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
			fh = strings.TrimSuffix(fh, ".csv")

			addressCSV, _ := easycsv.NewCSV(fmt.Sprintf("%s_addr.csv", fh))
			serviceCSV, _ := easycsv.NewCSV(fmt.Sprintf("%s_svcs.csv", fh))

			addressCSV.Write("# ADDRESS OBJECTS\n")
			addressCSV.Write("#Name,Type,Value,Description,Tag,DeviceGroup\n")

			addrs, err := pan.Addresses(dg)
			if err != nil {
				fmt.Println(err)
			}

			for _, a := range addrs.Addresses {
				addrObjTag := sliceToString(a.Tag)

				if len(a.IPAddress) != 0 {
					if pan.DeviceType == "panorama" && len(dg) > 0 {
						addressCSV.Write(fmt.Sprintf("\"%s\",%s,%s,\"%s\",\"%s\",%s\n",
							a.Name, "ip", a.IPAddress, a.Description, addrObjTag, dg))
					}

					if pan.DeviceType == "panorama" && len(dg) == 0 {
						addressCSV.Write(fmt.Sprintf("\"%s\",%s,%s,\"%s\",\"%s\",%s\n",
							a.Name, "ip", a.IPAddress, a.Description, addrObjTag, "shared"))
					}

					if pan.DeviceType == "panos" {
						addressCSV.Write(fmt.Sprintf("\"%s\",%s,%s,\"%s\",\"%s\",%s\n",
							a.Name, "ip", a.IPAddress, a.Description, addrObjTag, ""))
					}

					time.Sleep(5 * time.Millisecond)
				}

				if len(a.IPRange) != 0 {
					if pan.DeviceType == "panorama" && len(dg) > 0 {
						addressCSV.Write(fmt.Sprintf("\"%s\",%s,%s,\"%s\",\"%s\",%s\n",
							a.Name, "range", a.IPRange, a.Description, addrObjTag, dg))
					}

					if pan.DeviceType == "panorama" && len(dg) == 0 {
						addressCSV.Write(fmt.Sprintf("\"%s\",%s,%s,\"%s\",\"%s\",%s\n",
							a.Name, "range", a.IPRange, a.Description, addrObjTag, "shared"))
					}

					if pan.DeviceType == "panos" {
						addressCSV.Write(fmt.Sprintf("\"%s\",%s,%s,\"%s\",\"%s\",%s\n",
							a.Name, "range", a.IPRange, a.Description, addrObjTag, ""))
					}

					time.Sleep(5 * time.Millisecond)
				}

				if len(a.FQDN) != 0 {
					if pan.DeviceType == "panorama" && len(dg) > 0 {
						addressCSV.Write(fmt.Sprintf("\"%s\",%s,%s,\"%s\",\"%s\",%s\n",
							a.Name, "fqdn", a.FQDN, a.Description, addrObjTag, dg))
					}

					if pan.DeviceType == "panorama" && len(dg) == 0 {
						addressCSV.Write(fmt.Sprintf("\"%s\",%s,%s,\"%s\",\"%s\",%s\n",
							a.Name, "fqdn", a.FQDN, a.Description, addrObjTag, "shared"))
					}

					if pan.DeviceType == "panos" {
						addressCSV.Write(fmt.Sprintf("\"%s\",%s,%s,\"%s\",\"%s\",%s\n",
							a.Name, "fqdn", a.FQDN, a.Description, addrObjTag, ""))
					}

					time.Sleep(5 * time.Millisecond)
				}
			}

			addressCSV.Write("#\n")
			addressCSV.Write("# ADDRESS GROUPS\n")
			addressCSV.Write("#Name,Type,Value,Description,Tag,DeviceGroup\n")

			groups, err := pan.AddressGroups(dg)
			if err != nil {
				fmt.Println(err)
			}

			for _, g := range groups.Groups {
				addrGrpMembers := sliceToString(g.Members)
				addrGrpTag := sliceToString(g.Tag)

				if g.Type == "Static" {
					if pan.DeviceType == "panorama" && len(dg) > 0 {
						addressCSV.Write(fmt.Sprintf("\"%s\",%s,\"%s\",\"%s\",\"%s\",%s\n",
							g.Name, strings.ToLower(g.Type), addrGrpMembers, g.Description, addrGrpTag, dg))
					}

					if pan.DeviceType == "panorama" && len(dg) == 0 {
						addressCSV.Write(fmt.Sprintf("\"%s\",%s,\"%s\",\"%s\",\"%s\",%s\n",
							g.Name, strings.ToLower(g.Type), addrGrpMembers, g.Description, addrGrpTag, "shared"))
					}

					if pan.DeviceType == "panos" {
						addressCSV.Write(fmt.Sprintf("\"%s\",%s,\"%s\",\"%s\",\"%s\",%s\n",
							g.Name, strings.ToLower(g.Type), addrGrpMembers, g.Description, addrGrpTag, ""))
					}

					time.Sleep(5 * time.Millisecond)
				}

				if g.Type == "Dynamic" {
					if pan.DeviceType == "panorama" && len(dg) > 0 {
						addressCSV.Write(fmt.Sprintf("\"%s\",%s,\"%s\",\"%s\",\"%s\",%s\n",
							g.Name, strings.ToLower(g.Type), g.DynamicFilter, g.Description, addrGrpTag, dg))
					}

					if pan.DeviceType == "panorama" && len(dg) == 0 {
						addressCSV.Write(fmt.Sprintf("\"%s\",%s,\"%s\",\"%s\",\"%s\",%s\n",
							g.Name, strings.ToLower(g.Type), g.DynamicFilter, g.Description, addrGrpTag, "shared"))
					}

					if pan.DeviceType == "panos" {
						addressCSV.Write(fmt.Sprintf("\"%s\",%s,\"%s\",\"%s\",\"%s\",%s\n",
							g.Name, strings.ToLower(g.Type), g.DynamicFilter, g.Description, addrGrpTag, ""))
					}

					time.Sleep(5 * time.Millisecond)
				}
			}

			addressCSV.End()

			time.Sleep(1 * time.Second)

			serviceCSV.Write("# SERVICE OBJECTS\n")
			serviceCSV.Write("#Name,Type,Value,Description,Tag,DeviceGroup\n")

			svcs, err := pan.Services(dg)
			if err != nil {
				fmt.Println(err)
			}

			for _, s := range svcs.Services {
				svcObjTag := sliceToString(s.Tag)

				if len(s.TCPPort) != 0 {
					if pan.DeviceType == "panorama" && len(dg) > 0 {
						serviceCSV.Write(fmt.Sprintf("\"%s\",%s,\"%s\",\"%s\",\"%s\",%s\n",
							s.Name, "tcp", s.TCPPort, s.Description, svcObjTag, dg))
					}

					if pan.DeviceType == "panorama" && len(dg) == 0 {
						serviceCSV.Write(fmt.Sprintf("\"%s\",%s,\"%s\",\"%s\",\"%s\",%s\n",
							s.Name, "tcp", s.TCPPort, s.Description, svcObjTag, "shared"))
					}

					if pan.DeviceType == "panos" {
						serviceCSV.Write(fmt.Sprintf("\"%s\",%s,\"%s\",\"%s\",\"%s\",%s\n",
							s.Name, "tcp", s.TCPPort, s.Description, svcObjTag, ""))
					}

					time.Sleep(5 * time.Millisecond)
				}

				if len(s.UDPPort) != 0 {
					if pan.DeviceType == "panorama" && len(dg) > 0 {
						serviceCSV.Write(fmt.Sprintf("\"%s\",%s,\"%s\",\"%s\",\"%s\",%s\n",
							s.Name, "udp", s.UDPPort, s.Description, svcObjTag, dg))
					}

					if pan.DeviceType == "panorama" && len(dg) == 0 {
						serviceCSV.Write(fmt.Sprintf("\"%s\",%s,\"%s\",\"%s\",\"%s\",%s\n",
							s.Name, "udp", s.UDPPort, s.Description, svcObjTag, "shared"))
					}

					if pan.DeviceType == "panos" {
						serviceCSV.Write(fmt.Sprintf("\"%s\",%s,\"%s\",\"%s\",\"%s\",%s\n",
							s.Name, "udp", s.UDPPort, s.Description, svcObjTag, ""))
					}

					time.Sleep(5 * time.Millisecond)
				}
			}

			serviceCSV.Write("#\n")
			serviceCSV.Write("# SERVICE GROUPS\n")
			serviceCSV.Write("#Name,Type,Value,Description,Tag,DeviceGroup\n")

			svcg, err := pan.ServiceGroups(dg)
			if err != nil {
				fmt.Println(err)
			}

			for _, sg := range svcg.Groups {
				svcGrpMembers := sliceToString(sg.Members)
				svcGrpTag := sliceToString(sg.Tag)

				if pan.DeviceType == "panorama" && len(dg) > 0 {
					serviceCSV.Write(fmt.Sprintf("\"%s\",%s,\"%s\",\"%s\",\"%s\",%s\n",
						sg.Name, "service", svcGrpMembers, sg.Description, svcGrpTag, dg))
				}

				if pan.DeviceType == "panorama" && len(dg) == 0 {
					serviceCSV.Write(fmt.Sprintf("\"%s\",%s,\"%s\",\"%s\",\"%s\",%s\n",
						sg.Name, "service", svcGrpMembers, sg.Description, svcGrpTag, "shared"))
				}

				if pan.DeviceType == "panos" {
					serviceCSV.Write(fmt.Sprintf("\"%s\",%s,\"%s\",\"%s\",\"%s\",%s\n",
						sg.Name, "service", svcGrpMembers, sg.Description, svcGrpTag, ""))
				}

				time.Sleep(5 * time.Millisecond)
			}

			serviceCSV.End()
		}

		if action == "import" {
			err = pan.CreateObjectsFromCsv(fh)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		if action == "modify" {
			err = pan.ModifyGroupsFromCsv(fh)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		if action == "rename" {
			objs, err := easycsv.Open(fh)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			for _, line := range objs {
				var devicegroup string
				linelen := len(line)
				oldname := line[0]
				newname := line[1]

				if linelen > 2 && len(line[2]) > 0 {
					devicegroup = line[2]
				}

				if err = pan.RenameObject(oldname, newname, devicegroup); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				time.Sleep(20 * time.Millisecond)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(objectsCmd)

	objectsCmd.Flags().StringVarP(&action, "action", "a", "", "Action to perform - export, import, rename, or modify")
	objectsCmd.Flags().StringVarP(&fh, "file", "f", "", "Name of the CSV file to export/import or modify")
	objectsCmd.Flags().StringVarP(&dg, "devicegroup", "g", "", "Device group - only needed when exporting and run against a Panorama device")
	objectsCmd.Flags().StringVarP(&user, "user", "u", "", "User to connect to the device as")
	objectsCmd.Flags().StringVarP(&device, "device", "d", "", "Firewall or Panorama device to connect to")
	objectsCmd.MarkFlagRequired("user")
	objectsCmd.MarkFlagRequired("device")
	objectsCmd.MarkFlagRequired("action")
	objectsCmd.MarkFlagRequired("file")
}
