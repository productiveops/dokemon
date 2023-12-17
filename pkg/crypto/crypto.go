package crypto

import (
	"crypto/sha256"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

func HashBytes(m []byte) string {
	h := sha256.New()
	h.Write((m))
	bs := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(bs)
}

func HashString(m string) string {
	return base64.StdEncoding.EncodeToString([]byte(HashBytes([]byte(m))))
}

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}