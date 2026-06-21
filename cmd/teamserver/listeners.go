package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/RodKast/go-c2/pkg/listener"
)

func startNewListener(command string) {
	parts := strings.Fields(command)
	if len(parts) < 6 {
		fmt.Println("Usage: listen tcp --lhost 0.0.0.0 --lport 8080")
		return
	}

	protocol := parts[1]
	var host string
	var port int

	for i, p := range parts {
		if p == "--lhost" {
			host = parts[i+1]
		}
		if p == "--lport" {
			port, _ = strconv.Atoi(parts[i+1])
		}
	}

	l := listener.NewListener(protocol, host, port)

	ctx, cancel := context.WithCancel(context.Background())
	l.Cancel = cancel

	listenerMu.Lock()
	listeners[l.ID] = l
	listenerMu.Unlock()

	go func() {
		err := l.Start(ctx, handleConnection)
		if err != nil {
			log.Printf("listener %s stopped with error: %v", l.ID, err)
		} else {
			log.Printf("listener %s stopped", l.ID)
		}
	}()

	l.Status = "running"
	log.Printf("started listener %s on %s:%d", l.ID, l.Host, l.Port)
	fmt.Printf("started listener %s on %s:%d\n", l.ID, l.Host, l.Port)

}

func stopListener(command string) {
	parts := strings.Fields(command)
	if len(parts) != 2 {
		fmt.Println("Usage: stop <listener_id>")
		return
	}

	listenerID := parts[1]

	listenerMu.Lock()
	l, exists := listeners[listenerID]
	if !exists {
		listenerMu.Unlock()
		fmt.Printf("listener %s not found\n", listenerID)
		return
	}
	l.Cancel()
	delete(listeners, listenerID)
	listenerMu.Unlock()

	log.Printf("stopped listener %s", listenerID)
	fmt.Printf("stopped listener %s\n", listenerID)
}

func listListeners() {
	listenerMu.Lock()
	defer listenerMu.Unlock()

	if len(listeners) == 0 {
		fmt.Println("No listeners running.")
		return
	}

	fmt.Println("Active listeners:")
	for id, l := range listeners {
		fmt.Printf("ID: %s, Protocol: %s, Host: %s, Port: %d, Status: %s\n", id, l.Protocol, l.Host, l.Port, l.Status)
	}
}
