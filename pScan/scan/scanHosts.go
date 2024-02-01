package scan

import (
	"fmt"
	"net"
	"time"
)

//PortState represents the state of a single TCP port


type PortState struct {
	Port int 
	Open state
}

type state bool 

//String converts the boolean value of state to a human readable string
func(s state) String() string {
	if s {
		return "open"
	}
	return "closed"
}

//scanPort performs a port scan on a single TCP port 
func scanPort(host string, port int) PortState {
	p := PortState{
		Port: port,
	}
	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	scanConn, err := net.DialTimeout("tcp", address, 1*time.Second)

	if err != nil {
		return p
	}
	scanConn.Close()
	p.Open = true
	return p
}

// Results represent the scan results for a single host 
// The Run() function returns a slice of Results, one for each host in the list 

type Results struct {
	Host 		string
	NotFound 	bool
	PortStates  []PortState
}
//NotFound indicates whether the host can be resolved to a valid Ip Address in the network
//PortStates is a slice of the type PortState indicting the status ofr each port scanned 
//basically we are scanning different ports for one host 


//Run function performs a port scan on the hosts list.
//The function takes in a pointer to the HostList type 
//and a slice of integers representing the ports to scan. It returns a slice of Results 

func Run(hl *HostsList, ports []int) []Results {
	res := make([]Results, 0, len(hl.Hosts))

	//loop through the list of hosts and define an instance of Results for each host 
	for _, h := range hl.Hosts {
		r := Results{
			Host: h,
		}
		//use the net.LookupHost to resolve the host name into a valid Ip address
		if _, err := net.LookupHost(h); err != nil {
			r.NotFound = true
			res = append(res, r)
			continue
		}

		for _, p := range ports {
			r.PortStates = append(r.PortStates, scanPort(h, p))
		}
		res = append(res, r)

	}
	return res
	
}
