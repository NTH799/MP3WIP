package unicasts

import (
	"encoding/gob"
	"fmt"
	"net"
	"strconv"
	"sync"
	"PM3/utility"
)

//Send the message with the delay bounds from the config file
func Unicast_send(initmsg utility.Message, n int) {
	for {
		for NumNodes := 1; NumNodes <= n; NumNodes++ {
			StrNum := strconv.Itoa(NumNodes)
			ip, port, _ := utility.GetHostPort(StrNum)
			min, max := utility.GetDelay()
			unicast_send(StrNum, ip+":"+port, initmsg, min, max)
		}
	}

}

//Sends message to the destination process
func unicast_send(process string, destination string, initmsg utility.Message, min_delay int, max_delay int) {
	//dial to the TCP server using net library
	conn, err := net.Dial("tcp", destination)
	if err != nil {
		fmt.Println(err)
		return
	}
	encoder := gob.NewEncoder(conn)
	encoder.Encode(initmsg)
	decoder := gob.NewDecoder(conn) //begin gob process to send message through connection
	var updatemsg utility.Message
	_ = decoder.Decode(&updatemsg)
	groupTest := new(sync.WaitGroup)
	go utility.Delay(min_delay, max_delay, groupTest)
	//use the waitgroup to apply the delay
	groupTest.Add(1)
	groupTest.Wait()


}