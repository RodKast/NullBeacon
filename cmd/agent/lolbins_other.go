//go:build !windows

package main

func RunEncoded(command string) (string, error)     { return "", nil }
func DownloadFile(url, dest string) (string, error) { return "", nil }
