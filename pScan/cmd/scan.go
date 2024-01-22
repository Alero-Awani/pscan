/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"pragprog.com/rggo/cobra/pScan/scan"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Run a port scan on the hosts",
	RunE: func(cmd *cobra.Command, args []string) error {
		// hostsFile, err := cmd.Flags().GetString("host-file")
		// if err != nil {
		// 	return err
		// }
		hostsFile := viper.GetString("host-file")

		ports, err := cmd.Flags().GetIntSlice("ports")
		if err != nil {
			return err
		}

		rports, err := cmd.Flags().GetString("rports")
		if err != nil {
			return err
		}
		

		// check if the user set the port or rport then decide what to return 

		return scanAction(cmd *cobra.Command,os.Stdout, hostsFile, ports, rports)

	},
}



//scanAction function takes an io.Writer interface represnting where to print the output to 
//hostsFile contains the name of the file to load the hosts list from 
//and slice of integer ports representing the ports to scan 


func scanAction(cmd *cobra.Command, out io.Writer, hostsFile string, ports []int, rports string) error {
	//create instance of HostsList type from scan package
	hl := &scan.HostsList{}

	//load content of the hostsFile into the hosts list instance 
	if err := hl.Load(hostsFile); err != nil {
		return err
	}

	//disable ability to pass both ports and portRange
	//this function wont work because rport and port are never empty because of the defaults 
	// //but it works with empty defaults 
	// if (len(ports) > 3) && (rports != "") {
	// 	flagErr := errors.New("error: Specify either ports or portRange and not both")
	// 	fmt.Println(ports)
	// 	fmt.Println(rports)
	// 	return flagErr
	// }

	// define result variable 
	var results []scan.Results 

	// if cmd.Flags().Changed("rports") {
	// 	fmt.Println("The 'myflag' flag was changed by the user.")
	// } else {
	// 	fmt.Println("The 'myflag' flag was not changed.")
	// }

	if cmd.Flags().Changed("ports") {
		results = scan.Run(hl, ports)
	}
	
	//if portRange is provided loop through it and populate ports
	if cmd.Flags().Changed("rports") {
		portInt := strings.Split(rports, "-")

		start, err := strconv.Atoi(portInt[0])
		if err != nil {
			fmt.Println("Error converting start:", err)
			return err
		}
		end, err := strconv.Atoi(portInt[1])
		if err != nil {
			fmt.Println("Error converting end:", err)
			return err
		}
		if (start >= 1 && end <= 65535) && (end > start) {
			rangeports := []int{}
			for i := start; i <= end; i++ {
				rangeports = append(rangeports, i)
			}
			results = scan.Run(hl, rangeports)
		} else {
			flagErr := errors.New("error: port range should be between 1-65535 | upper port number must be greater than lower port number")
			return flagErr
		}
		
	}


	// results := scan.Run(hl, ports)
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

	scanCmd.Flags().IntSliceP("ports", "p", []int{22,33,44}, "ports to scan")
	
	//multiple port scan input 
	scanCmd.Flags().String("rports", "1-15", "Scan a range of ports")



	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
