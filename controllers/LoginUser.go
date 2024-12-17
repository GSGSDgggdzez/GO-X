package controllers

import (
	"GO-X/models"
	"GO-X/utils"
	"log"

	"github.com/gofiber/fiber/v2"
)

// LoginRequest struct defines the expected login data
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginUser handles the user login process
func LoginUser(c *fiber.Ctx) error {
	// Check if the request body is empty
	if c.Body() == nil || len(c.Body()) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Request body cannot be empty",
		})
	}

	// Parse the incoming request body into LoginRequest struct
	var loginRequest LoginRequest
	if err := c.BodyParser(&loginRequest); err != nil {
		log.Println("BodyParser error:", err)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input format",
			"errors":  err.Error(),
		})
	}

	// Validate input
	if loginRequest.Username == "" || loginRequest.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Username and password are required",
		})
	}

	// Retrieve the user from the database
	user, err := models.GetUserByUsernameAndPassword(db, loginRequest.Username, loginRequest.Password)
	if err != nil {
		log.Println("Error fetching user:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to authenticate user",
			"Error":   err.Error(),
		})
	}

	// Check if user exists and password matches
	// if user == nil || !models.CheckPasswordHash(loginRequest.Password, user.Password) {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
	// 		"status":  "error",
	// 		"message": "Invalid credentials",
	// 	})
	// }

	// Generate JWT token
	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		log.Println("Error generating JWT:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to generate token",
			"Error":   err.Error(),
		})
	}

	// Return the token to the user
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Login successful",
		"token":   token,
	})
}
