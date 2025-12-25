package utility

func ValidatePasswordStrength(password string) bool {
	return true
}

func ComputePHC(password string) (string, error) {
	return password, nil
}

func VerifyPassword(password, phc string) bool {
	return password == phc
}
