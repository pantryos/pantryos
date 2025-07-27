package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	t.Run("Hash Valid Password", func(t *testing.T) {
		password := "mySecurePassword123"
		hash, err := HashPassword(password)

		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEqual(t, password, hash)
		assert.Len(t, hash, 60) // bcrypt hashes are 60 characters
	})

	t.Run("Hash Empty Password", func(t *testing.T) {
		password := ""
		hash, err := HashPassword(password)

		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.Len(t, hash, 60)
	})

	t.Run("Hash Special Characters", func(t *testing.T) {
		password := "!@#$%^&*()_+-=[]{}|;':\",./<>?"
		hash, err := HashPassword(password)

		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEqual(t, password, hash)
	})

	t.Run("Hash Long Password", func(t *testing.T) {
		password := "thisIsAVeryLongPasswordThatIsUnderTheBcryptLimitOf72Bytes"
		hash, err := HashPassword(password)

		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEqual(t, password, hash)
	})

	t.Run("Hash Password Near Bcrypt Limit", func(t *testing.T) {
		// bcrypt has a 72-byte limit, so we test with a password just under that limit
		password := "thisIsAVeryLongPasswordThatIsJustUnderTheBcryptLimitOf72Bytes"
		hash, err := HashPassword(password)

		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEqual(t, password, hash)
	})

	t.Run("Hash Unicode Password", func(t *testing.T) {
		password := "pässwörd123🚀"
		hash, err := HashPassword(password)

		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEqual(t, password, hash)
	})
}

func TestCheckPasswordHash(t *testing.T) {
	t.Run("Valid Password and Hash", func(t *testing.T) {
		password := "mySecurePassword123"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		result := CheckPasswordHash(password, hash)
		assert.True(t, result)
	})

	t.Run("Invalid Password with Valid Hash", func(t *testing.T) {
		password := "mySecurePassword123"
		wrongPassword := "wrongPassword"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		result := CheckPasswordHash(wrongPassword, hash)
		assert.False(t, result)
	})

	t.Run("Valid Password with Invalid Hash", func(t *testing.T) {
		password := "mySecurePassword123"
		invalidHash := "invalidHashString"

		result := CheckPasswordHash(password, invalidHash)
		assert.False(t, result)
	})

	t.Run("Empty Password with Valid Hash", func(t *testing.T) {
		password := ""
		hash, err := HashPassword(password)
		require.NoError(t, err)

		result := CheckPasswordHash(password, hash)
		assert.True(t, result)
	})

	t.Run("Empty Password with Invalid Hash", func(t *testing.T) {
		password := ""
		invalidHash := "invalidHashString"

		result := CheckPasswordHash(password, invalidHash)
		assert.False(t, result)
	})

	t.Run("Unicode Password and Hash", func(t *testing.T) {
		password := "pässwörd123🚀"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		result := CheckPasswordHash(password, hash)
		assert.True(t, result)
	})

	t.Run("Case Sensitive Password Check", func(t *testing.T) {
		password := "MySecurePassword123"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		// Different case should fail
		wrongCasePassword := "mysecurepassword123"
		result := CheckPasswordHash(wrongCasePassword, hash)
		assert.False(t, result)

		// Correct case should pass
		result = CheckPasswordHash(password, hash)
		assert.True(t, result)
	})
}

func TestPasswordHashConsistency(t *testing.T) {
	t.Run("Same Password Produces Different Hashes", func(t *testing.T) {
		password := "mySecurePassword123"

		hash1, err := HashPassword(password)
		require.NoError(t, err)

		hash2, err := HashPassword(password)
		require.NoError(t, err)

		// bcrypt should produce different hashes for the same password due to salt
		assert.NotEqual(t, hash1, hash2)

		// But both should validate correctly
		assert.True(t, CheckPasswordHash(password, hash1))
		assert.True(t, CheckPasswordHash(password, hash2))
	})

	t.Run("Hash Verification Round Trip", func(t *testing.T) {
		testCases := []string{
			"simple",
			"complexPassword123!@#",
			"",
			"pässwörd🚀",
			"veryLongPasswordButUnderBcryptLimit",
		}

		for _, password := range testCases {
			t.Run("Password: "+password, func(t *testing.T) {
				hash, err := HashPassword(password)
				require.NoError(t, err)

				result := CheckPasswordHash(password, hash)
				assert.True(t, result, "Password verification failed for: %s", password)
			})
		}
	})
}

func TestPasswordSecurity(t *testing.T) {
	t.Run("Hash Contains Salt", func(t *testing.T) {
		password := "testPassword"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		// bcrypt hashes should start with $2a$ or $2b$ indicating they contain salt
		assert.True(t, len(hash) >= 60, "Hash should be at least 60 characters")
		assert.Contains(t, hash, "$", "Hash should contain bcrypt format indicators")
	})

	t.Run("Hash Is Not Reversible", func(t *testing.T) {
		password := "originalPassword"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		// The hash should not contain the original password
		assert.NotContains(t, hash, password)

		// The hash should not be the same as the password
		assert.NotEqual(t, hash, password)
	})
}
