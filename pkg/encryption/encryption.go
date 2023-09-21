package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

func GenerateEncryptionKey() (string, error) {
	key := make([]byte, 16) // 16 bytes for a 32-character hexadecimal string
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	hexKey := hex.EncodeToString(key)
	return hexKey, nil
}

func EncryptData(input io.Reader, key []byte) (io.Reader, error) {

	plainText, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Creating GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Generating random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Decrypt file
	cipherText := gcm.Seal(nonce, nonce, plainText, nil)
	return bytes.NewReader(cipherText), nil

}

func DecryptData(input io.Reader, key []byte) (io.Reader, error) {

	cipherText, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Creating GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Deattached nonce and decrypt
	nonce := cipherText[:gcm.NonceSize()]
	cipherText = cipherText[gcm.NonceSize():]

	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(plainText), nil
}

func ScrambleString(originalString string, key []byte) (string, error) {

	iv := make([]byte, aes.BlockSize)
	copy(iv, key[:aes.BlockSize])

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	stream := cipher.NewCTR(block, iv)

	encryptedString := make([]byte, len(originalString))
	stream.XORKeyStream(encryptedString, []byte(originalString))

	scrambled := hex.EncodeToString(iv) + hex.EncodeToString(encryptedString)

	return scrambled, nil
}

func UnscrambleString(scrambledString string, key []byte) (string, error) {

	scrambledData, err := hex.DecodeString(scrambledString)
	if err != nil {
		return "", err
	}

	if len(scrambledData) < aes.BlockSize {
		return "", errors.New("scrambled data is too short")
	}

	iv := scrambledData[:aes.BlockSize]
	encryptedString := scrambledData[aes.BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	stream := cipher.NewCTR(block, iv)

	decryptedString := make([]byte, len(encryptedString))
	stream.XORKeyStream(decryptedString, encryptedString)

	return string(decryptedString), nil
}

func ValidateEncryptionKey(key []byte) error {

	// Check if the encryption key has the correct size
	keySize := len(key)
	if keySize != 16 && keySize != 24 && keySize != 32 {
		return fmt.Errorf("invalid encryption key size, encryption key must be 16, 24, or 32 bytes in length")
	}

	return nil
}
