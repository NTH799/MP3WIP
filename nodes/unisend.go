package nodes

import (
	"MP3/utils"
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

//Handle all the nodes
func handleNodes(wg *sync.WaitGroup, node utils.Node, config utils.Config) {
	defer wg.Done()
	vals := make(chan utils.Message)
	go handleConnections(vals, node.Server)
	handleValues(vals, node, config)
}

//create a goroutine for each connection
func handleConnections(vals chan utils.Message, ln net.Listener) {
	for {
		conn, err := ln.Accept()
		utils.Error(err)
		go unicast_receive(vals, conn)
	}
}

//push messages to channel from connection
func unicast_receive(vals chan utils.Message, conn net.Conn) {
	for {
		decoder := gob.NewDecoder(conn)
		var pmessage utils.Message
		err := decoder.Decode(&pmessage)
		utils.Error(err)
		vals <- pmessage
	}
}

//handle messages from vals channel
func handleValues(vals chan utils.Message, node utils.Node, config utils.Config) {
	id := node.Id
	n := len(config.Nodes)
	faulty := config.Faulty
	nfaulty := 0
	round := 1
	var val float64
	receivedMsg := 0
	sum := 0.
	canCrash := true

	//send initial state to all other nodes
	err := multicast(node.Conns,
		utils.Message{From: id, Round: round, Value: node.Input},
		true,
		config.Min,
		config.Max)
	//handle the crash scenario
	if err != nil {
		fmt.Println("Node "+node.Id+":", err.Error())
		//send to other nodes that the node has crashed
		err = multicast(node.Conns,
			utils.Message{From: id, Fail: true},
			false,
			0, 1)
		utils.Error(err)
		return
	}
	for {
		//go through every val in the channel
		message := <-vals

		//check message flags and current round
		if message.Output {
			println(val)
			break
		} else if message.Fail {
			nfaulty++
			if nfaulty >= faulty {
				canCrash = false
			}
		} else if message.Round < round {
			continue
		} else if message.Round > round {
			vals <- message
		} else {
			//message is on current round and we update count of messages
			receivedMsg++
			sum += message.Value

			//if node has received all other values this round, update it and send its value to all other nodes
			if receivedMsg >= n-nfaulty {
				val = sum / float64(n-nfaulty)
				fmt.Printf("Node %s has finished round %d. With value %f\n", id, round, val)
				err = multicast(node.Conns,
					utils.Message{From: id, Round: round + 1, Value: val},
					canCrash,
					config.Min,
					config.Max)
				//if a node crashes during the multicast send its faulty status to all other nodes
				if err != nil {
					fmt.Println("Node "+node.Id+":", err.Error())
					err = multicast(node.Conns,
						utils.Message{From: id, Fail: true},
						false,
						0, 1)
					utils.Error(err)
					return
				}
				round++
				receivedMsg = 0
				sum = 0.
			}
		}
	}
}

// Send message to every connection in conns, with min and max delay in ms
func multicast(conns []net.Conn, message utils.Message, canCrash bool, min, max int) error {

	//send to master server first
	encoder := gob.NewEncoder(conns[0])
	err := encoder.Encode(message)
	utils.Error(err)

	//send to other connections with chance to crash
	rand.Seed(time.Now().UnixNano())
	for _, conn := range conns[1:] {
		if canCrash {
			//use a random number to generate crash chance which is fixed at 5 percent but chan change
			n := rand.Intn(100)
			crashChance := 2
			if n < crashChance {
				return errors.New("node has crashed")
			}
		}
		// Send as goroutine to avoid bottleneck
		go unicast_send(conn, message, min, max)
	}

	return nil
}

//send the message using a delay defined in config
func unicast_send(conn net.Conn, message utils.Message, min, max int) {
	delay := rand.Intn(max-min) + min
	time.Sleep(time.Duration(delay) * time.Millisecond)
	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(message)
	utils.Error(err)
}
