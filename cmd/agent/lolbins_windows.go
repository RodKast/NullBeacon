package main

import (
	"encoding/base64"
	"os/exec"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func RunEncoded(command string) (string, error) {
	encoder := unicode.UTF16(unicode.LittleEndian,
		unicode.IgnoreBOM).NewEncoder()
	utf16, _, err := transform.Bytes(encoder, []byte(command))
	if err != nil {
		return "", err
	}
	b64 := base64.StdEncoding.EncodeToString(utf16)
	cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive",
		"-EncodedCommand", b64)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func DownloadFile(url, dest string) (string, error) {
	cmd := exec.Command("certutil.exe", "-urlcache", "-split", "-f", url, dest)
	out, err := cmd.CombinedOutput()
	return string(out), err
}
