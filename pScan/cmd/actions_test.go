package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"pragprog.com/rggo/cobra/pScan/scan"
)

//auxillairy function to set up test environment.
//includes, creating a temp file and initializing a list if required

//function accepts a slice of strings representing hosts to initilise a list, bool to indicate whether the list should be initialized
//returns the name of the temp file as string and cleanup function that deletes the temp file after it has been used
func setup(t *testing.T, hosts []string, initList bool) (string, func()) {
	// Create temp file 
	tf, err := ioutil.TempFile("", "pScan")
	if err != nil {
		t.Fatal(err)
	}
	tf.Close()

	// Initilize list if needed 
	if initList {
		hl := &scan.HostsList{}

		for _ , h := range hosts {
			hl.Add(h)
		}

		if err := hl.Save(tf.Name()); err != nil {
			t.Fatal(err)
		}
	}

	// Return the temp file name and the cleanup function 
	return tf.Name(), func() {
		os.Remove(tf.Name())
	}

}

//Test function to test the action functions

func TestHostActions(t *testing.T) {
	//Define hosts for actions test 

	// these are used to initialise the list in the setup function 
	hosts := []string{
		"host1",
		"host2",
		"host3",
	}

	// Test cases for Action test 
	// args - a list of args to pass to the action function 
	//initList - indicates whether the list must be initialised before the test 
	//actionFunction - represents which action function to test 
	testCases := []struct {
		name 			string
		args			[]string
		expectedOut		string
		initList 		bool 
		actionFunction	func(io.Writer, string, []string) error
	} {
		{
			name: 			"AddAction",
			args: 			hosts,
			expectedOut:	"Added host: host1\nAdded host: host2\nAdded host: host3\n",
			initList: 		false, 
			actionFunction: addAction,
		}, 
		{
			name : 			"ListAction",
			expectedOut: 	"host1\nhost2\nhost3\n",
			initList:		true, 
			actionFunction: listAction,
		},
		{
			name:			"DeleteAction",
			args:			[]string{"host1", "host2"},
			expectedOut:    "Deleted host: host1\nDeleted host: host2\n",
			initList: 		true,
			actionFunction: deleteAction,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T){
			//Setup Action test 
			tf, cleanup := setup(t, hosts, tc.initList)
			defer cleanup()

			// Define var to capture Action output 
			var out bytes.Buffer 

			// Execute Action and capture the output
			if err := tc.actionFunction(&out, tf, tc.args); err != nil {
				t.Fatalf("Expected no error, got %q\n", err)
			}

			// compare the output of the action fucntion with the expected output and fail the test if they dont match 
			// Test Actions outout 
			if out.String() != tc.expectedOut {
				t.Errorf("Expected output %q, got %q\n", tc.expectedOut, out.String())
			}
		})
	}
}


//INTEGRATION TEST - The goal is to execute all commands in sequence, simulating what a user would do with the tool.
// we will simulate where a user
// 1. adds three hosts to the list 
// 2. Prints them out
// 3. Deletes a host from the list
// 4. Prints the list again 


func testIntegration(t *testing.T) {
	// Define hosts for integration test 

	hosts := []string{
		"host1", 
		"host2",
		"host3",
	}

	// Set up the test using the setup function
	tf, cleanup := setup(t, hosts, false)
	// the cleanup function is defered to ensure the file is deleted after the tests 
	defer cleanup()

	// create variable to hold the name of the host that will be deleted with the delete operation
	delHost := "host2"

	// this variable represents the end state of the list of hosts after the delete operation
	hostsEnd := []string{
		"host1",
		"host3",
	}

	// define variable bytes.buffer to capture output for the Integrated test 
	// later gotten from out.string 

	var out bytes.Buffer

	// Define the expected output by concatenating the ouput of all the operations that will be executed during the test 
	// Remember that the goal of the execution test is to execute all the commands in sequence so the EXPECTED OUTPUT will be all the outputs in sequence concatenated together

	// RECALL 
	// we will simulate where a user
	// 1. adds three hosts to the list (add)
	// 2. Prints them out(list)
	// 3. Deletes a host from the list(Delete)
	// 4. Prints the list again (List)

	expectedOut := ""

	// loop through the hosts slice to create the output for the add operation(this is because the add operation adds everything from the host list)
	//(add)
	for _,  v := range hosts{
		expectedOut += fmt.Sprintf("Added host: %s\n", v)
	}
	//(list)
	// join the items of the hosts slice with a newline character \n as the output of the list operation
	expectedOut += strings.Join(hosts, "\n")
	expectedOut += fmt.Sprintln()
	//(delete)
	expectedOut += fmt.Sprintf("Deleted host: %s\n", delHost)
	//(list after delete)
	expectedOut += strings.Join(hostsEnd, "\n")
	expectedOut += fmt.Sprintln()


	// Now execute all the operations in the defined sequence add -> list -> delete -> list 
	// use the buffer variable out to capture the output of all operations. if any of the operations result in an error we fail the test immediately.

	// Add hosts to the list
	if err := addAction(&out, tf, hosts); err != nil {
		t.Fatalf("Expected no error, got %q\n", err)
	}

	// List hosts 
	if err := listAction(&out, tf, nil); err != nil {
		t.Fatalf("Expected no error, got %q\n", err)
	}

	// Delete host2
	if err := deleteAction(&out, tf, []string{delHost}); err != nil {
		t.Fatalf("Expected no error, got %q\n", err)
	}

	// List hosts after delete 
	if err := listAction(&out, tf, nil); err != nil {
		t.Fatalf("Expected no error, got %q\n", err)
	}

	// compare the output of all the operations with the expected output

	if out.String() != expectedOut {
		t.Errorf("Expected output %q got %q\n", expectedOut,out.String())
	}

}



