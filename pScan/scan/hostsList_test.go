package scan_test

import (
	"testing"

	"pragprog.com/rggo/cobra/pScan/scan"
)

// Define two test cases, one to add a new host and another to add an existing host which should return an error

func TestAdd(t *testing.T) {
	testCases := []struct {
		name		string
		host		string
		expectLen 	int
		expectErr	error
	}{
		{"AddNew", "host2", 2, nil},
		{"AddExisting", "host1", 1, scan.ErrExists},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T)){
			hl := &scan.HostsList{}

			//Initialize list //this adding host 1 to the list so that we can check later if 
			//the code works well to check for existing hosts  
			if err := hl.Add("host1")err != nil {				
				t.Fatal(err)
			}

			// here we actually use the add function for both hosts 
			err := hl.Add(tc.host)

			//TEST FOR HOST1(AddExisting) for expected error 
			// if its not nil then its host1, so for host 1, the expected error should be not be nil, in the next if statement 
			//we check if it is nil , if it is then we throw an error. 
			if tc.expectErr != nil {
				if err == nil {
					t.Fatalf("Expected error, got nil instead\n")
				}

				//error.Is is used to check if an error is of a particular type and returns a bool
				//so here we check if err is not of type tc.expectErr, this means if we get a different typ eof error that is not ErrExists for host1 
				if ! errors.Is(err, tc.expectErr) {
					t.Errorf("Expected error %q, got %q instead\n",
					tc.expectErr, err)
				}

				//Why is there a return here?
				return
			}

			//TEST FOR HOST2(AddNew) for expected err

			//the expected error for host2 is nil , so we check and if it is not nil ,we throw an error 
			if err != nil {
				t.Fatalf("Expected no error, got %q instead\n", err)
			}


			// TEST FOR Addnew AND Addzexisting FOR EXPECTED LENGTH
			//the len of Add new should be 2 if host2 was successfully added and the len of AddExisting should remain 1 if everything works well 
			if len(hl.Hosts) != tc.expectLen {
				t.Errorf("Expected list length %d, got %d instead\n",
				tc.expectLen, len(hl.Hosts))
			}

			//For addnew, we check if the second value is host2, if it isnt we throw an error. 
			// For addexiting , the list will be [host1, host2], because we ran the addnew first which appended host2 to the list (so this is still the same thing for both)

			if hl.Hosts[1] != tc.host {
				t.Errorf("Expected host name %q as index 1, got %q instead\n",
				tc.host, hl.Hosts[1])
			}
		}
	}

}


//EXPLANATION OF THE TEST CASE ABOVE 
// Each entry in the struct represents an entry for the test case 


// In this test case we are using two expected cases to test 
//1. Expected err 
//2. Expected len



// REMOVE() TEST FUNCTION (This has two test cases)
func TestRemove(t *testing.T) {
	testCases := []struct {
		name 		string 
		host 		string 
		expectedLen int 
		expectedErr error 
	}{
		{"RemoveExisting", "host1", 1, nil},
		{"RemoveNotFound", "host3", 1, scan.ErrNotExists},

	}

	for _, tc := rnage testCases {
		t.Run(tc.name, func(t *testing.T) {
			hl := &scan.HostsList{}

			// Initialize list //add these to the hl hosts list 
			for _, h := range []string{"host1", "host2"} {
				if err := hl.Add(h); err != nil {
					t.Fatal(err)
				}
			}

			err := hl.Remove(tc.host)

			//Checking for RemoveNotFound (if the expected error in the struct is not nil but you got nil)
			if tc.expectErr != nil {
				if err == nil {
					t.Fatalf("Expected error, got nil instead\n")
				}

				//check for the kind of error it returned 
				if ! errors.Is(err, tc.expectedErr) {
					t.Errorf("Expected err %q, got %q instead\n",
					tc.expectErr, err)
				}
				return // ends the program and moves on to the next iteration
			}

			// Check for RemoveExisting
			if err != nil {
				t.Fatalf("Expected no error, got %q instead\n" err)
			} 

			
			//Fatal doesnt terminate the program so we move on to the other test cases that depend on it 
			// if th lenght of the list is not equal to the expected len after removing the host 

			if len(hl.Hosts) != tc.expectLen {
				t.Errorf("Expected list length %d, got %d instead\n",
				tc.expectLen, len(hl.Hosts))
			}

			//check if it eas actually removed 
			if hl.Hosts[0] == tc.host {
				t.Errorf("Host name %q shoul dnot be in the list\n", tc.host)
			}
		} )
	}
}


// Create test function to tes Save and Load methods 
// function creates two HostsList instances, initializes the first list and uses the Save()
// method to save it to temp file 
// then it  uses the load method to load the contents of the temporary file into the second lis t
//then it compares both lists 
// The test fails if the contents of the list dont match