package main

import (
	"encoding/json"
	"os"
)

type Profile struct {
	BeaconURL string `json:"beacon_url"`
	UserAgent string `json:"user_agent"`
	Interval  int    `json:"interval"`
}

var activeProfile = Profile{
	BeaconURL: "/beacon",
	UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
	Interval:  10,
}

func init() {
	data, err := os.ReadFile("profile.json")
	if err != nil {
		return
	}
	if err := json.Unmarshal(data, &activeProfile); err != nil {
		return
	}
}
