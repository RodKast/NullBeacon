package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/RodKast/go-c2/pkg/agent"
	"github.com/RodKast/go-c2/pkg/listener"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
)

var (
	agents     = make(map[string]*agent.Agent)
	agentMu    sync.Mutex
	listeners  = make(map[string]*listener.Listener)
	listenerMu sync.Mutex
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
	cyan.Println()
	cyan.Println("  NullBeacon C2 | For authorized use only")
}

func main() {
	printBanner()
	logFile, err := os.OpenFile("teamserver.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	log.SetOutput(logFile)
	operatorShell()
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
		case command == "listeners":
			listListeners()
		case strings.HasPrefix(command, "listen"):
			startNewListener(command)
		case strings.HasPrefix(command, "stop"):
			stopListener(command)
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
		case strings.HasPrefix(command, "generate"):
			generateAgent(command)
		case command == "help":
			printHelp()
		case strings.HasPrefix(command, "remove"):
			removeAgent(command)
		default:
			fmt.Println("Unknown command.")
		}
	}
}
