package main

import (
	"fmt"

	"github.com/fatih/color"
)

func printHelp() {
	cyan := color.New(color.FgCyan).Add(color.Bold)
	green := color.New(color.FgGreen)
	white := color.New(color.FgWhite)

	fmt.Println()
	cyan.Println("  LISTENERS")
	white.Printf("  %-50s", "listen tcp --lhost <host> --lport <port>")
	green.Println("Start a TCP listener")
	white.Printf("  %-50s", "listeners")
	green.Println("List active listeners")
	white.Printf("  %-50s", "stop <listenerID>")
	green.Println("Stop a listener")

	fmt.Println()
	cyan.Println("  AGENTS")
	white.Printf("  %-50s", "list")
	green.Println("List connected agents")
	white.Printf("  %-50s", "interact <agentID>")
	green.Println("Enter agent shell")
	white.Printf("  %-50s", "generate --os <os> --arch <arch> --lhost <host> --lport <port>")
	green.Println("Generate an agent binary")

	fmt.Println()
	cyan.Println("  AGENT SHELL")
	white.Printf("  %-50s", "<command>")
	green.Println("Queue a task and wait for output")
	white.Printf("  %-50s", "tasks")
	green.Println("List all tasks and their output")
	white.Printf("  %-50s", "back")
	green.Println("Return to main shell")

	fmt.Println()
	cyan.Println("  GENERAL")
	white.Printf("  %-50s", "help")
	green.Println("Show this menu")
	white.Printf("  %-50s", "exit")
	green.Println("Exit NullBeacon")
	fmt.Println()
}
