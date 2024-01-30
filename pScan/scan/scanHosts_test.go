package scan_test

import (
	"fmt"
	"net"
	"strconv"
	"testing"

	"pragprog.com/rggo/cobra/pScan/scan"
)

// test function to test the String() method of the state type

func TestStateString(t *testing.T) {
	ps := scan.PortState{}

	//because Open is "closed(false)" by default
	if ps.Open.String() != "closed" {
		t.Errorf("Expected %q, got %q instead\n", "closed", ps.Open.String())
	}
	ps.Open = true

	if ps.Open.String() != "open" {
		t.Errorf("Expected %q, got %q instead\n", "open", ps.Open.String())
	}
}

// To test the Run function // This tests for when the host is found (go through and see how this test was generated based on the host being found )
// this test has two cases, Open port and closed port 

func TestRunHostFound(t *testing.T) {
	testCases := []struct {
		name 	 		string
		expectedState	string
	}{
		{"OpenPort", "open"},
		{"ClosedPort", "closed"},
	}

	// create an instance of the scan.Hosts List and add localhost to it 

	host := "localhost"
	hl := &scan.HostsList{}
	hl.Add(host)

	//ports slice 
	ports := []int{}// we should have one open port and one closed port appended here 

	// Initialise ports, 1 open, 1 closed 
	for _, tc := range testCases {

		// using 0 means you want a free available port selected by the system 
		ln, err := net.Listen("tcp", net.JoinHostPort(host, "0"))

		if err != nil {
			t.Fatal(err)
		}

		defer ln.Close()

		//extract port from the Listener address using the Addr() method 

		_, portStr, err := net.SplitHostPort(ln.Addr().String())
		if err != nil {
			t.Fatal(err)
		}
		port, err := strconv.Atoi(portStr)
		if err != nil {
			t.Fatal(err)
		}

		ports = append(ports, port)

		if tc.name == "ClosedPort" {
			ln.Close() 
		}
	}

	fmt.Println(ports)

	// execute the Run() method using the ports slice 
	res := scan.Run(hl, ports)
	fmt.Println(res) // because theres only one host with different ports appended to port slice 

	// Verify results for HostFound test 
	// there should be only one element in the result slice returned by the Run function.

	if len(res) != 1 {
		t.Fatalf("Expected 1 results, got %d instead\n", len(res))
	}

	// the host name in the result should match the variable host
	if res[0].Host != host {
		t.Errorf("Expected host %q got %q instead\n", host, res[0].Host)
	}


	// the property NotFound should be false since we expect this host to exist 
	if res[0].NotFound {
		t.Errorf("Expected host %q to be found\n", host)
	}

	// verify that two ports are present in th eportzstates slice 
	if len(res[0].PortStates) != 2 {
		t.Fatalf("Expected 2 port states, got %d instead\n", len(res[0].PortStates))
	}

	// verify each port state by looping through each test case 
	// and check if the port number and state match the expected values

	for i, tc := range testCases {
		if res[0].PortStates[i].Port != ports[i] {
			t.Errorf("Expected port %d, got %d instead\n", ports[0],
				res[0].PortStates[i].Port)
		}

		if res[0].PortStates[i].Open.String() != tc.expectedState {
			t.Errorf("Expected port %d to be %s\n", ports[i], tc.expectedState)
		}
	}
}

// Test for when the host is not found
func TestRunHostNotFound(t *testing.T) {
	// create an instance of scan.HostsList and add 389.389.389.389
	host := "389.389.389.389"
	hl := &scan.HostsList{}

	hl.Add(host)

	// execute the Run() method using an empty slice as the port argument since the host doesnt exist 

	res := scan.Run(hl, []int{})

	// verify results for the HostNotFound test 

	// there should be only on eelement in the results slice retruned by the Run() function
	if len(res) != 1 {
		t.Fatalf("Expected 1 result, got %d instead\n", len(res))
	}

	// the host name in the result should amtch the variable host name 
	if res[0].Host != host {
		t.Errorf("Expected host %q, got %q instead\n", host, res[0].Host)
	}

	//the property NotFound should be true since we do not expect this host to exist 
	if !res[0].NotFound {
		t.Errorf("Expected host %q NOT to be found\n", host)
	}

	// The PortStates slice should contain no elements as the scan should be skipped for this host 
	//This can be seen when we use net.Lookuphost, continue 

	if len(res[0].PortStates) != 0 {
		t.Fatalf("Expected 0 port states, got %d instead\n", len(res[0].PortStates))
	}
} 
