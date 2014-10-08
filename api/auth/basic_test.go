package auth

import (
	"testing"
)

func Test_basicCredentials(t *testing.T) {
	var tests = []struct {
		input    string
		username string
		password string
		err      *AuthError
	}{
		// Empty input
		{"", "", "", errNoAuthorizationHeader},
		// No Authorization type
		{"abcdef012346789", "", "", errNoAuthorizationType},
		// Not HTTP Basic
		{"Hello abcdef012346789", "", "", errNotBasicAuthorization},
		// Invalid base64
		{"Basic !@#", "", "", errInvalidBase64Authorization},
		// Invalid credential pair (empty password)
		{"Basic dGVzdDo=", "", "", errInvalidBasicCredentialPair},
		// Valid pair
		{"Basic dGVzdDp0ZXN0", "test", "test", nil},
	}

	for _, test := range tests {
		// Split input header into credentials
		username, password, err := basicCredentials(test.input)
		if err != nil {
			// Check for expected error
			if err == test.err {
				continue
			}
		}

		// Verify username
		if username != test.username {
			t.Fatalf("unexpected username: %v != %v", username, test.username)
		}

		// Verify password
		if password != test.password {
			t.Fatalf("unexpected password: %v != %v", password, test.password)
		}
	}
}
