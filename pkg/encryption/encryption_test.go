package encryption_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/LordMathis/GitEcho/pkg/encryption"
)

func TestEncryptDecryptData(t *testing.T) {
	// Define a test key (replace with your key generation logic)
	key := []byte("0123456789abcdef0123456789abcdef")

	// Create a test data string
	testData := "This is a test data string."

	// Encrypt the test data
	encryptedData, err := encryption.EncryptData(bytes.NewBufferString(testData), key)
	if err != nil {
		t.Fatalf("Error encrypting data: %v", err)
	}

	// Decrypt the encrypted data
	decryptedData, err := encryption.DecryptData(encryptedData, key)
	if err != nil {
		t.Fatalf("Error decrypting data: %v", err)
	}

	// Convert the decrypted data back to a string
	decryptedText, err := io.ReadAll(decryptedData)
	if err != nil {
		t.Fatalf("Error converting decrypted data: %v", err)
	}

	// Compare the original data and decrypted data
	if testData != string(decryptedText) {
		t.Errorf("Original data and decrypted data do not match.\nOriginal: %s\nDecrypted: %s", testData, decryptedText)
	}
}

func TestStringScrambleUnscramble(t *testing.T) {
	// Define a test key (replace with your key generation logic)
	key := []byte("0123456789abcdef0123456789abcdef")

	// Create a test string
	testString := "This is a test string."

	// Scramble the test string
	scrambledString, err := encryption.ScrambleString(testString, key)
	if err != nil {
		t.Fatalf("Error scrambling string: %v", err)
	}

	// Unscramble the scrambled string
	unscrambledString, err := encryption.UnscrambleString(scrambledString, key)
	if err != nil {
		t.Fatalf("Error unscrambling string: %v", err)
	}

	// Compare the original string and unscrambled string
	if testString != unscrambledString {
		t.Errorf("Original string and unscrambled string do not match.\nOriginal: %s\nUnscrambled: %s", testString, unscrambledString)
	}

	// Scramble the test string again
	scrambledString2, err := encryption.ScrambleString(testString, key)
	if err != nil {
		t.Fatalf("Error scrambling string: %v", err)
	}

	// Ensure that multiple calls to scrambling lead to the same scrambled string
	if scrambledString != scrambledString2 {
		t.Errorf("Multiple scramblings with the same input do not produce the same result.\nScrambled1: %s\nScrambled2: %s", scrambledString, scrambledString2)
	}
}
