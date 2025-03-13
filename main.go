package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// Thanks to otraore for the code example
// https://gist.github.com/otraore/4b3120aa70e1c1aa33ba78e886bb54f3

const (
	userkey = "user"   // key used to store the username in the session
	secret  = "secret" // random and secure key used to encrypt the session cookie
)

func main() {
	// Initialize the engine
	e := engine()

	// Run the engine on port 8080
	if err := e.Run(":8080"); err != nil {
		log.Fatal("Unable to start:", err)
	}
}

func engine() *gin.Engine {
	// Create a new gin engine
	r := gin.New()

	// Enable Gin's logging middleware for request logging
	r.Use(gin.Logger())

	// Setup the cookie store for session management
	// This middleware will automatically handle session cookie reading/writing
	r.Use(sessions.Sessions("mysession", cookie.NewStore([]byte(secret))))

	// Public routes that don't require authentication
	r.POST("/login", login)  // Handles user login
	r.GET("/logout", logout) // Handles user logout

	// Private route group, protected by AuthRequired middleware
	// All routes within this group require a valid session
	private := r.Group("/private")
	private.Use(AuthRequired) // Enable the middleware on these routes
	{
		private.GET("/me", me)         // Returns current user info
		private.GET("/status", status) // Returns login status
	}

	return r
}

// AuthRequired is a middleware that checks if the user has a valid session.
// It should be used on routes that require authentication.
// If no valid session exists, it aborts the request with 401 Unauthorized.
func AuthRequired(c *gin.Context) {
	// Get the session from the request context
	session := sessions.Default(c)

	// Try to get the user from the session
	if user := session.Get(userkey); user == nil {
		// No user in session, abort the request
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// User is authenticated, continue to the next handler
	c.Next()
}

// login is a handler that parses a form and checks for specific data.
func login(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Validate form input
	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	// Check for username and password match, usually from a database
	if username != "hello" || password != "itsme" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	// Save the username in the session
	session.Set(userkey, username) // In real world usage you'd set this to the users ID
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user"})
}

// logout is the handler called for the user to log out.
func logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}
	session.Delete(userkey)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

// me is the handler that will return the user information stored in the
// session.
func me(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// status is the handler that will tell the user whether it is logged in or not.
func status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "You are logged in"})
}
