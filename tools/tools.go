package tools

import (
	"crypto/rand"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/auremsinistram/go-errors"
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

func GenerateCode(charset string, length int) (string, error) {
	if length <= 0 || length > 100 {
		return "", errors.New("length is invalid")
	}

	var result strings.Builder

	result.Grow(length)

	max := big.NewInt(int64(len(charset)))

	for range length {
		random, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", errors.Wrap(err, "tools - GenerateCode - #1")
		}

		result.WriteByte(charset[random.Int64()])
	}

	return result.String(), nil
}
