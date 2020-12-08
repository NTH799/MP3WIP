package utils

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Min int
	Max int
	Faulty        int
	Nodes    []Node
	Master  Server
}

type Node struct {
	Id     string
	Input  float64
	Ip     string
	Port   string
	Conns  []net.Conn
	Server net.Listener
}

type Server struct {
	Ip    string
	Port  string
	Conns []net.Conn
}

type Message struct {
	From   string
	Round  int
	Value  float64
	Fail   bool
	Output bool
}


func ReadConfig() Config {
	file, err := os.Open("config.txt")
	Error(err)
	//use a scanner to read the config file
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string
	//append the text so it becomes readable by the program
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	err = file.Close()
	Error(err)

	//get delay and the number of faulty nodes
	line := strings.Split(text[0], " ")
	min, _ := strconv.Atoi(line[0])
	max, _ := strconv.Atoi(line[1])
	faulty, _ := strconv.Atoi(line[2])

	//get master server port and ip
	line = strings.Split(text[1], " ")
	server := Server{Ip: line[0], Port: line[1], Conns: []net.Conn{}}

	//get all the node information while also skipping the information in the file that is for delay and faulty
	var nodes []Node
	for _, line := range text[2:] {
		//add nodes to list
		n := strings.Split(line, " ")
		input, err := strconv.ParseFloat(n[1], 64)
		Error(err)
		node := Node{Id: n[0], Input: input, Ip: n[2], Port:n[3], Conns: []net.Conn{}}
		nodes = append(nodes, node)
	}
	return Config{Min: min, Max: max, Faulty: faulty, Nodes: nodes, Master: server}
}

func Error(err error) {
	if err != nil {
		fmt.Println("Error. This program will now exit with code: ", err.Error())
		os.Exit(1)
	}
}
