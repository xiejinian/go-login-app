package main

import (
	"html/template"
	"log"

	"github.com/gin-gonic/gin"
	"go-login-app/database"
	"go-login-app/handlers"
)

var (
	templates *template.Template
	db        *database.DB
)

func init() {
	templates = template.Must(template.ParseGlob("templates/*.html"))

	// Initialize database
	var err error
	db, err = database.Initialize("users.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create default admin user if no users exist
	if err := db.InitializeDefaultUser(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Create Gin router
	router := gin.Default()

	// Load HTML templates
	router.LoadHTMLGlob("templates/*")

	// Serve static files
	router.Static("/static", "./static")

	// Initialize auth handler
	authHandler := handlers.NewAuthHandler(db, templates)

	// Public routes
	router.GET("/", authHandler.HandleLogin)
	router.POST("/login", authHandler.HandleLoginSubmit)

	// Protected routes
	authorized := router.Group("/")
	authorized.Use(authHandler.RequireAuth())
	{
		authorized.GET("/dashboard", authHandler.HandleDashboard)
		authorized.GET("/logout", authHandler.HandleLogout)
	}

	log.Println("Server starting on http://localhost:8080")
	router.Run(":8080")
}
