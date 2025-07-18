package sappress

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// ============ DECRYPTION PART ====================

// hexToBytes converts a hex string to a byte slice.
func hexToBytes(hexStr string) ([]byte, error) {
	return hex.DecodeString(hexStr)
}

// xorBytes XORs two byte slices and returns the result.
func xorBytes(b1, b2 []byte) []byte {
	n := len(b1)
	if len(b2) < n {
		n = len(b2)
	}
	out := make([]byte, n)
	for i := 0; i < n; i++ {
		out[i] = b1[i] ^ b2[i]
	}
	return out
}

// deriveDecryptionKey generates the AES key using content and user secrets.
func deriveDecryptionKey(contentSecretHex, userSecretHex string) ([]byte, error) {
	userKey, err := hexToBytes(userSecretHex)
	if err != nil {
		return nil, fmt.Errorf("invalid user secret hex: %w", err)
	}
	contentKey, err := hexToBytes(contentSecretHex)
	if err != nil {
		return nil, fmt.Errorf("invalid content secret hex: %w", err)
	}

	fixedHex := "b6df5f479367881c42c7800534e7d8dddae6d186"
	fixedBytes, err := hexToBytes(fixedHex)
	if err != nil {
		return nil, fmt.Errorf("invalid fixed bytes: %w", err)
	}

	temp := xorBytes(userKey, fixedBytes)
	temp2 := xorBytes(contentKey, temp)
	finalKey := xorBytes(temp2, fixedBytes)

	return finalKey, nil
}

// pkcs7Unpad removes PKCS7 padding.
func pkcs7Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}
	padding := int(data[len(data)-1])
	if padding == 0 || padding > len(data) {
		return nil, errors.New("invalid padding size")
	}
	return data[:len(data)-padding], nil
}

// DecryptIt decrypts AES-CBC encrypted data using secrets.
func DecryptIt(encryptedData []byte, userHex, contentHex string) ([]byte, error) {
	aesKey, err := deriveDecryptionKey(contentHex, userHex)
	if err != nil {
		return nil, fmt.Errorf("key derivation failed: %w", err)
	}

	// Static IV used for decryption (should be passed optionally for flexibility)
	ivHex := "581bcaf3f2d6281b5a2d2873caff7259"
	iv, err := hexToBytes(ivHex)
	if err != nil {
		return nil, fmt.Errorf("invalid IV: %w", err)
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	if len(encryptedData)%aes.BlockSize != 0 {
		return nil, errors.New("encrypted data length is not a multiple of AES block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(encryptedData))
	mode.CryptBlocks(decrypted, encryptedData)

	// Remove padding
	decryptedBytes, err := pkcs7Unpad(decrypted)
	if err != nil {
		return nil, fmt.Errorf("padding error: %w", err)
	}

	return decryptedBytes, nil
}


// ================= OTHER FUNCTIONS =======================

// CleanFilename sanitizes a string to be a safe ASCII-only filename for Windows.
func CleanFilename(name string) string {
	// Reserved Windows filenames (case-insensitive)
	reserved := map[string]bool{
		"CON": true, "PRN": true, "AUX": true, "NUL": true,
		"COM1": true, "COM2": true, "COM3": true, "COM4": true, "COM5": true, "COM6": true, "COM7": true, "COM8": true, "COM9": true,
		"LPT1": true, "LPT2": true, "LPT3": true, "LPT4": true, "LPT5": true, "LPT6": true, "LPT7": true, "LPT8": true, "LPT9": true,
	}

	// Remove or replace Windows-forbidden characters
	invalidChars := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`)
	name = invalidChars.ReplaceAllString(name, "_")

	// Remove non-ASCII and non-safe characters (allow letters, digits, dash, dot, space, underscore)
	var clean strings.Builder
	for _, r := range name {
		if r > unicode.MaxASCII {
			continue
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '-' || r == '_' || r == ' ' {
			clean.WriteRune(r)
		} else {
			clean.WriteRune('_')
		}
	}
	name = clean.String()

	// Trim leading/trailing periods and spaces
	name = strings.Trim(name, " .")

	// Avoid reserved Windows names
	upper := strings.ToUpper(name)
	if reserved[upper] {
		name = "_" + name
	}

	if name == "" {
		name = "untitled"
	}

	return name
}