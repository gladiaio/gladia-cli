package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const CONFIG_FILENAME = ".gladia"

func GetGladiaConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, CONFIG_FILENAME), nil
}

func SaveGladiaKeyToFile(gladiaKey string) error {
	configPath, err := GetGladiaConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(gladiaKey + "\n")
	if err != nil {
		return err
	}

	err = writer.Flush()
	if err != nil {
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
