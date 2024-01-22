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

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List hosts in hosts list",
	RunE: func(cmd *cobra.Command, args []string) error {
		// hostsFile, err := cmd.Flags().GetString("host-file")
		// if err != nil {
		// 	return err
		// }
		hostsFile := viper.GetString("host-file")
		return listAction(os.Stdout, hostsFile, args)
	},
	Aliases: []string{"l"},
}

func listAction(out io.Writer, hostsFile string, args []string) error {
	//create instance of HostsList type from scan package
	hl := &scan.HostsList{}

	//load the content of hostsFile into the hosts list 
	if err := hl.Load(hostsFile); err != nil {
		return err
	}

	//iterate over list and print each item 
	for _, h := range hl.Hosts {
		if _, err := fmt.Fprintln(out, h); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	hostsCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
