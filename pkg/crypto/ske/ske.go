package ske // Symmetric Key Encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
)

type Ske struct {
	key []byte
}

var (
	ske Ske
)

func Init(key string) {
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if(err != nil) {
		panic(err)
	}

	if(len(keyBytes) != 32) {
		panic("key should be 32 byte long")
	}

	ske = Ske{key: keyBytes}
}

func Encrypt(plaintext string) (string, error) {
    aes, err := aes.NewCipher(ske.key)
    if err != nil {
        return "", err
    }

    gcm, err := cipher.NewGCM(aes)
    if err != nil {
        return "", err
    }

    // We need a 12-byte nonce for GCM (modifiable if you use cipher.NewGCMWithNonceSize())
    // A nonce should always be randomly generated for every encryption.
    nonce := make([]byte, gcm.NonceSize())
    _, err = rand.Read(nonce)
    if err != nil {
        return "", err
    }

    // ciphertext here is actually nonce+ciphertext
    // So that when we decrypt, just knowing the nonce size
    // is enough to separate it from the ciphertext.
    ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(ciphertext string) (string, error) {
    aes, err := aes.NewCipher(ske.key)
    if err != nil {
        return "", err
    }

    gcm, err := cipher.NewGCM(aes)
    if err != nil {
        return "", err
    }

	ct, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
        return "", err
    }
    // Since we know the ciphertext is actually nonce+ciphertext
    // And len(nonce) == NonceSize(). We can separate the two.
    nonceSize := gcm.NonceSize()
    nonce, ct := ct[:nonceSize], ct[nonceSize:]

    plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ct), nil)
    if err != nil {
        return "", err
    }

    return string(plaintext), nil
}

func GenerateRandomKey() (string, error) {
    b := make([]byte, 32)
    _, err := rand.Read(b)
    if err != nil {
        return "", err
    }
    return base64.StdEncoding.EncodeToString(b), nil
}