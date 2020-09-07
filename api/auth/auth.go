package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword creates a hash string for a given input string using Bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// ValidatePassword compares a hash with a plain password to see if they are equal
func ValidatePassword(hash, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}

	return true
}
