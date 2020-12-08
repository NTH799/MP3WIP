package unicasts

import (
	"MP3/utils"
	"encoding/gob"
	"math"
	"net"
)

//creates a channel that the master server uses to communicate
func handleServer(config *utils.Config, ln net.Listener) {
	c := make(chan utils.Message)
	go handleConnections(c, ln)
	handleValues(c, config)
}

//use goroutine to handle all connections like in other MPs
func handleConnections(c chan utils.Message, ln net.Listener) {
	for {
		conn, err := ln.Accept()
		utils.Error(err)
		go unicast_receive(c, conn)
	}
}

//send the message over the channel
func unicast_receive(c chan utils.Message, conn net.Conn) {
	for {
		decoder := gob.NewDecoder(conn)
		var pmessage utils.Message
		err := decoder.Decode(&pmessage)
		utils.Error(err)
		c <- pmessage
	}
}

//handle values coming in through the channel
func handleValues(ch chan utils.Message, config *utils.Config) {

	n := len(config.Nodes)
	nfaulty := 0
	nodemap := make(map[string]float64)

	for {
		//for every message in our channel
		message := <-ch

		//node has failed, add to number of faulty nodes and send the failure state to all other nodes
		if message.Fail {
			nfaulty++
			delete(nodemap, message.From)
			continue
		} else {
			//update the value with the new value from the channel
			nodemap[message.From] = message.Value
			if len(nodemap) == n-nfaulty && Consensus(nodemap) {
				//send its output to all other nodes its output and break out of loop
				message := utils.Message{Output: true}
				multicast(message, config.Master.Conns)
				break
			}
		}
	}
}

//iterate through all values and make sure list is sorted
//previous method only checked first and last node which led to issues in checking consensus
func Consensus(nodemap map[string]float64) bool {
	for _, val1 := range nodemap {
		for _, val2 := range nodemap {
			if math.Abs(val1-val2) > .001 {
				return false
			}
		}
	}
	return true
}

//update message for all other connections
func multicast(message utils.Message, conns []net.Conn) {
	for _, conn := range conns {
		encoder := gob.NewEncoder(conn)
		err := encoder.Encode(message)
		utils.Error(err)
	}
}
