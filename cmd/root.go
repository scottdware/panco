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
	"log"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

var cfgFile, srcrule, mvwhere, targetrule, pass, txt string
var action, dg, user, device, f, t, query, l, v, source, ostype, config string
var load, multiple, xlate, hit bool

var tag2color = map[string]string{
	"None":           "color0",
	"Red":            "color1",
	"Green":          "color2",
	"Blue":           "color3",
	"Yellow":         "color4",
	"Copper":         "color5",
	"Orange":         "color6",
	"Purple":         "color7",
	"Gray":           "color8",
	"Light Green":    "color9",
	"Cyan":           "color10",
	"Light Gray":     "color11",
	"Blue Gray":      "color12",
	"Lime":           "color13",
	"Black":          "color14",
	"Gold":           "color15",
	"Brown":          "color16",
	"Olive":          "color17",
	"_":              "color18",
	"Maroon":         "color19",
	"Red-Orange":     "color20",
	"Yellow-Orange":  "color21",
	"Forest Green":   "color22",
	"Turquoise Blue": "color23",
	"Azure Blue":     "color24",
	"Cerulean Blue":  "color25",
	"Midnight Blue":  "color26",
	"Medium Blue":    "color27",
	"Cobalt Blue":    "color28",
	"Blue Violet":    "color30",
	"Medium Violet":  "color31",
	"Medium Rose":    "color32",
	"Lavender":       "color33",
	"Orchid":         "color34",
	"Thistle":        "color35",
	"Peach":          "color36",
	"Salmon":         "color37",
	"Magenta":        "color38",
	"Red Violet":     "color39",
	"Mahogany":       "color40",
	"Burnt Sienna":   "color41",
	"Chestnut":       "color42",
}

var color2tag = map[string]string{
	"color0":  "None",
	"color1":  "Red",
	"color2":  "Green",
	"color3":  "Blue",
	"color4":  "Yellow",
	"color5":  "Copper",
	"color6":  "Orange",
	"color7":  "Purple",
	"color8":  "Gray",
	"color9":  "Light Green",
	"color10": "Cyan",
	"color11": "Light Gray",
	"color12": "Blue Gray",
	"color13": "Lime",
	"color14": "Black",
	"color15": "Gold",
	"color16": "Brown",
	"color17": "Olive",
	"color18": "_",
	"color19": "Maroon",
	"color20": "Red-Orange",
	"color21": "Yellow-Orange",
	"color22": "Forest Green",
	"color23": "Turquoise Blue",
	"color24": "Azure Blue",
	"color25": "Cerulean Blue",
	"color26": "Midnight Blue",
	"color27": "Medium Blue",
	"color28": "Cobalt Blue",
	"color30": "Blue Violet",
	"color31": "Medium Violet",
	"color32": "Medium Rose",
	"color33": "Lavender",
	"color34": "Orchid",
	"color35": "Thistle",
	"color36": "Peach",
	"color37": "Salmon",
	"color38": "Magenta",
	"color39": "Red Violet",
	"color40": "Mahogany",
	"color41": "Burnt Sienna",
	"color42": "Chestnut",
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "panco",
	Short: "Command-line tool that interacts with Palo Alto firewalls and Panorama using CSV files",
	Long: `Command-line tool that interacts with Palo Alto firewalls and Panorama using CSV files
	
See https://panco.dev for complete documentation`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.panco.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".panco" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".panco")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func getPassword() string {
	// fmt.Print("Enter Password: ")
	// bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// password := string(bytePassword)
	// fmt.Println()
	fmt.Printf("Enter password: ")
	password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Printf("Error reading password: %s", err)
	}

	return strings.TrimSpace(string(password))
}

func sliceToString(slice []string) string {
	var str string

	for _, item := range slice {
		str += fmt.Sprintf("%s, ", item)
	}

	return strings.TrimRight(str, ", ")
}

func stringToSlice(str string) []string {
	var slice []string

	list := strings.FieldsFunc(str, func(r rune) bool { return strings.ContainsRune(",;", r) })
	for _, item := range list {
		slice = append(slice, strings.TrimSpace(item))
	}

	return slice
}

func userSliceToString(slice []string) string {
	var str string

	for _, item := range slice {
		str += fmt.Sprintf("%s; ", item)
	}

	return strings.TrimRight(str, "; ")
}

func userStringToSlice(str string) []string {
	var slice []string

	list := strings.FieldsFunc(str, func(r rune) bool { return strings.ContainsRune(";", r) })
	for _, item := range list {
		slice = append(slice, item)
	}

	return slice
}

func duplicateObjects(objects map[string]string) (map[string]string, map[string]string) {
	unique := map[string]string{}
	dups := map[string]string{}

	for k, v := range objects {
		_, exists := unique[v]

		if exists {
			dups[k] = v
		} else {
			unique[v] = k
		}
	}

	result := map[string]string{}
	for key := range unique {
		result[unique[key]] = key
	}

	return result, dups
}

func formatDesc(str string) string {
	s := strings.ReplaceAll(str, ",", "")
	s = strings.ReplaceAll(s, "\n", " ")

	return s
}
