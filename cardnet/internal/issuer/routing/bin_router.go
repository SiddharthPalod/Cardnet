package routing

import (
	"cardnet/internal/issuer/model"
	"strings"
	"sync"
)

var (
	mu       sync.RWMutex
	binTable map[string]model.IssuerProfile
)

func Load(table map[string]model.IssuerProfile) {
	mu.Lock()
	defer mu.Unlock()
	binTable = table
}

// NormalizeCardNumber removes all non-digit characters from the card number
// This is useful for extracting BIN from card tokens that may contain dashes, spaces, etc.
func NormalizeCardNumber(cardNumber string) string {
	var digits strings.Builder
	for _, r := range cardNumber {
		if r >= '0' && r <= '9' {
			digits.WriteRune(r)
		}
	}
	return digits.String()
}

// ExtractBIN extracts the first 6 digits from a card number (after normalization)
func ExtractBIN(cardNumber string) string {
	normalized := NormalizeCardNumber(cardNumber)
	if len(normalized) >= 6 {
		return normalized[:6]
	}
	return normalized
}

func ResolveIssuer(cardNumber string) (model.IssuerProfile, bool) {
	// Extract BIN (first 6 digits after normalization)
	bin := ExtractBIN(cardNumber)
	if len(bin) < 6 {
		return model.IssuerProfile{}, false
	}

	mu.RLock()
	defer mu.RUnlock()

	profile, ok := binTable[bin]
	return profile, ok
}
