package controllers

import (
	"GO-X/models"  // Import the models package where we define and interact with the database models
	"GO-X/utils"   // Import the utils package for utility functions like generating JWT tokens
	"database/sql" // Import the sql package to interact with the SQL database
	"log"          // Import the log package to print error messages

	"github.com/go-playground/validator/v10" // Import Go validator package for input validation
	"github.com/gofiber/fiber/v2"            // Import the Fiber web framework to handle HTTP requests
)

var db *sql.DB // Declare a variable to store the database connection

// SetDB sets the database connection in the controllers package
// This function is called from the main app to initialize the db connection in the controllers
func SetDB(database *sql.DB) {
	db = database
}

// RegisterRequest struct defines the expected user registration data
// This structure represents the format of data we expect when a user registers
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"` // Validate that the username is required and between 3 and 50 characters
	Email    string `json:"email" validate:"required,email"`           // Validate that the email is required and in correct format
	Password string `json:"password" validate:"required,min=6,max=50"` // Validate that the password is required and between 6 and 50 characters
}

// RegisterUser handles the user registration
// This function is responsible for registering new users in the database
func RegisterUser(c *fiber.Ctx) error {
	// Check if the request body is empty
	// If the body is empty, return a 400 Bad Request with an error message
	if c.Body() == nil || len(c.Body()) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Request body cannot be empty", // Error message for empty body
			"errors":  "Invalid request format",       // Error description
		})
	}

	// Parse the incoming request body into the RegisterRequest struct
	// This converts the JSON body into a structured object so we can work with it
	var registerRequest RegisterRequest
	if err := c.BodyParser(&registerRequest); err != nil {
		log.Println("BodyParser error:", err)
		// If parsing fails, return a 422 Unprocessable Entity with the error
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "error",
			"message": "Review your input", // Generic error message for invalid input
			"errors":  err.Error(),         // Include the specific error for debugging
		})
	}

	// Validate the parsed input using the Go validator package
	// This checks if the input values (username, email, password) meet the required rules
	validate := validator.New()
	if err := validate.Struct(registerRequest); err != nil {
		// If validation fails, return a 400 Bad Request with the validation error message
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body", // Generic error message for validation failure
			"errors":  err.Error(),            // Return the detailed error message
		})
	}

	// Sanitize the inputs to prevent XSS or other malicious attacks
	// This ensures that the input doesn't contain harmful characters or scripts
	registerRequest.Username = sanitizeInput(registerRequest.Username)
	registerRequest.Email = sanitizeInput(registerRequest.Email)

	// Check if the user already exists by querying the database for the username
	existingUser, err := models.GetUserByUsernameAndPassword(db, registerRequest.Username, registerRequest.Password)
	if err != nil {
		log.Println("Error checking for existing user:", err)
		// If there’s an error checking the database, return a 500 Internal Server Error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal server error", // General error message for server issues
		})
	}
	if existingUser != nil {
		// If the username is already taken, return a 400 Bad Request with an error message
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Username already taken", // Message indicating the username is already in use
		})
	}

	// Hash the password before storing it in the database
	// This is a security measure to protect the user’s password from being stored in plain text
	hashedPassword, err := models.HashPassword(registerRequest.Password)
	if err != nil {
		log.Println("Error hashing password:", err)
		// If there’s an error hashing the password, return a 500 Internal Server Error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to process registration", // General error message for registration failure
		})
	}

	// Create a new User object with the sanitized input and hashed password
	user := models.User{
		Username: registerRequest.Username,
		Email:    registerRequest.Email,
		Password: hashedPassword,
	}

	// Register the new user in the database
	// This will save the user’s data in the database
	if err := user.Register(db); err != nil {
		log.Println("Error registering user:", err)
		// If there’s an error registering the user, return a 500 Internal Server Error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to register user", // Message for any issue during registration
		})
	}

	// Generate a JWT token after successful registration
	// This token will be used to authenticate the user in future requests
	token, err := utils.GenerateJWT(registerRequest.Username)
	if err != nil {
		log.Println("Error generating JWT:", err)
		// If there’s an error generating the token, return a 500 Internal Server Error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to generate JWT", // Message indicating a failure to create the token
		})
	}

	// Return the JWT token along with a success message
	// This response lets the user know they’ve successfully registered and now have a token
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "User registered successfully", // Success message for registration
		"token":   token,                          // Send the generated JWT token to the client
	})
}

// sanitizeInput sanitizes inputs to prevent XSS or other malicious attacks
// This function should remove unwanted characters or sanitize the input before it’s saved to the database
func sanitizeInput(input string) string {
	// Basic sanitization: You can expand this to handle more advanced sanitization
	// For example, you can use a library to escape HTML or remove scripts
	return input
}
