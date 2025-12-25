package service

import (
	"fmt"
)

func ErrUserNotFound(id uint) error {
	return fmt.Errorf("user with ID %v not found", id)
}

func ErrDatabaseUnknown(err error) error {
	return fmt.Errorf("unknown database error: %v", err)
}
