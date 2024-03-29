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

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete <host1>...<host n>",
	Aliases: []string{"d"},
	Short: "Delete hosts(s) from list",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// hostsFile, err := cmd.Flags().GetString("host-file")
		// if err != nil {
		// 	return err
		// }
		hostsFile := viper.GetString("host-file")
		return deleteAction(os.Stdout, hostsFile, args)
	},
}

//Define delete function that uses Remove() method to remove hosts 

func deleteAction(out io.Writer, hostsFile string, args []string) error {
	hl := &scan.HostsList{}

	if err := hl.Load(hostsFile); err != nil {
		return err
	}

	for _, h := range args {
		if err := hl.Remove(h); err != nil {
			return err
		}
		fmt.Fprintln(out, "Deleted host:", h)
	}
	return hl.Save(hostsFile)
}

func init() {
	hostsCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
