package main

import (
	"MP3/nodes"
	"MP3/unicasts"
	"MP3/utils"
	"sync"
	"time"
)

func main() {
	start := time.Now()

	//read the config file and return its parameters
	config := utils.ReadConfig()
	//start the master server given the port and ip
	unicasts.MasterServer(&config)
	//start getting the values for the nodes from config
	nodes.StartNodes(&config)
	//start connecting between nodes and master server
	unicasts.ConnectToServer(&config)
	//start connecting between the nodes themselves
	nodes.ConnectNodes(&config)
	//begin the sorting algorithm
	var wg sync.WaitGroup
	nodes.StartSort(&wg, config)
	//wait for all go routines to finish before printing
	wg.Wait()
	duration := time.Since(start)
	finaltime := duration.Seconds()
	println("done in ", finaltime)
}
