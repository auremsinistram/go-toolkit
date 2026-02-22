package tools

import (
	"os"
	"strconv"
)

func GetenvInt(key string, defaultValue int) int {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return result
}

func GetenvFloat(key string, defaultValue float64) float64 {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}

	return result
}

func GetenvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	result, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return result
}
