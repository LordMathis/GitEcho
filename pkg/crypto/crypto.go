package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"os"
)

var encryptionKey = []byte(os.Getenv("ENCRYPTION_KEY"))

func Encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}

	// Generate a random initialization vector (IV)
	iv := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return nil, err
	}

	// Encrypt the data
	stream := cipher.NewCFBEncrypter(block, iv)
	encrypted := make([]byte, len(data))
	stream.XORKeyStream(encrypted, data)

	// Append the IV to the encrypted data
	encryptedWithIV := append(iv, encrypted...)

	// Encode the encrypted data as base64 before returning
	return []byte(base64.StdEncoding.EncodeToString(encryptedWithIV)), nil
}

func Decrypt(encryptedData []byte) ([]byte, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}

	// Decode the base64 encoded encrypted data
	encryptedWithIV, err := base64.StdEncoding.DecodeString(string(encryptedData))
	if err != nil {
		return nil, err
	}

	// Extract the IV from the encrypted data
	iv := encryptedWithIV[:aes.BlockSize]
	encrypted := encryptedWithIV[aes.BlockSize:]

	// Decrypt the data
	stream := cipher.NewCFBDecrypter(block, iv)
	decrypted := make([]byte, len(encrypted))
	stream.XORKeyStream(decrypted, encrypted)

	return decrypted, nil
}
