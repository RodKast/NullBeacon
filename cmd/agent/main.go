package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/user"
	"strings"
)

var (
	serverAddr = flag.String("addr", "localhost:8080", "Server address")
)

func main() {
	flag.Parse()
	conn, err := net.Dial("tcp", *serverAddr)
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("failed to get hostname: %v", err)
	}
	user, err := user.Current()
	if err != nil {
		log.Fatalf("failed to get current user: %v", err)
	}
	message := fmt.Sprintf("%s:%s\n", user.Username, hostname)
	_, err = conn.Write([]byte(message))
	if err != nil {
		log.Fatalf("failed to send message: %v", err)
	}

	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("failed to read response: %v", err)
	}

	log.Printf("received response: %s", strings.TrimSpace(response))
}
