package utility

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"sort"
	"time"
)

// StartSort begins the sort by giving the first value and setting the round to 1
func StartSort() (float64, int) {
	round := 1
	rand.Seed(time.Now().UnixNano())
	node := rand.Float64()
	return node, round
}

// GetAverage takes the average of all the nodes in the array
func GetAverage(vals []float64) float64 {
	var sum float64

	for _, val := range vals{
		sum = sum + val
	}
	return sum / float64(len(vals))
}

// Server rounds will show all the values in the array and will repeatedly show in the array so we can see it update
func ServerRounds(r int, values []float64) {
	fmt.Printf("Round:%v\n", r)
	fmt.Println(values)
}

//Consensus will check to see if consensus has been reached and will determine when we can stop the program
func Consensus(vals []float64) bool{
	sort.Float64s(vals)
	if len(vals) == 0{
		return false
	}
	if vals[len(vals)-1]-vals[0] > 0.001{ //this will check for anything other than consensus
		return false
	}
	return true
}

//The ReadMap function will show each key value pair in the map so it can be read in terminal
func ReadMap(nodes map[string]net.Conn){
	for key, val := range nodes {
		fmt.Println("the ID is", key, "with a connection of", val)
	}
}

//The crash function uses a set chance to simulate a node crash
func Crash() {
	//current chance to crash is 10
	num := rand.Intn(100)
	if num <= 10 {
		fmt.Println("node has crashed.")
		os.Exit(0)
	}
}

//RoundCheck will go through the value of the node, round number, and all the states the server has received from the client
func RoundCheck(n float64, round int, states []float64){
	fmt.Printf("round: %v \n", round)
	fmt.Printf("current value %v \n", n)
	fmt.Println("States are as follows:")
	fmt.Println(states)
}

