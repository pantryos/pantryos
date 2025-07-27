// Package utils provides utility functions for common operations across the application.
// This package includes password hashing utilities that use bcrypt for secure password storage.
package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword generates a bcrypt hash of the provided password.
//
// This function uses bcrypt with a cost factor of 14, which provides a good balance
// between security and performance. The cost factor determines the number of rounds
// used in the hashing process - higher values are more secure but slower.
//
// Parameters:
//   - password: The plain text password to hash
//
// Returns:
//   - string: The bcrypt hash of the password (60 characters)
//   - error: Any error that occurred during hashing
//
// Security notes:
//   - bcrypt automatically generates a random salt for each password
//   - The same password will produce different hashes due to the salt
//   - bcrypt has a 72-byte limit on password length
//   - Cost factor 14 provides approximately 2^14 = 16,384 rounds
func HashPassword(password string) (string, error) {
	// Use cost factor 14 for a good balance of security and performance
	// This provides approximately 2^14 = 16,384 rounds of hashing
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares a plain text password with a bcrypt hash to verify
// if the password matches the hash.
//
// This function safely compares the password against the stored hash without
// revealing timing information that could be used in timing attacks.
//
// Parameters:
//   - password: The plain text password to verify
//   - hash: The bcrypt hash to compare against
//
// Returns:
//   - bool: true if the password matches the hash, false otherwise
//
// Security notes:
//   - Uses constant-time comparison to prevent timing attacks
//   - Returns false for invalid hash formats
//   - Safe to use even if the hash is malformed
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
