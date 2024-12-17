package middleware

import (
	"GO-X/utils" // Import utility functions (such as ValidateJWT)
	"log"        // For logging errors or other information
	"strings"    // For manipulating strings (e.g., trimming prefixes)

	"github.com/gofiber/fiber/v2" // Import the Fiber web framework
)

// ProtectRoute is a middleware function that protects routes by verifying the JWT token
// This function is called before any protected route handler
func ProtectRoute(c *fiber.Ctx) error {
	// Get the "Authorization" header from the incoming request
	authHeader := c.Get("Authorization")

	// If no Authorization header is present, return a 401 Unauthorized error
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Missing authorization token",
		})
	}

	// Check if the Authorization header starts with "Bearer "
	// This is the standard way to send tokens in the "Authorization" header
	if !strings.HasPrefix(authHeader, "Bearer ") {
		// If it doesn't, return a 401 Unauthorized error with an appropriate message
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid authorization token format",
		})
	}

	// Extract the actual JWT token by removing the "Bearer " prefix
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Call the utility function to validate the JWT token
	// This function checks if the token is valid and hasn't expired
	claims, err := utils.ValidateJWT(tokenString)
	if err != nil {
		// If there was an error (e.g., the token is invalid or expired), log the error and return a 401 Unauthorized response
		log.Println("Error validating token:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid or expired token",
		})
	}

	// If the token is valid, store the claims (user info and other data) in the Fiber context
	// This allows access to the claims in any handler that follows
	c.Locals("claims", claims)

	// Proceed with the request by calling the next handler in the chain
	return c.Next()
}
