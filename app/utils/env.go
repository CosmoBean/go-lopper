package utils

import "os"

// GetEnvWithDefault returns the value of the environment variable with the given key.
func GetEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
