package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/rand"
)

func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
}

func DecodeToJSON(src io.Reader, dst any) error {
	data, err := io.ReadAll(src)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, dst); err != nil {
		return err
	}
	return nil
}

func PrintInJSONFormat(data interface{}) {
	dataJson, _ := json.MarshalIndent(data, "", "\t")
	fmt.Println(string(dataJson))
}

func GenerateAlphaNumericString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func EncryptString(target string) (string, error) {
	encryptStr, err := bcrypt.GenerateFromPassword([]byte(target), bcrypt.DefaultCost)
	return string(encryptStr), err
}
