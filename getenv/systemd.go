package getenv

import (
	"fmt"
	"os"
)

func ConfigDirectory() (string, error) {
	return directoryFromEnv("CONFIGURATION_DIRECTORY")
}

func StateDirectory() (string, error) {
	return directoryFromEnv("STATE_DIRECTORY")
}

func RuntimeDirectory() (string, error) {
	return directoryFromEnv("RUNTIME_DIRECTORY")
}

func CacheDirectory() (string, error) {
	return directoryFromEnv("CACHE_DIRECTORY")
}

func LogsDirectory() (string, error) {
	return directoryFromEnv("LOGS_DIRECTORY")
}

func directoryFromEnv(name string) (string, error) {
	dir := os.Getenv(name)
	if dir != "" {
		if err := validDir(dir); err != nil {
			return "", fmt.Errorf("invalid dir: :w", err)
		}
		return dir, nil
	}
	return "", fmt.Errorf("missing %s", name)
}

func validDir(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("cannot stat dir %s: %w", path, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is a a file not a dir", path)
	}

	return nil
}
