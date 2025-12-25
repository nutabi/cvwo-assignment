package utility

import "strconv"

func ValidatePasswordStrength(password string) bool {
	return true
}

func ComputePHC(password string) (string, error) {
	return password, nil
}

func VerifyPassword(password, phc string) bool {
	return password == phc
}

func ValidateUsername(username string) bool {
	return true
}

func TryParseUserId(idStr string) (uint, bool) {
	num, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, false
	}
	return uint(num), true
}

func ValidateEmail(email string) bool {
	return true
}

func ValidateAvatarUrl(avatarUrl string) bool {
	return true
}

func ValidateBio(bio string) bool {
	return true
}
