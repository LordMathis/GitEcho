package database_test

import (
	"testing"

	"github.com/LordMathis/GitEcho/pkg/encryption"
)

func setEncryptionKey(t *testing.T) {
	encryption.SetEncryptionKey([]byte("12345678901234567890123456789012"))
}
