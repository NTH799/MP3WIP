package main

import (
	"MP3/config"
	"MP3/message"
	"MP3/unicasts"
	"MP3/utility"
	"fmt"
	"net"
	"time"
)

func main() {
	nodes, _, port := config.ServerConfig()
	nodemap := make(map[string]net.Conn)    //create a node map of all the connections


	for i := 1; i < nodes +1; i++ {
		unicasts.ServerListen(port, nodemap)
	}

	utility.ReadMap(nodemap)

	start := time.Now() // begin timing how long the sort will take

	var round int

	// keep going until consensus is reached
	for {
		vals := []float64{} //vals will store all the states

		for key, value := range nodemap {
			var nodeRec message.Message
			unicasts.UniReceive(value, &nodeRec)
			if nodeRec.State == 0 {
				delete(nodemap, key)
				fmt.Println("Process", key, "crashed")
			} else {
				vals = append(vals, nodeRec.State)
				round = nodeRec.Round
			}
		}

		utility.ServerRounds(round, vals)

		if utility.Consensus(vals) {
			break
		}
	}

	elapsed := time.Since(start) //stop the timer when consensus is reached and take the time elapsed so it can be printed

	unicasts.NodeExit(nodemap, round)
	fmt.Println("Consensus reached after ", round, " and ", elapsed)
	return
}
