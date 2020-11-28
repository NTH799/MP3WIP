package config

import(
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//the ReadConfig function will read the config file and return it in string format so it can be referenced in the later functions

func ReadConfig() []string{
	config, err := os.Open("config.txt")
	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(config)
	var text []string

	for scanner.Scan(){
		text = append(text, scanner.Text())
	}
	config.Close()
	return text
}

//order takes the info from read config and returns the ordered list of variables it extracts from the file

func Order() (int, int, []string, []string, []string, []string) {
	config := ReadConfig()
	N, f, ID, IP, ports, faulty := Extract(config)
	ShowConfig(N, f, ID, IP, ports)
	return N, f, ID, IP, ports, faulty
}

func Extract(text []string) (int, int, []string, []string, []string, []string) {
	vals := strings.Split(text[0], " ")
	N, err := strconv.Atoi(vals[0])
	if err != nil {
		fmt.Println(err)
	}
	f, err := strconv.Atoi(vals[1])
	if err != nil {
		fmt.Println(err)
	}
	var ID, IP, ports, faulty []string

	for i := 1; i < len(text); i++ {
		temp := strings.Split(text[i], " ")
		ID = append(ID, temp[0])
		IP = append(IP, temp[1])
		ports = append(ports, temp[2])
		faulty = append(faulty, temp[3])
	}

	return N, f, ID, IP, ports, faulty
}

//ShowConfig displays the config so you can see the data and was helpful in debugging this code
func ShowConfig(N, f int, ID, IP, ports []string) {
	fmt.Println("Read the config successfully. Here is the result:")
	fmt.Printf("N: %v\n", N)
	fmt.Printf("f: %v\n", f)
	fmt.Printf("id: %v\n", ID)
	fmt.Printf("ip: %v\n", IP)
	fmt.Printf("ports: %v\n", ports)
}

//the Faulty function is looking to see whether or not the node is faulty but automatically assumes the claim is false
func Faulty(port string, ports []string, faulty []string) bool {
	isFaulty := false

	for j := 0; j < len(ports); j++ {
		if ports[j] == port {
			if faulty[j] == "f" {
				isFaulty = true
			}
		}
	}
	return isFaulty
}

//the GetID function retrieves the ID from the process given the port is connected to
func GetID(port string, ports []string, IDs []string) string {
	var ID string

	for j := 0; j < len(ports); j++ {
		if ports[j] == port {
			ID = IDs[j]
		}
	}

	return ID
}
// ServerConfig reads the config file and returns information for server node to start
func ServerConfig() (int, int, string) {
	file := ReadConfig()
	N, f, _, _, ports, _ := Extract(file)
	serverPort := ports[0]
	return N, f, serverPort
}
