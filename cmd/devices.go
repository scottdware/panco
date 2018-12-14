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

	easycsv "github.com/scottdware/go-easycsv"
	panos "github.com/scottdware/go-panos"
	"github.com/spf13/cobra"
)

// devicesCmd represents the devices command
var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Device specific functions such as exporting data from Panorama or local firewalls",
	Long: `The devices command will provide information about devices connected/managed
by Panorama, as well as other device (firewall) specific information. Currently, the only action
is to "export"

When ran against a Panorama device, the following "types" are available: applications and devices.
Using "devices" will export a list of all managed firewalls, with data from the "Panorama > Managed Devices"
tab. When using "applications" it will export a list of every predefined application and all of
it's characteristics, such as category, subcategory, etc.

Using the "interfaces" type will export a list of all of the logical and physical (hardware)
interfaces on the device, with all of their information, such as IP address, MAC address, zone, etc.`,
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

		if action == "export" && t == "devices" {
			devs, err := pan.Devices()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if !strings.Contains(fh, ".csv") {
				fh = fmt.Sprintf("%s.csv", fh)
			}

			dfh, err := easycsv.NewCSV(fh)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			dfh.Write("Hostname,Model,Family,Serial,IP Address,Connected,Connected At,Operational Mode,")
			dfh.Write("HA State,HA Peer,Software Version,Apps and Threat,Antivirus,URL Filtering,Wildfire,Uptime\n")
			for _, device := range devs.Devices {
				dfh.Write(fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,\"%s\"\n",
					device.Hostname, device.Model, device.Family, device.Serial, device.IPAddress, device.Connected,
					device.ConnectedAt, device.OperationalMode, device.HAState, device.HAPeer,
					device.SoftwareVersion, device.AppVersion, device.AntiVirusVersion, device.URLFilteringVersion,
					device.WildfireVersion, device.Uptime))
			}

			dfh.End()
		}

		if action == "export" && t == "applications" {
			apps, err := pan.ApplicationInfo()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if !strings.Contains(fh, ".csv") {
				fh = fmt.Sprintf("%s.csv", fh)
			}

			afh, err := easycsv.NewCSV(fh)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			afh.Write("#ID,Name,OriginCountry,OriginLanguage,Category,Subcategory,Technology,VirusIdent,FileTypeIdent,TunnelOtherApplications,")
			afh.Write("DataIdent,FileForward,IsSAAS,EvasiveBehavior,ConsumeBigBandwidth,UsedByMalware,AbleToTransferFile,HasKnownVulnerability,")
			afh.Write("ProneToMisuse,PervasiveUse,References,DefaultPort,UseApplications,ImplicitUseApplications,Risk\n")
			for _, app := range apps.Applications {
				var references []string
				for _, r := range app.References {
					ref := fmt.Sprintf("%s: %s", r.Name, r.Link)
					references = append(references, ref)
				}

				afh.Write(fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,\"%s\",\"%s\",\"%s\",\"%s\",%d\n", app.ID,
					app.Name, app.OriginCountry, app.OriginLanguage, app.Category, app.Subcategory, app.Technology, app.VirusIdent,
					app.FileTypeIdent, app.TunnelOtherApplications, app.DataIdent, app.FileForward, app.IsSAAS, app.EvasiveBehavior,
					app.ConsumeBigBandwidth, app.UsedByMalware, app.AbleToTransferFile, app.HasKnownVulnerability, app.ProneToMisuse, app.PervasiveUse,
					sliceToString(references), sliceToString(app.DefaultPort), sliceToString(app.UseApplications),
					sliceToString(app.ImplicitUseApplications), app.Risk))
			}

			afh.End()
		}

		if action == "export" && t == "interfaces" {
			ifinfo, err := pan.InterfaceInfo()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if !strings.Contains(fh, ".csv") {
				fh = fmt.Sprintf("%s.csv", fh)
			}

			ifh, err := easycsv.NewCSV(fh)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			ifh.Write("#LOGICAL INTERFACE INFORMATION\n")
			ifh.Write("#Name, Zone, VirtualRouter, Vsys, DynamicAddress, IPv6Addres, Tag, IPAddress, ID, Address\n")
			for _, l := range ifinfo.Logical {
				ifh.Write(fmt.Sprintf("%s,%s,%s,%d,%s,%s,%d,%s,%d,%s\n", l.Name, l.Zone, l.VirtualRouter, l.Vsys, l.DynamicAddress,
					l.IPv6Addres, l.Tag, l.IPAddress, l.ID, l.Address))
			}

			ifh.Write("#HARDWARE INTERFACE INFORMATION\n")
			ifh.Write("#Name, Duplex, Type, State, AEMember, Settings, MACAddress, Mode, Speed, ID\n")
			for _, h := range ifinfo.Hardware {
				ifh.Write(fmt.Sprintf("%s,%s,%s,%s,\"%s\",%s,%s,%s,%s,%d\n", h.Name, h.Duplex, h.Type, h.State,
					sliceToString(h.AEMember), h.Settings, h.MACAddress, h.Mode, h.Speed, h.ID))
			}

			ifh.End()
		}
	},
}

func init() {
	rootCmd.AddCommand(devicesCmd)
	devicesCmd.Flags().StringVarP(&action, "action", "a", "", "Action to perform - e.g. export")
	devicesCmd.Flags().StringVarP(&user, "user", "u", "", "User to connect to the device as")
	devicesCmd.Flags().StringVarP(&device, "device", "d", "", "Firewall or Panorama device to connect to")
	devicesCmd.Flags().StringVarP(&fh, "file", "f", "", "Name of the CSV file to export to")
	devicesCmd.Flags().StringVarP(&t, "type", "t", "", "Type of information to export - e.g. applications or interfaces")
	devicesCmd.MarkFlagRequired("user")
	devicesCmd.MarkFlagRequired("device")
	devicesCmd.MarkFlagRequired("action")
	devicesCmd.MarkFlagRequired("type")
	devicesCmd.MarkFlagRequired("file")
}
