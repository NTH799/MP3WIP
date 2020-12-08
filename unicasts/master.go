package unicasts

import (
	"MP3/utils"
	"net"
)

//create the master server that handles all the other processes
func MasterServer(config *utils.Config) {
	ln, err := net.Listen("tcp", ":"+config.Master.Port)
	utils.Error(err)
	//use goroutine so it stays active
	go handleServer(config, ln)
}

//start connecting to all the nodes servers
func ConnectToServer(config *utils.Config) {
	nodes := config.Nodes
	for _, n := range nodes {
		ip := n.Ip
		port := n.Port
		CONN := ip + ":" + port
		//connect to the server of the node
		conn, err := net.Dial("tcp", CONN)
		utils.Error(err)
		//add server to list of connections
		config.Master.Conns = append(config.Master.Conns, conn)
	}
}
