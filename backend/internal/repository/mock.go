package repository

import (
	"context"
	"errors"

	"github.com/nutabi/cvwo-assignment/backend/internal/model"
	"gorm.io/gorm"
)

// MockRepository is a mock implementation of the Repository interface for testing
type MockRepository struct {
	// Storage for mock data
	Users map[uint]*model.User
	// Maps for quick lookups
	UsernameToID map[string]uint
	EmailToID    map[string]uint
	// Counter for auto-incrementing IDs
	NextID uint

	// Mock error responses (set these to simulate errors)
	MigrateError             error
	GetUserByIDError         error
	GetUserByUsernameError   error
	UpdateUserError          error
	CheckUsernameExistsError error
	CheckEmailExistsError    error
	CreateUserError          error
}

// NewMockRepository creates a new mock repository with initialized maps
func NewMockRepository() Repository {
	return &MockRepository{
		Users:        make(map[uint]*model.User),
		UsernameToID: make(map[string]uint),
		EmailToID:    make(map[string]uint),
		NextID:       1,
	}
}

// Migrate simulates database migration
func (m *MockRepository) Migrate() error {
	if m.MigrateError != nil {
		return m.MigrateError
	}
	return nil
}

// GetUserByID retrieves a user by ID
func (m *MockRepository) GetUserByID(ctx context.Context, id uint) (model.User, error) {
	if m.GetUserByIDError != nil {
		return model.User{}, m.GetUserByIDError
	}

	user, exists := m.Users[id]
	if !exists {
		return model.User{}, gorm.ErrRecordNotFound
	}

	return *user, nil
}

// GetUserByUsername retrieves a user by username
func (m *MockRepository) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	if m.GetUserByUsernameError != nil {
		return model.User{}, m.GetUserByUsernameError
	}

	userID, exists := m.UsernameToID[username]
	if !exists {
		return model.User{}, gorm.ErrRecordNotFound
	}

	user := m.Users[userID]
	return *user, nil
}

// UpdateUser updates an existing user
func (m *MockRepository) UpdateUser(ctx context.Context, user *model.User) error {
	if m.UpdateUserError != nil {
		return m.UpdateUserError
	}

	existingUser, exists := m.Users[user.ID]
	if !exists {
		return errors.New("user not found")
	}

	// Update fields (only AvatarURL and Bio as per the sql implementation)
	existingUser.AvatarURL = user.AvatarURL
	existingUser.Bio = user.Bio

	return nil
}

// CheckUsernameExists checks if a username already exists
func (m *MockRepository) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	if m.CheckUsernameExistsError != nil {
		return false, m.CheckUsernameExistsError
	}

	_, exists := m.UsernameToID[username]
	return exists, nil
}

// CheckEmailExists checks if an email already exists
func (m *MockRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	if m.CheckEmailExistsError != nil {
		return false, m.CheckEmailExistsError
	}

	_, exists := m.EmailToID[email]
	return exists, nil
}

// CreateUser creates a new user
func (m *MockRepository) CreateUser(ctx context.Context, user *model.User) error {
	if m.CreateUserError != nil {
		return m.CreateUserError
	}

	// Check if username already exists
	if _, exists := m.UsernameToID[user.Username]; exists {
		return errors.New("username already exists")
	}

	// Check if email already exists
	if _, exists := m.EmailToID[user.Email]; exists {
		return errors.New("email already exists")
	}

	// Assign ID if not set
	if user.ID == 0 {
		user.ID = m.NextID
		m.NextID++
	}

	// Store the user
	m.Users[user.ID] = user
	m.UsernameToID[user.Username] = user.ID
	m.EmailToID[user.Email] = user.ID

	return nil
}

// Helper methods for testing

// AddUser adds a user directly to the mock repository (useful for test setup)
func (m *MockRepository) AddUser(user *model.User) {
	if user.ID == 0 {
		user.ID = m.NextID
		m.NextID++
	}

	m.Users[user.ID] = user
	m.UsernameToID[user.Username] = user.ID
	m.EmailToID[user.Email] = user.ID
}

// Reset clears all data from the mock repository
func (m *MockRepository) Reset() {
	m.Users = make(map[uint]*model.User)
	m.UsernameToID = make(map[string]uint)
	m.EmailToID = make(map[string]uint)
	m.NextID = 1

	// Reset error fields
	m.MigrateError = nil
	m.GetUserByIDError = nil
	m.GetUserByUsernameError = nil
	m.UpdateUserError = nil
	m.CheckUsernameExistsError = nil
	m.CheckEmailExistsError = nil
	m.CreateUserError = nil
}
