package shortener

import (
	"crypto/sha256"
	"fmt"
	"github.com/google/uuid"
	"github.com/itchyny/base58-go"
	"math/big"
	"os"
)

// Computes the SHA256 hash for a given input string.
func sha256Of(input string) []byte {
	algorithm := sha256.New()
	algorithm.Write([]byte(input))
	return algorithm.Sum(nil)
}

// Encodes a byte slice into a Base58 string using Bitcoin encoding.
func base58Encoded(bytes []byte) string {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return string(encoded)
}

// Creates a short link based on the initial link and user ID.
func GenerateShortLink(longUrl string) string {
	// Generate a random UUID.
	uniquifier := uuid.New().String()

	// Compute SHA256 of the combined string.
	urlHashBytes := sha256Of(longUrl + uniquifier)

	// Convert the SHA256 hash into a base 10 integer and truncate to 64 bits.
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()

	// Encode the number into a Base58 string and truncate it to 8 characters.
	shortUrl := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))[:8]

	// Return the final short link.
	return shortUrl
}
