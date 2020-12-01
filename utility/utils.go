package utility

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

//Gets the ports. My previous function was getting too much data and was hard to index through the outputs
func GetPorts() []string {
	line := 0
	f, err := os.Open("config.txt")
	var ports []string
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		if line != 0 {
			port := strings.Split(scanner.Text(), " ")[2]
			ports = append(ports, port)
		}
		line = line + 1
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return ports
}

//parses config.txt and returns ip and host
func GetHostPort(destination string) (string, string, float64) {
	line := 0
	f, err := os.Open("config.txt")
	if err != nil {

		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		//this scanner will run through config and seperate all the info so it can be returned and used by our unicasts
		if line != 0 {
			process := strings.Split(scanner.Text(), " ")[0]
			ip := strings.Split(scanner.Text(), " ")[1]
			port := strings.Split(scanner.Text(), " ")[2]
			strState := strings.Split(scanner.Text(), " ")[3]
			state, _ := strconv.ParseFloat(strState, 64)
			if process == destination {
				return ip, port, state

			}
		}

		line = line + 1
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return "", "done", 0.745
}

//Simulate network delay by adding an extra layer before sending the message via the TCP channel
func Delay(min int, max int, wg *sync.WaitGroup) {
	num := rand.Intn(max-min) + min
	time.Sleep(time.Duration(num) * time.Millisecond)

	//decrement value of waitgroup and relay the flow of execution back to main
	wg.Done()
}

func StartSort(finalqueue []float64){
	round := 1
 sort.Float64s(finalqueue)
	fmt.Println(finalqueue)
	start := time.Now()
	for Consensus(finalqueue) == false{
		round += round + 1
		for i := 0; i <= len(finalqueue)-1; i++{
			finalqueue[i] = rand.Float64()
		}
	}
	duration := time.Since(start)
	fmt.Println("sorted in time ", duration, "after ", round, "rounds\n")
	fmt.Println("Exiting...")
}

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
func GetDelay() (int, int) {
	f, err := os.Open("./config.txt")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	delays := strings.Fields(scanner.Text())
	min_delay, _ := strconv.Atoi(delays[0])
	max_delay, _ := strconv.Atoi(delays[1])
	f.Close()
	return min_delay, max_delay
}