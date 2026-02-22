package config

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// LoadAPIKey returns BRAVE_API_KEY from the first source that has it:
//  1. Environment variable (already exported into the process environment)
//  2. ~/.config/brave-cli/.env  (tool-specific dotenv, KEY=value)
//  3. ~/.secrets                (export KEY=value, shared secrets file)
//  4. ~/.zshenv                 (zsh: loaded for all sessions, best place for exported vars)
//  5. ~/.zshrc                  (zsh: interactive sessions)
//  6. ~/.bashrc                 (bash: interactive non-login shells)
//  7. ~/.bash_profile           (bash: login shells)
//  8. ~/.profile                (POSIX sh login shells)
//  9. ~/.env                    (generic dotenv fallback)
func LoadAPIKey() (string, error) {
	const key = "BRAVE_API_KEY"

	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v, nil
	}

	home, _ := os.UserHomeDir()

	paths := []string{
		filepath.Join(home, ".config", "brave-cli", ".env"),
		filepath.Join(home, ".secrets"),
		filepath.Join(home, ".zshenv"),
		filepath.Join(home, ".zshrc"),
		filepath.Join(home, ".bashrc"),
		filepath.Join(home, ".bash_profile"),
		filepath.Join(home, ".profile"),
		filepath.Join(home, ".env"),
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
