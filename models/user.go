package models

import (
	"database/sql" // Import the database/sql package to interact with SQL databases

	"golang.org/x/crypto/bcrypt" // Import bcrypt package for securely hashing passwords
)

// User struct represents a user in the system
// This is a Go struct that holds user information
// The struct tags `json:"username"` are used to specify how the struct fields should be named when converted to or from JSON
type User struct {
	ID       int    `json:"id"`       // The ID of the user, typically auto-generated in the database
	Username string `json:"username"` // The username of the user, unique in the system
	Email    string `json:"email"`    // The email address of the user
	Password string `json:"password"` // The password of the user (should be hashed before storing)
}

// Register a new user in the database
// This function is responsible for saving a new user's information into the "users" table in the database
func (u *User) Register(db *sql.DB) error {
	// The SQL query to insert the new user into the "users" table
	// It takes the username, email, and password from the User struct and inserts them into the table
	query := `INSERT INTO users (username, email, password) VALUES (?, ?, ?)`
	_, err := db.Exec(query, u.Username, u.Email, u.Password) // Execute the query
	if err != nil {
		// If there’s an error with the query (e.g., a database issue), return the error
		return err
	}
	// If successful, return nil (no error)
	return nil
}

// HashPassword hashes the user's password before saving it
// Passwords should never be saved in plain text, so we hash them before storing them
func HashPassword(password string) (string, error) {
	// bcrypt.GenerateFromPassword hashes the password using bcrypt with a default cost factor
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// If there’s an error while hashing, return an empty string and the error
		return "", err
	}
	// Return the hashed password as a string
	return string(hash), nil
}

// func CheckPasswordHash(Password string) (string, error) {

// 	HashChecks, err := becypt.CompareHashAndPassword([]byte(password), []byte(Password))

// 	if err != nil {
// 		return "", err
// 	}

// }

// GetUserByUsername retrieves a user by their username
// This function queries the database to find a user by their username
func GetUserByUsernameAndPassword(db *sql.DB, username string, password string) (*User, error) {
	var user User // Declare a User variable to hold the data from the database

	// Execute a SQL query to select the user's details from the "users" table based on the username
	err := db.QueryRow("SELECT id, username, email, password FROM users WHERE username = ?", username).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		// If no user is found (i.e., `sql.ErrNoRows`), return nil to indicate no user exists with that username
		// If there’s another error (e.g., database issue), return the error
		if err == sql.ErrNoRows {
			return nil, nil // No user found
		}
		return nil, err // Return any other database error
	}

	// Check if the provided password matches the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err // Passwords don't match
	}

	// Return the user found in the database
	return &user, nil
}
