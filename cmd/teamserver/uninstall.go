package main

import (
	"fmt"
	"os"
)

func uninstall() {
	err := os.Remove("/usr/local/bin/nullbeacon")
	if err != nil {
		fmt.Printf("Failed to uninstall: %v\n", err)
	} else {
		fmt.Println("Uninstallation successful.")
	}
}
