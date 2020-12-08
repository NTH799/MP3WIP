package nodes

import (
	"MP3/utils"
	"net"
	"sync"
)

//start listening for each node in the config file
func StartNodes(config *utils.Config) {
	for j, node := range config.Nodes {
		ln, err := net.Listen("tcp", ":"+node.Port)
		utils.Error(err)
		config.Nodes[j].Server = ln
	}
}

//connect all nodes to master server and to each other
func ConnectNodes(config *utils.Config) {
	nodes := config.Nodes
	ip := config.Master.Ip
	port := config.Master.Port
	CONN := ip + ":" + port

//connect to master server for all nodes and then connect to each other
	for j := range nodes {
		//dial into master server
		conn, err := net.Dial("tcp", CONN)
		utils.Error(err)
		nodes[j].Conns = append(nodes[j].Conns, conn)
	}
	for i := range nodes {
		//have it connect to all nodes including itself
		for _, serverNode := range nodes {
			ip := serverNode.Ip
			port := serverNode.Port
			CONN := ip + ":" + port
			//connect to all other nodes
			conn, err := net.Dial("tcp", CONN)
			utils.Error(err)
			nodes[i].Conns = append(nodes[i].Conns, conn)
		}
	}
}

func StartSort(wg *sync.WaitGroup, config utils.Config) {
	for _, node := range config.Nodes {
		wg.Add(1)
		go handleNodes(wg, node, config)
	}
}
