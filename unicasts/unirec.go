package unicasts

import (
	"PM3/utility"
	"encoding/gob"
	"fmt"
	"net"
	//"os"
)

func unicast_receive(msg utility.Message, pmessage *[]utility.Message, states *[]float64) (newmsg *utility.Message, newqueue []float64) {
	var updatemsg utility.Message
	if msg.State > 0 {
		*states = append(*states, msg.State)
		*pmessage = append(*pmessage, msg)
		sum := 0.00
		l := len(*pmessage)
		if l == 2 {
			for i := 0; i < l; i++ {
				sum += ((*states)[i])
			}
			updateState := (float64(sum)) / (float64(l))
			R := msg.R + 1
			updatemsg.State = updateState
			updatemsg.R = R
			return &updatemsg, *states
		}
	}
	return &updatemsg, *states
}

func handleConnection(c net.Conn, states *[]float64, pmessage *[]utility.Message) {
	var finalqueue []float64

	for {
		decoder := gob.NewDecoder(c) //use gob to send the message
		msg := new(utility.Message)

		_ = decoder.Decode(msg)
		updatemsg, nodemap := unicast_receive(*msg, pmessage, states)

		encoder := gob.NewEncoder(c)
		encoder.Encode(*updatemsg)

		if len(nodemap) == 3{
			// this line is a WIP as the number on the right of the equation is equal to the amount of entries in the config file
			finalqueue = nodemap
			//I am looking for a better place to call this as since this is called in a goroutine it does not output correctly
			utility.StartSort(finalqueue)

		}
	}

	c.Close()
}

func StartServer(NodeNum string) {
	var states []float64
	p := &states
	var message []utility.Message
	pmessage := &message
	_, port, _ := utility.GetHostPort(NodeNum)

	//receives port from config and adds colon to make it work with net.Listen
	PORT := ":" + port
	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		//use handleconnection as goroutine to begin feeding messages through channel
		go handleConnection(c, p, pmessage)
	}
}