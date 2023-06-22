package Hash

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
)

func StringToSHA512(message string) (string, error) {
	hash := sha512.New()
	_, err := hash.Write([]byte(message))
	if err != nil {
		return "", fmt.Errorf("generating SHA512 error: %w", err)
	}
	bytes := hash.Sum(nil)
	hashCode := hex.EncodeToString(bytes)
	return hashCode, nil
}
