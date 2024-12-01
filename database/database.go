package database

import (
	"fmt"

	"go-login-app/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	*gorm.DB
}

// Initialize creates a new database connection and creates the user table if it doesn't exist
func Initialize(dbPath string) (*DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// Auto Migrate the schema
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

// CreateUser adds a new user to the database
func (db *DB) CreateUser(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		Username: username,
		Password: string(hashedPassword),
	}

	result := db.Create(&user)
	return result.Error
}

// GetUser retrieves a user from the database by username
func (db *DB) GetUser(username string) (*models.User, error) {
	var user models.User
	result := db.Where("username = ?", username).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// AuthenticateUser checks if the provided username and password are correct
func (db *DB) AuthenticateUser(username, password string) (bool, error) {
	user, err := db.GetUser(username)
	if err != nil {
		return false, fmt.Errorf("error getting user: %w", err)
	}
	if user == nil {
		return false, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, fmt.Errorf("error comparing passwords: %w", err)
	}
	return true, nil
}

// InitializeDefaultUser creates a default admin user if no users exist
func (db *DB) InitializeDefaultUser() error {
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count == 0 {
		return db.CreateUser("admin", "password123")
	}
	return nil
}
