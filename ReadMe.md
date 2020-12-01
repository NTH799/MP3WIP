This version of the project is a higher level abstraction of my previous as I was having trouble getting the maps to work through tcp

To run this version of the sorting algorithm you need to enter the following lines in terminal
* go run node.go 1
* go run node.go 2
* go run node.go 3

you can run this as many times in various terminals but there must be the same amount of lines in config to run correctly.

this is an obvious bottleneck as it cannot run efficiently or test many nodes over the connection.

an important note is that on line 45 on unirec.go I make the statement:
`if len(nodemap) == 3`

this is a workaround fix as I currently have no way of telling this when it needs to stop receiving values since I am using a goroutine for the server

my last problem is that the `StartSort()` function is in the `handleConnection()` function as of right now.

this is a problem as the program will run infinitely until `command + C` is entered to force quit

this has also led to a lot of issues regarding the round and time keeping, so I am currently looking for input regarding abstraction that could help me accurately fix my calculations

given these challenges I have not yet come up with a reliable way to present analysis but I will continue testing methods until I can pull the sort method out of the go routine to prevent these issues.


Here are some resources I used to create this and I will add more as I continue to work:
https://github.com/hashicorp/raft \
https://gobyexample.com/waitgroups \
https://gobyexample.com/command-line-arguments \
https://medium.com/@felipedutratine/manage-config-in-golang-to-get-variables-from-file-and-env-variables-33d876887152 \
https://tutorialedge.net/golang/go-sorting-with-sort-tutorial \
https://golang.org/pkg/time/
https://gobyexample.com/random-numbers