package utils

import (
	"errors" // To handle errors
	"time"   // To work with time (e.g., expiration of the JWT)

	"github.com/golang-jwt/jwt/v4" // Import the JWT library to create and verify tokens
)

// Define a secret key for signing JWT tokens (use a more secure key in production)
var jwtSecretKey = []byte("your-secret-key")

// GenerateJWT generates a JWT token for the user
// It takes the username as input and returns a signed JWT token as a string.
func GenerateJWT(username string) (string, error) {
	// Create a new JWT token using the HMAC SHA-256 signing method
	token := jwt.New(jwt.SigningMethodHS256)

	// Set the claims (the data embedded inside the token)
	// Claims can hold any data you want to store in the token.
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username                         // Add the username to the claims
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Set the expiration time to 24 hours from now

	// Sign the token using the secret key
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		// If there is an error signing the token, return it
		return "", err
	}

	// Return the signed JWT token string
	return tokenString, nil
}

// ValidateJWT validates the JWT token
// It checks if the token is valid and returns the claims if it is.
func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	// Parse the token string and verify the token using the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure that the token is using the correct signing method (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method") // If the signing method is not HMAC, return an error
		}
		// Return the secret key used to sign the token so it can be validated
		return jwtSecretKey, nil
	})

	// If there was an error parsing the token, return the error
	if err != nil {
		return nil, err
	}

	// Check if the token is valid and extract the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// If the token is valid, return the claims (the data inside the token)
		return claims, nil
	}

	// If the token is invalid or expired, return an error
	return nil, errors.New("invalid or expired token")
}
