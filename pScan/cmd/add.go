/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"pragprog.com/rggo/cobra/pScan/scan"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <host1>...<hostn>",
	Short: "Add new host(s) to list",
	Aliases:    []string{"a"},
	Args: cobra.MinimumNArgs(1), 
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error{
		// hostsFile, err := cmd.Flags().GetString("host-file")
		// if err != nil {
		// 	return err
		// }

		//obtain the value from Viper
		hostsFile := viper.GetString("host-file")
		return addAction(os.Stdout, hostsFile, args)
	},
}

//Implement the function addAction to execute the command's action in RunE

func addAction(out io.Writer, hostsFile string, args []string) error {
	//create an empty instance of scan.HostsList
	hl := &scan.HostsList{}

	//Load the hosts from the file into List
	if err := hl.Load(hostsFile); err != nil {
		return err
	}

	//loop through the args(hosts to be added)
	for _, h := range args {
		if err := hl.Add(h); err != nil {
			return err 
		}
		fmt.Fprintln(out, "Added host:", h)
	}
	return hl.Save(hostsFile)
}



func init() {
	hostsCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
