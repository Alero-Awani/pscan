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
}

for _, tc := range testCases {
	t.Run(tc.name, func(t *testing.T)){
		hl := &scan.HostsList{}

		//Initialize list //this adding host 1 to the list so that we can check later if 
		//te code works well to check for existing hosts  
		if err := hl.Add("host1")err != nil {
			t.Fatal(err)
		}
		err := hl.Add(tc.host)

		if tc.expectErr !- nil {
			if err == nil {
				t.Fatalf("Expected error, got nil instead\n")
			}

			if ! errors.Is(err, tc.expectErr) {
				t.Errorf("Expected error %q, got %q instead\n",
				tc.expectErr, err)
			}

			return
		}

		if err != nil {
			t.Fatalf("Expected n error, got %q instead\n", err)
		}

		if len(hl.Hosts) != tc.expectLen {
			t.Errorf("Expected list length %d, got %d instead\n",
			tc.expectLen, len(hl.Hosts))
		}

		if hl.Hosts[1] != tc.host {
			t.Errorf("Expected host name %q as index 1, got %q instead\n",
			tc.host, hl.Hosts[1])
		}
	}
}
