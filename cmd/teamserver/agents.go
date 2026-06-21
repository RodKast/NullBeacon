package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/RodKast/go-c2/pkg/task"
	"github.com/chzyer/readline"
)

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
