/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// modifyCmd represents the modify command
var policyModifyCmd = &cobra.Command{
	Use:   "modify",
	Short: "A brief description of your command",
	Long: ``
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("modify called")
	},
}

func init() {
	policyCmd.AddCommand(policyModifyCmd)

}
