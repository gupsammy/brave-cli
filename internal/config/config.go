package config

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// LoadAPIKey returns BRAVE_API_KEY from the first source that has it:
//  1. Environment variable (already exported into shell)
//  2. ~/.config/brave-cli/.env  (KEY=value lines)
//  3. ~/.secrets                (export KEY=value lines, same format as reddit-cli)
func LoadAPIKey() (string, error) {
	const key = "BRAVE_API_KEY"

	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v, nil
	}

	home, _ := os.UserHomeDir()

	paths := []string{
		filepath.Join(home, ".config", "brave-cli", ".env"),
		filepath.Join(home, ".secrets"),
	}
	for _, p := range paths {
		if v := readKeyFromFile(p, key); v != "" {
			return v, nil
		}
	}

	return "", errors.New("missing BRAVE_API_KEY")
}

// readKeyFromFile parses shell-style env files (KEY=value or export KEY=value).
func readKeyFromFile(path, key string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// strip leading "export "
		line = strings.TrimPrefix(line, "export ")
		if strings.HasPrefix(line, key+"=") {
			val := strings.TrimPrefix(line, key+"=")
			val = strings.Trim(val, `"'`)
			return strings.TrimSpace(val)
		}
	}
	return ""
}
