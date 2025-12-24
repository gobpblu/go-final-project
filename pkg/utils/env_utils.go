package utils

import "os"

func GetEnvOrDefault(key, defaultValue string) (string) {
	value := os.Getenv(key)

	if value == "" {
		value = "7540"
	}
	
	return value
}