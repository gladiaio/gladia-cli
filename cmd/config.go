package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	configFilename  = ".gladia"
	envGladiaAPIKey = "GLADIA_API_KEY"
)

func GetGladiaConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, configFilename), nil
}

func SaveGladiaKeyToFile(gladiaKey string) error {
	configPath, err := GetGladiaConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	if _, err := writer.WriteString(strings.TrimSpace(gladiaKey) + "\n"); err != nil {
		return err
	}
	if err := writer.Flush(); err != nil {
		return err
	}
	if err := os.Chmod(configPath, 0o600); err != nil {
		return err
	}

	fmt.Printf("Gladia API key saved to %s\n", configPath)
	return nil
}

func GetGladiaKeyFromFile() (string, error) {
	configPath, err := GetGladiaConfigFilePath()
	if err != nil {
		return "", err
	}

	file, err := os.Open(configPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	gladiaKey, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(gladiaKey), nil
}

// ResolveAPIKey returns a key from GLADIA_API_KEY, then ~/.gladia, then flagKey.
func ResolveAPIKey(flagKey string) (string, error) {
	if k := strings.TrimSpace(os.Getenv(envGladiaAPIKey)); k != "" {
		return k, nil
	}
	if k, err := GetGladiaKeyFromFile(); err == nil && k != "" {
		return k, nil
	}
	if k := strings.TrimSpace(flagKey); k != "" {
		return k, nil
	}
	return "", errors.New(missingAPIKeyMessage())
}

func missingAPIKeyMessage() string {
	configPath, _ := GetGladiaConfigFilePath()
	return fmt.Sprintf(`no Gladia API key found.

  • export GLADIA_API_KEY=<your-key>
  • gladia auth set <your-key>  (writes %s)
  • gladia transcribe <source> --gladia-key <your-key>

Get a key at https://app.gladia.io/account`, configPath)
}
