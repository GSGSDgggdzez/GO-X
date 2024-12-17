package controller

import (
	"GO-X/models"
	"GO-X/utils"
	"database/sql"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var db *sgl.dB

func SetDB(database *sql.DB) {
	db = Database
}

func RemoveUse(c *fiber.Ctx) error {

	type RemoveRequest struct {
		Username string `json:"username" validate:"required,min=3, max=50"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6, max=50"`

		Jwt string `json:"jwt"`
	}

	if c.body() == nill || lin(c.Body()) == 0 {
		log.Println("BodyParser error:", err)

		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "Error",
			"Massage": "Problem with Input",
			"error":   err.Error(),
		})
	}

	validate := validator.New()
	if err := validate.Struct(registerRequest); err != nil {
		// If validation fails, return a 400 Bad Request with the validation error message
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body", // Generic error message for validation failure
			"errors":  err.Error(),            // Return the detailed error message
		})
	}

	RemoveRequest.Username = sanitizeInput(RemoveRequest.Username)
	RemoveRequest.Email = sanitizeInput(RemoveRequest.Email)
	RemoveRequest.Password = sanitizeInput(RemoveRequest.Password)
	RemoveRequest.JWT = sanitizeInput(RemoveRequest.JWT)

	existingUser, err := models.GetUserByUsernameAndPassword(db, RemoveRequest.username, RemoveRequest.Password)

	if err != nil {
		log.Println("Error checking for existing user:", err)
		// If there’s an error checking the database, return a 500 Internal Server Error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal server error", // General error message for server issues
			"Error":   err.Error(),
		})
	}
	if existingUser != nil {
		// If the username is already taken, return a 400 Bad Request with an error message
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Username already taken", // Message indicating the username is already in use
			"Error":   err.Error(),
		})
	}

	ValidJWT, err := utils.ValidateJWT(RemoveRequest.JWT)

	log.Println("Error Problem with the JWT:", err)

	if err != nill {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Status":  "error",
			"message": "JWT Invalid",
			"Error":   err.Error(),
		})
	}

	if ok != nill {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Status":  "error",
			"message": "JWT Invalid",
			"Error":   err.Error(),
		})
	}

	return nill, error.New(" impossible to delete the user data")

}
