/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"pragprog.com/rggo/cobra/pScan/scan"
)

var flagErr = errors.New("Put in correct value for ports, should be in format '1-15', '22,33' or single port '22'")

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Run a port scan on the hosts",
	RunE: func(cmd *cobra.Command, args []string) error {

		hostsFile := viper.GetString("host-file")

		ports, err := cmd.Flags().GetString("ports")
		if err != nil {
			return err
		}

		// issue, we want to pass filter to print results without having to pass it through
		filter, err := cmd.Flags().GetString("filter")
		if err != nil {
			return err
		}

		portSlice, err := portAction(ports)
		if err != nil {
			return err
		}
        
		return scanAction(os.Stdout, hostsFile, portSlice, &filter)

	},
}

func portAction(ports string) ([]int, error) {

	// if person passes just one num , it should be converted to int and stored in the list 

	comma_match, _ := regexp.MatchString(`\d,\d`, ports)
	dash_match, _ := regexp.MatchString(`\d-\d`, ports)
	num_match, _ := regexp.MatchString(`\d`, ports)

	rangeports := []int{}

	if dash_match {
		portStr := strings.Split(ports, "-")
		start, err := strconv.Atoi(portStr[0])
		if err != nil {
			fmt.Println("Error converting start:", err)
			return nil ,err
		}
		end, err := strconv.Atoi(portStr[1])
		if err != nil {
			fmt.Println("Error converting end:", err)
			return nil, err
		}
		if (start >= 1 && end <= 65535) && (end > start) {
			for i := start; i <= end; i++ {
				rangeports = append(rangeports, i)
			}

		} else {
			portErr := errors.New("port range should be between 1-65535 | upper port number must be greater than lower port number")
			return nil, portErr 
		}
	} else if comma_match {
		portStr := strings.Split(ports, ",")
		for _, i := range portStr {
			j, err := strconv.Atoi(i)
			if err != nil {
				fmt.Println("Error converting string to integer")
			}
			rangeports = append(rangeports, j)
		}

	} else if num_match {
		num, err := strconv.Atoi(ports)
		if err != nil {
			return nil, flagErr
		}
		rangeports = append(rangeports, num)
	} else {
		return nil, flagErr
	}

	return rangeports, nil
}


func scanAction(out io.Writer, hostsFile string, portSlice []int, filter *string) error {
	hl := &scan.HostsList{}

	if err := hl.Load(hostsFile); err != nil {
		return err
	}

	results := scan.Run(hl, portSlice)
 
	return printResults(out, results, *filter)
}


// we will add it in printresult, hardcoding filter

func printResults(out io.Writer, results []scan.Results, filter string) error {
	//compose the output message
	message := ""

	for _, r := range results {
		message += fmt.Sprintf("%s:", r.Host)

		if r.NotFound {
			message += fmt.Sprintf(" Host not found\n\n")
			continue
		}

		message += fmt.Sprintln()

		for _, p := range r.PortStates {
			if (*&filter == "open" && p.Open) || (*&filter == "closed" && !p.Open) || (filter == "") {
				message += fmt.Sprintf("\t%d: %s\n", p.Port, p.Open)
			}
		}

		message += fmt.Sprintln()

	}

	_, err := fmt.Fprint(out, message)
	return err
}

// add a local flag --ports to allow user specify a slice of ports to be scanned
func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().StringP("ports", "p", "", "Scan ports")
	scanCmd.Flags().String("filter", "", "Show either open or closed ports")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
