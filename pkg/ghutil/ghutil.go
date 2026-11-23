package ghutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

var client = &http.Client{Timeout: 30 * time.Second}

func get(url, accept string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if accept != "" {
		req.Header.Set("Accept", accept)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}
	return resp, nil
}

func RunInstaller(repo string) error {
	url := fmt.Sprintf("https://github.com/%s/releases/latest/download/install.sh", repo)
	resp, err := get(url, "")
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	script, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed reading installer script: %w", err)
	}

	cmd := exec.Command("sh", "-s")
	cmd.Stdin = bytes.NewReader(script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("installer execution failed: %w", err)
	}
	return nil
}

func GetLatestTag(repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	resp, err := get(url, "application/vnd.github+json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}
	if release.TagName == "" {
		return "", fmt.Errorf("no tag_name found in release response")
	}
	return release.TagName, nil
}
