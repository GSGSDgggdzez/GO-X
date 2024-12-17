package main

import (
	"GO-X/controllers" // Import the controllers package where the database logic is handled
	"GO-X/routes"      // Import the routes package where the HTTP routes are defined
	"database/sql"     // Import the database/sql package to interact with the SQL database
	"log"              // Import the log package for logging errors and info

	_ "github.com/go-sql-driver/mysql" // Blank import to initialize the MySQL driver (this allows us to interact with MySQL databases)
	"github.com/gofiber/fiber/v2"      // Import the Fiber web framework for building the web server
)

func main() {
	// 1. Create a new Fiber app. This app will handle incoming HTTP requests and responses.
	app := fiber.New()

	// 2. Set up the MySQL database connection.
	// We are defining the Data Source Name (DSN) here, which contains the necessary credentials
	// to connect to our MySQL database (change these values to match your MySQL setup).
	dsn := "root:@tcp(localhost:3306)/GO-X" // Change this string to your MySQL username, password, and database name
	db, err := sql.Open("mysql", dsn)       // Attempt to open the database connection using the MySQL driver
	if err != nil {                         // If there's an error opening the database, we log and exit
		log.Fatal("Error opening the database: ", err)
	}
	defer db.Close() // Ensures that the database connection is closed when the function exits

	// 3. Ping the database to check if the connection is successful.
	// This is like saying "Hey, are you there?" to the database.
	err = db.Ping()
	if err != nil { // If the ping fails (database isn't reachable), log the error and stop the program
		log.Fatal("Error connecting to the database: ", err)
	}

	// 4. If the connection is successful, print a message.
	log.Println("Successfully connected to the MySQL database!")

	// 5. Now that the database is connected, we pass it to the controllers package.
	// This makes sure that the controllers have access to the database.
	controllers.SetDB(db)

	// 6. Next, we set up all the routes for the web application using the routes package.
	// Routes define how the app should handle incoming requests (like what happens when someone visits a URL).
	routes.SetupRoutes(app, db)

	// 7. Finally, start the server and listen for incoming HTTP requests.
	// The server will listen on port 8000, and handle requests as per the defined routes.
	err = app.Listen(":8000")
	if err != nil { // If there's an error starting the server, log it and stop the program
		log.Fatal("Error starting the server: ", err)
	}
}
