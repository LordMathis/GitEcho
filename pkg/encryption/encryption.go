package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

var encryptionKey []byte

// SetEncryptionKey sets the encryption key manually
func SetEncryptionKey(key []byte) {
	encryptionKey = key
}

func GenerateEncryptionKey() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	fmt.Print(string(key))
	keyBase64 := base64.StdEncoding.EncodeToString(key)
	return keyBase64, nil
}

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

func ValidateEncryptionKey() ([]byte, error) {

	key := os.Getenv("GITECHO_ENCRYPTION_KEY")

	if key == "" {
		return nil, fmt.Errorf(`encryption key not set, please set the GITECHO_ENCRYPTION_KEY environment variable or run with -g flag to genereate an encryption key`)
	}

	decodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	// Check if the encryption key has the correct size
	keySize := len(decodedKey)
	if keySize != 16 && keySize != 24 && keySize != 32 {
		return nil, fmt.Errorf("invalid encryption key size, encryption key must be 16, 24, or 32 bytes in length")
	}

	return decodedKey, nil
}
