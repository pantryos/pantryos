// Package auth provides JWT (JSON Web Token) authentication functionality for the application.
// This package handles token generation, validation, and extraction from HTTP requests.
// It uses HMAC-SHA256 signing method and supports configurable expiration times.
package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret holds the secret key used for signing and validating JWT tokens.
// It is initialized from the JWT_SECRET environment variable, with a fallback
// to a default secret for development/testing purposes.
//
// Security note: In production, always set the JWT_SECRET environment variable
// to a strong, randomly generated secret key.
var jwtSecret = func() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Use a default secret for testing/development only
		// This should never be used in production
		return []byte("test-secret-key-for-development-only")
	}
	return []byte(secret)
}()

// JWTClaims represents the custom claims structure for JWT tokens.
// It extends the standard JWT registered claims with application-specific fields.
type JWTClaims struct {
	// UserID is the unique identifier of the authenticated user
	UserID int `json:"user_id"`
	// RegisteredClaims contains standard JWT claims like expiration, issued at, etc.
	jwt.RegisteredClaims
}

// GenerateJWT creates a new JWT token for the specified user ID.
// The token includes standard claims (expiration, issued at) and custom claims (user ID).
// Tokens are signed using HMAC-SHA256 and have a default expiration of 24 hours.
//
// Parameters:
//   - userID: The unique identifier of the user to create a token for
//
// Returns:
//   - string: The signed JWT token
//   - error: Any error that occurred during token generation
//
// Security notes:
//   - Tokens expire after 24 hours for security
//   - Uses HMAC-SHA256 signing method
//   - Includes standard JWT claims for validation
func GenerateJWT(userID int) (string, error) {
	// Set token expiration to 24 hours from now
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create claims with user ID and standard JWT claims
	claims := &JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create and sign the token using HMAC-SHA256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken validates a JWT token string and extracts the claims.
// This function verifies the token signature, checks expiration, and validates
// the signing method to ensure the token is authentic and not expired.
//
// Parameters:
//   - tokenString: The JWT token string to validate
//
// Returns:
//   - *JWTClaims: The extracted claims if the token is valid
//   - error: Any error that occurred during validation (invalid signature, expired, etc.)
//
// Security notes:
//   - Validates token signature to prevent tampering
//   - Checks token expiration automatically
//   - Verifies signing method to prevent algorithm confusion attacks
//   - Returns detailed error messages for debugging
func ValidateToken(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}

	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify that the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// Additional validation check
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ExtractToken extracts the JWT token from the Authorization header of an HTTP request.
// It expects the token to be in the "Bearer <token>" format and validates the header structure.
//
// Parameters:
//   - r: The HTTP request containing the Authorization header
//
// Returns:
//   - string: The extracted JWT token
//   - error: Any error that occurred during extraction (missing header, invalid format, etc.)
//
// Security notes:
//   - Validates the Authorization header format
//   - Case-insensitive "Bearer" prefix validation
//   - Returns clear error messages for debugging
//   - Handles missing or malformed headers gracefully
func ExtractToken(r *http.Request) (string, error) {
	// Get the Authorization header
	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		return "", errors.New("authorization header required")
	}

	// Split the header into "Bearer" and the token
	parts := strings.Split(bearerToken, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("invalid authorization header format")
	}

	return parts[1], nil
}
