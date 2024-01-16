/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"pragprog.com/rggo/cobra/pScan/scan"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Run a port scan on the hosts",
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile, err := cmd.Flags().GetString("host-file")
		if err != nil {
			return err
		}

		ports, err := cmd.Flags().GetIntSlice("ports")
		if err != nil {
			return err
		}
		return scanAction(os.Stdout, hostsFile, ports)
	},
}



//scanAction function takes an io.Writer interface represnting where to print the output to 
//hostsFile contains the name of the file to load the hosts list from 
//and slice of integer ports representing the ports to scan 


func scanAction(out io.Writer, hostsFile string, ports []int) error {
	//create instance of HostsList type from scan package
	hl := &scan.HostsList{}

	//load content of the hostsFile into the hosts list instance 
	if err := hl.Load(hostsFile); err != nil {
		return err
	}

	results := scan.Run(hl, ports)
	return printResults(out, results)
}

//PrintResults prints the results out, takes in io.Writer and slice of scan.Results as input and returns an error 

func printResults(out io.Writer, results []scan.Results) error {
	//compose the ouput message
	message := ""

	//loop through all the results in the result slice, add the host name and 
	//the list of ports with each status to the message variable.

	for _, r := range results {
		message += fmt.Sprintf("%s:", r.Host)

		if r.NotFound {
			message += fmt.Sprintf(" Host not found\n\n")
			continue
		}

		message += fmt.Sprintln()

		for _, p := range r.PortStates {
			message += fmt.Sprintf("\t%d: %s\n", p.Port, p.Open)
		}

		message += fmt.Sprintln()
	}

	//print the contents of message to io.Writer and return the error
	_, err := fmt.Fprint(out, message)
	return err
}


//add a local flag --ports to allow user specify a slice of ports to be scanned
func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().IntSliceP("ports", "p", []int{22, 80, 443}, "ports to scan")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
