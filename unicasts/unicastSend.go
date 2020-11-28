package unicasts

import (
	"MP3/message"
	"encoding/gob"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

func Dial(x int, ID string, IPs, IDs, ports []string, nodes map[string]net.Conn) {
	//this line wants to dial the network using the combination of the host ip and port which we will take from the config file
	y, err := net.Dial("tcp", IPs[x] + ":" +ports[x])
	if err != nil {
		fmt.Println("unable to establish connection")
		return
	}
	enc := gob.NewEncoder(y)
	enc.Encode(ID)
	nodes[IDs[x]] = y
}
//we need a delay function in the send to create the upper and lower bound referenced in the MP
func delay(minDelay, maxDelay int, wg *sync.WaitGroup){
	delay := rand.Intn(maxDelay-minDelay)+minDelay
	time.Sleep(time.Duration(delay)*time.Second)
	wg.Done()
}

//unicast send will send the message over the network but we need to implement a wait group so it goes for every node in config before it stops
func UniSend(dest net.Conn, message message.Message){
	waitTime := new(sync.WaitGroup)
	go delay(100, 200, waitTime)
	waitTime.Add(1)
	waitTime.Wait()
	enc := gob.NewEncoder(dest)
	enc.Encode(message)
}

//lastly we need an exit function to tell the nodes when to exit
func NodeExit(nodes map[string]net.Conn, i int){
	for _, val := range nodes {
		signal := message.Message{0, i}
		UniSend(val, signal)
	}
}