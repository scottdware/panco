/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"

	easycsv "github.com/scottdware/go-easycsv"
	"github.com/spf13/cobra"
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Generate CSV templates for object and policy importing",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		switch t {
		case "objects":
			output, err := easycsv.NewCSV("Objects.csv")
			if err != nil {
				log.Printf("Error creating CSV file - %s", err)
				os.Exit(1)
			}

			output.Write("#Name,Type,Value,Description,Tag,Device-Group/Vsys\n")
			output.Write("IP_Object,ip,10.1.1.1,,,vsys1\n")
			output.Write("FQDN_Object,fqdn,paloaltonetworks.com,,,Some_Device_Group\n")
			output.Write("Range_Object,range,10.1.1.1-10.1.1.10,,,Some_Device_Group\n")
			output.Write("TCP_443,tcp,443,,,vsys1\n")
			output.Write("UDP_53,udp,53,,,vsys1\n")
			output.Write("Servers,tag,Orange,,,Some_Device_Group\n")
			output.Write("DMZ,tag,Yellow,,,Some_Device_Group\n")
			output.Write("Address_Group,static,\"IP_Object, Range_Object\",,,Some_Device_Group\n")
			output.Write("Dynamic_Group,dynamic,\"'Servers' and 'DMZ'\",,,Some_Device_Group\n")
			output.Write("Service_Group,service,\"TCP_443, UDP_53\",,,vsys1\n")
			output.Write("Add_to_Custom_URL,urladd,google.com/,,,vsys1\n")
			output.Write("Remove_from_Custom_URL,urlremove,bing.com/,,,vsys1\n")
			output.Write("Object_to-be_Deleted,ip,delete,,,Some_Device_Group\n")
			output.Write("Remove_Object_from_This_Group,remove-address,IP_Object,,,vsys1\n")
			output.Write("Remove_Service_from_This_Group,remove-service,UDP_53,,,vsys1\n")
			output.Write("Change_Me_Address,rename-address,New_Address_Name,,,Some_Device_Group\n")
			output.Write("Change_Me_AddrGroup,rename-addressgroup,New_AddrGroup_Name,,,Some_Device_Group\n")
			output.Write("Change_Me_Service,rename-service,New_Service_Name,,,vsys1\n")
			output.Write("Change_Me_ServiceGroup,rename-servicegroup,New_ServiceGroup_Name,,,vsys1\n")
			output.Write("Object_to_Tag,ip,10.1.1.1,,Servers,Some_Device_Group\n")
			output.End()
		case "policy":
			output, err := easycsv.NewCSV("Policy.csv")
			if err != nil {
				log.Printf("Error creating CSV file - %s", err)
				os.Exit(1)
			}
			output.Write("#Name,Type,Description,Tags,SourceZones,SourceAddresses,NegateSource,SourceUsers,HipProfiles,DestinationZones,DestinationAddresses,NegateDestination,Applications,Services,Categories,Action,LogSetting,LogStart,LogEnd,Disabled,Schedule,IcmpUnreachable,DisableServerResponseInspection,Group,Virus,Spyware,Vulnerability,UrlFiltering,FileBlocking,WildFireAnalysis,DataFiltering\n")
			output.Write(",universal,,,,,FALSE,any,,,,FALSE,any,,any,allow,default,FALSE,TRUE,FALSE,,FALSE,FALSE,,,,,,,,\n")
			output.End()
		case "all":
			obj, err := easycsv.NewCSV("Objects.csv")
			if err != nil {
				log.Printf("Error creating CSV file - %s", err)
				os.Exit(1)
			}

			obj.Write("#Name,Type,Value,Description,Tag,Device-Group/Vsys\n")
			obj.Write("IP_Object,ip,10.1.1.1,,,vsys1\n")
			obj.Write("FQDN_Object,fqdn,paloaltonetworks.com,,,Some_Device_Group\n")
			obj.Write("Range_Object,range,10.1.1.1-10.1.1.10,,,Some_Device_Group\n")
			obj.Write("TCP_443,tcp,443,,,vsys1\n")
			obj.Write("UDP_53,udp,53,,,vsys1\n")
			obj.Write("Servers,tag,Orange,,,Some_Device_Group\n")
			obj.Write("DMZ,tag,Yellow,,,Some_Device_Group\n")
			obj.Write("Address_Group,static,\"IP_Object, Range_Object\",,,Some_Device_Group\n")
			obj.Write("Dynamic_Group,dynamic,\"'Servers' and 'DMZ'\",,,Some_Device_Group\n")
			obj.Write("Service_Group,service,\"TCP_443, UDP_53\",,,vsys1\n")
			obj.Write("Add_to_Custom_URL,urladd,google.com/,,,vsys1\n")
			obj.Write("Remove_from_Custom_URL,urlremove,bing.com/,,,vsys1\n")
			obj.Write("Object_to-be_Deleted,ip,delete,,,Some_Device_Group\n")
			obj.Write("Remove_Object_from_This_Group,remove-address,IP_Object,,,vsys1\n")
			obj.Write("Remove_Service_from_This_Group,remove-service,UDP_53,,,vsys1\n")
			obj.Write("Change_Me_Address,rename-address,New_Address_Name,,,Some_Device_Group\n")
			obj.Write("Change_Me_AddrGroup,rename-addressgroup,New_AddrGroup_Name,,,Some_Device_Group\n")
			obj.Write("Change_Me_Service,rename-service,New_Service_Name,,,vsys1\n")
			obj.Write("Change_Me_ServiceGroup,rename-servicegroup,New_ServiceGroup_Name,,,vsys1\n")
			obj.Write("Object_to_Tag,ip,10.1.1.1,,Servers,Some_Device_Group\n")
			obj.End()

			pol, err := easycsv.NewCSV("Policy.csv")
			if err != nil {
				log.Printf("Error creating CSV file - %s", err)
				os.Exit(1)
			}
			pol.Write("#Name,Type,Description,Tags,SourceZones,SourceAddresses,NegateSource,SourceUsers,HipProfiles,DestinationZones,DestinationAddresses,NegateDestination,Applications,Services,Categories,Action,LogSetting,LogStart,LogEnd,Disabled,Schedule,IcmpUnreachable,DisableServerResponseInspection,Group,Virus,Spyware,Vulnerability,UrlFiltering,FileBlocking,WildFireAnalysis,DataFiltering\n")
			pol.Write(",universal,,,,,FALSE,any,,,,FALSE,any,,any,allow,default,FALSE,TRUE,FALSE,,FALSE,FALSE,,,,,,,,\n")
			pol.End()
		}
	},
}

func init() {
	rootCmd.AddCommand(templateCmd)

	templateCmd.Flags().StringVarP(&t, "type", "t", "all", "Type of the template to generate <objects|policy|all>")
	templateCmd.MarkFlagRequired("type")
}
