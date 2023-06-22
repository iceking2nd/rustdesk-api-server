package Hash

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func StringToMD5(message string) (string, error) {
	hash := md5.New()
	_, err := hash.Write([]byte(message))
	if err != nil {
		return "", fmt.Errorf("generating MD5 error: %w", err)
	}
	bytes := hash.Sum(nil)
	hashCode := hex.EncodeToString(bytes)
	return hashCode, nil
}
