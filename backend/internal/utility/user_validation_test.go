package utility

import (
	"strings"
	"testing"
)

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     bool
	}{
		// Valid passwords
		{
			name:     "valid password with minimum length",
			password: "Pass123w",
			want:     true,
		},
		{
			name:     "valid password with maximum length",
			password: "Pass1234word1234word1234word1234",
			want:     true,
		},
		{
			name:     "valid password with all requirements",
			password: "MyPassword123",
			want:     true,
		},
		// Invalid passwords - length
		{
			name:     "too short password",
			password: "Pass12",
			want:     false,
		},
		{
			name:     "empty password",
			password: "",
			want:     false,
		},
		{
			name:     "password too long",
			password: "Pass1234word1234word1234word12345",
			want:     false,
		},
		// Invalid passwords - missing requirements
		{
			name:     "no uppercase letter",
			password: "password123",
			want:     false,
		},
		{
			name:     "no lowercase letter",
			password: "PASSWORD123",
			want:     false,
		},
		{
			name:     "no digit",
			password: "PasswordABC",
			want:     false,
		},
		{
			name:     "only lowercase",
			password: "password",
			want:     false,
		},
		{
			name:     "only uppercase",
			password: "PASSWORD",
			want:     false,
		},
		{
			name:     "only digits",
			password: "12345678",
			want:     false,
		},
		{
			name:     "special characters but missing requirements",
			password: "pass!@#$",
			want:     false,
		},
		{
			name:     "valid with special characters",
			password: "Pass123!@#",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidatePassword(tt.password); got != tt.want {
				t.Errorf("ValidatePassword(%q) = %v, want %v", tt.password, got, tt.want)
			}
		})
	}
}

func TestComputePHC(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "MyPassword123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false,
		},
		{
			name:     "long password",
			password: "ThisIsAVeryLongPasswordThatShouldStillWork123!",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ComputePHC(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ComputePHC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Verify PHC format
				if !strings.HasPrefix(got, "$argon2id$") {
					t.Errorf("ComputePHC() PHC string should start with $argon2id$, got %s", got)
				}
				parts := strings.Split(got, "$")
				if len(parts) != 6 {
					t.Errorf("ComputePHC() PHC string should have 6 parts, got %d", len(parts))
				}
			}
		})
	}
}

func TestComputePHC_Uniqueness(t *testing.T) {
	// Test that same password generates different hashes (due to random salt)
	password := "TestPassword123"
	phc1, err1 := ComputePHC(password)
	phc2, err2 := ComputePHC(password)

	if err1 != nil || err2 != nil {
		t.Fatalf("ComputePHC() failed: err1=%v, err2=%v", err1, err2)
	}

	if phc1 == phc2 {
		t.Error("ComputePHC() should generate different hashes for same password due to random salt")
	}
}

func TestVerifyPassword(t *testing.T) {
	// First generate a PHC hash for testing
	password := "MyPassword123"
	phc, err := ComputePHC(password)
	if err != nil {
		t.Fatalf("ComputePHC() failed: %v", err)
	}

	tests := []struct {
		name     string
		password string
		phc      string
		want     bool
		wantErr  bool
	}{
		{
			name:     "correct password",
			password: password,
			phc:      phc,
			want:     true,
			wantErr:  false,
		},
		{
			name:     "incorrect password",
			password: "WrongPassword123",
			phc:      phc,
			want:     false,
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			phc:      phc,
			want:     false,
			wantErr:  false,
		},
		{
			name:     "invalid PHC format - too few parts",
			password: password,
			phc:      "$argon2id$v=19$m=65536",
			want:     false,
			wantErr:  true,
		},
		{
			name:     "invalid PHC format - malformed parameters",
			password: password,
			phc:      "$argon2id$v=19$invalid$salt$hash",
			want:     false,
			wantErr:  true,
		},
		{
			name:     "invalid PHC format - invalid base64 salt",
			password: password,
			phc:      "$argon2id$v=19$m=65536,t=1,p=4$!!!invalid!!!$hash",
			want:     false,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := VerifyPassword(tt.password, tt.phc)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("VerifyPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerifyPassword_RoundTrip(t *testing.T) {
	// Test that we can hash and verify multiple passwords
	passwords := []string{
		"Password123",
		"AnotherPass456",
		"Test123ABC",
		"MyS3cureP@ss",
	}

	for _, password := range passwords {
		t.Run(password, func(t *testing.T) {
			phc, err := ComputePHC(password)
			if err != nil {
				t.Fatalf("ComputePHC() failed: %v", err)
			}

			verified, err := VerifyPassword(password, phc)
			if err != nil {
				t.Fatalf("VerifyPassword() failed: %v", err)
			}

			if !verified {
				t.Error("VerifyPassword() should verify the correct password")
			}

			// Try with wrong password
			wrongVerified, err := VerifyPassword(password+"wrong", phc)
			if err != nil {
				t.Fatalf("VerifyPassword() failed: %v", err)
			}

			if wrongVerified {
				t.Error("VerifyPassword() should not verify an incorrect password")
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		want     bool
	}{
		// Valid usernames
		{
			name:     "valid username with minimum length",
			username: "abc",
			want:     true,
		},
		{
			name:     "valid username with maximum length",
			username: "abcdefghij1234567890",
			want:     true,
		},
		{
			name:     "valid username with numbers",
			username: "user123",
			want:     true,
		},
		{
			name:     "valid username all uppercase",
			username: "USER",
			want:     true,
		},
		{
			name:     "valid username mixed case",
			username: "UserName123",
			want:     true,
		},
		// Invalid usernames - length
		{
			name:     "too short username",
			username: "ab",
			want:     false,
		},
		{
			name:     "empty username",
			username: "",
			want:     false,
		},
		{
			name:     "username too long",
			username: "abcdefghij12345678901",
			want:     false,
		},
		// Invalid usernames - format
		{
			name:     "starts with number",
			username: "123user",
			want:     false,
		},
		{
			name:     "contains special characters",
			username: "user_name",
			want:     false,
		},
		{
			name:     "contains spaces",
			username: "user name",
			want:     false,
		},
		{
			name:     "contains hyphen",
			username: "user-name",
			want:     false,
		},
		{
			name:     "contains dot",
			username: "user.name",
			want:     false,
		},
		{
			name:     "only numbers",
			username: "12345",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateUsername(tt.username); got != tt.want {
				t.Errorf("ValidateUsername(%q) = %v, want %v", tt.username, got, tt.want)
			}
		})
	}
}
