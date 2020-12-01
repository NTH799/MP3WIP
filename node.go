package main

import (
	"PM3/utility"
	"PM3/unicasts"
	"fmt"
	"net"
	"os"
	"sync"
)
//Starts communicating with every other project
func dial(processPort string, port string, ip string) (status string) {
	if port != processPort {
		address := ip + ":" + port

		_, err := net.Dial("tcp", address)
		if err != nil {
			return "error establishing connection"
		}
	}
	return "successfully established connection"
}

//The function starts the dialing but since we use more than one terminal theres a delay to make sure it works
func initialize(NumNodes string) {
	processIP, host, _ := utility.GetHostPort(NumNodes)
	ports := utility.GetPorts()

	//loop through every port in the config.txt and create a TCP connection between current process' port and others
	for port := range ports {
		//keeps dialing until a successful connection was made
		for {
			status := dial(host, ports[port], processIP)
			if status == "success" {
				break
			}
			fmt.Println("looking for another connection... will try again in a few seconds")

			//create a delay a goroutine and waitgroups
			wg := new(sync.WaitGroup)
			go utility.Delay(4000, 4001, wg)
			wg.Add(1)
			wg.Wait()
		}
	}
}


func main() {
	args := os.Args
	NumNodes := args[1]
	_, _, initialState := utility.GetHostPort(NumNodes)
	initmsg := utility.Message{}

	initmsg.State = initialState
	initmsg.R = 1

	//Begins the server for each node
	go unicasts.StartServer(NumNodes)

	initialize(NumNodes)
	//Sends each of the nodes to the server via unicast send
	unicasts.Unicast_send(initmsg, 3)

}