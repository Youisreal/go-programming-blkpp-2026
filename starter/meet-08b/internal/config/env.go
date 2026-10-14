package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

// ponytail: KEY=value satu baris, kutip literal, tanpa ekspansi/escape shell.
// Gunakan parser dotenv lengkap jika membutuhkan multiline atau interpolasi.
func readEnvFile(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil, nil // Production dapat memakai environment OS saja.
	}
	if err != nil {
		return nil, fmt.Errorf("baca .env: %w", err)
	}
	defer file.Close()
	values := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for number := 1; scanner.Scan(); number++ {
		line := scanner.Text()
		if number == 1 {
			line = strings.TrimPrefix(line, "\ufeff") // UTF-8 BOM dari editor Windows.
		}
		if !utf8.ValidString(line) || strings.ContainsRune(line, '\x00') {
			return nil, fmt.Errorf(".env baris %d: harus teks UTF-8 tanpa NUL", number)
		}
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		key, value, ok := strings.Cut(line, "=")
		key, value = strings.TrimSpace(key), strings.TrimSpace(value)
		if !ok || !validKey(key) {
			return nil, fmt.Errorf(".env baris %d: gunakan KEY=value", number)
		}
		if len(value) > 0 && (value[0] == '\'' || value[0] == '"') {
			if len(value) < 2 || value[len(value)-1] != value[0] {
				return nil, fmt.Errorf(".env baris %d: kutip tidak ditutup", number)
			}
			value = value[1 : len(value)-1]
		}
		values[key] = value
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("baca .env: %w", err)
	}
	return values, nil
}

func validKey(key string) bool {
	if key == "" {
		return false
	}
	for i, ch := range key {
		if ch != '_' && !(ch >= 'A' && ch <= 'Z') && !(ch >= 'a' && ch <= 'z') && !(i > 0 && ch >= '0' && ch <= '9') {
			return false
		}
	}
	return true
}
