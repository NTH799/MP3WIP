package unicasts

import(
	"MP3/message"
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

//Have the server listen for incoming connections and store the incoming data into the map
func Listen(port string, ID int, IDs []string, nodes map[string]net.Conn){
	for{
		l, err := net.Listen("tcp", ":"+port)
		if err != nil {
			fmt.Println("could not find a client")
			return
		}
	c, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	var temp string
	dec := gob.NewDecoder(c)
	dec.Decode(&temp)
	nodes[temp]=c
	l.Close()
	}
}

func ServerListen(port string, nodes map[string]net.Conn) {
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}

	c, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	var ID string
	dec := gob.NewDecoder(c)
	dec.Decode(&ID)
	nodes[ID] = c
	l.Close()
}

//the unicast receive function uses gob to receive the message from the client and then decode it
func UniReceive(src net.Conn, message *message.Message){
	dec := gob.NewDecoder(src)
	dec.Decode(message)
}

//this function will receive a signal once we reach consensus and then will terminate the connection using the kill command
func Exit(nodes map[string]net.Conn) {
	var exit message.Message
	conn := nodes["0"]
	UniReceive(conn, &exit)
	if exit.State == 0 {
		fmt.Printf("consensus has been reached. current round is %v \n", exit.Round)
		os.Exit(0)
	}
}