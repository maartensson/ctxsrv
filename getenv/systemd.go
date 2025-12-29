package getenv

import (
	"fmt"
	"os"
)

func ConfigDirectory() (string, error) {
	dir := os.Getenv("CONFIGURATION_DIRECTORY")
	if dir != "" {
		return dir, nil
	}
	return "", fmt.Errorf("could not find CONFIGURATION_DIRECTORY")
}

func StateDirectory() (string, error) {
	dir := os.Getenv("STATE_DIRECTORY")
	if dir != "" {
		return dir, nil
	}
	return "", fmt.Errorf("could not find STATE_DIRECTORY")
}

func RuntimeDirectory() (string, error) {
	dir := os.Getenv("RUNTIME_DIRECTORY")
	if dir != "" {
		return dir, nil
	}
	return "", fmt.Errorf("could not find RUNTIME_DIRECTORY")
}

func CacheDirectory() (string, error) {
	dir := os.Getenv("CACHE_DIRECTORY")
	if dir != "" {
		return dir, nil
	}
	return "", fmt.Errorf("could not find CACHE_DIRECTORY")
}

func LogsDirectory() (string, error) {
	dir := os.Getenv("LOGS_DIRECTORY")
	if dir != "" {
		return dir, nil
	}
	return "", fmt.Errorf("could not find LOGS_DIRECTORY")
}
