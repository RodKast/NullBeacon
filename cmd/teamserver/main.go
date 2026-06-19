package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/RodKast/go-c2/pkg/agent"
)

var (
	agents  = make(map[string]*agent.Agent)
	agentMu sync.Mutex
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to start listener: %v", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("new connection from %s", conn.RemoteAddr())

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("failed to read message: %v", err)
		return
	}

	parts := strings.Split(strings.TrimSpace(message), ":")
	if len(parts) != 2 {
		log.Printf("invalid message format: %s", message)
		return
	}
	username := parts[0]
	hostname := parts[1]

	a := agent.NewAgent(username, hostname, conn.RemoteAddr().String())
	agentMu.Lock()
	agents[a.ID] = a
	agentMu.Unlock()

	log.Printf("agent registered: %s", a.ID)

	ackMesg := strings.ToUpper(strings.TrimSpace(message))
	response := fmt.Sprintf("ACK: %s\n", ackMesg)
	_, err = conn.Write([]byte(response))
	if err != nil {
		log.Printf("failed to write response: %v", err)
		return
	}
}
