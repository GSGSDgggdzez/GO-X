package routes

import (
	"GO-X/controllers" // Import the controllers package where the logic for handling user requests is defined
	"GO-X/middleware"  // Import the middleware package for adding additional functionality (e.g., security or authentication)
	"database/sql"     // Import the sql package to interact with the database

	"github.com/gofiber/fiber/v2" // Import the Fiber web framework to handle HTTP requests
)

// SetupRoutes sets up the routes and accepts the *sql.DB for database access
// This function defines all the HTTP routes that the application will handle and connects the database
func SetupRoutes(app *fiber.App, db *sql.DB) {
	// Post route for user registration
	// This route listens for POST requests to /auth/register and calls the RegisterUser function from the controllers package
	app.Post("/auth/register", controllers.RegisterUser)

	app.Post("/auth/login", controllers.LoginUser)

	// app.Post("/auth/forgot_password", controllers.forgot-password)

	// Route to check if the API is working
	// This route listens for GET requests to /api and sends a welcome message as a response
	app.Get("/api", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the Twitter Clone API!")
	})

	// Route to check the database connection status
	// This route listens for GET requests to /dbstatus and checks if the app can connect to the database
	app.Get("/dbstatus", func(c *fiber.Ctx) error {
		// Ping the database to check if the connection is live (it’s like asking the database, “Are you there?”)
		err := db.Ping()
		if err != nil { // If the database connection fails, we send an error response
			// Return a 500 Internal Server Error with the message "Database connection failed"
			return c.Status(500).SendString("Database connection failed: " + err.Error())
		}
		// If the connection is successful, return a success message with a 200 OK status
		return c.SendString("Successfully connected to the MySQL database!")
	})

	// Protected routes (require JWT)
	// This route listens for GET requests to /protected and checks if the user is authenticated using JWT (JSON Web Token)
	app.Get("/protected", middleware.ProtectRoute, func(c *fiber.Ctx) error {
		// Access the claims (user information) from the JWT token stored in the request context
		claims := c.Locals("claims").(map[string]interface{})
		// Return the username from the JWT claims in a JSON response
		return c.JSON(fiber.Map{
			"status": "success",
			"user":   claims["username"], // Send back the username stored in the JWT
		})
	})
}
