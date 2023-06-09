package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

func RunServer(args []string) {
	var port string
	if len(args) == 1 {
		port = "8989"
	}
	if len(args) == 2 {
		if portChecker(args[1]) {
			port = args[1]
		} else {
			log.Println("The port numbers in the range 0-1023 are system ports which are reserved for system services.\nThe allowed ports for your server are between 1024-65535, try entering a number in that range.")
			return
		}
	}
	if len(args) > 2 {
		fmt.Println("[USAGE]: go run . $port")
		return
	}
	listener, err := net.Listen("tcp", "localhost:"+port)
	log.Println("The server is running on localhost:" + port)
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	go Broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}

		Mu.Lock()
		if NumberOfClients > 9 {
			fmt.Fprintf(conn, "[SERVER IS FULL - TRY AGAIN LATER]")
			conn.Close()
		} else {
			NumberOfClients++
			go HandleConn(conn)
		}
		Mu.Unlock()

	}
}

func portChecker(arg string) bool {
	port, err := strconv.Atoi(arg)
	if err != nil {
		return false
	}
	if port < 1024 || port > 65535 {
		return false
	}
	return true
}
