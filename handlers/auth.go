package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"html/template"
	"log"

	"go-login-app/database"
	"go-login-app/models"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	db        *database.DB
	templates *template.Template
	sessions  map[string]*models.User
}

func NewAuthHandler(db *database.DB, templates *template.Template) *AuthHandler {
	log.Println("Initializing AuthHandler")
	return &AuthHandler{
		db:        db,
		templates: templates,
		sessions:  make(map[string]*models.User),
	}
}

func (h *AuthHandler) HandleLogin(c *gin.Context) {
	log.Println("Handling login request")
	c.HTML(200, "login.html", gin.H{
		"error": c.Query("error"),
	})
}

func (h *AuthHandler) HandleLoginSubmit(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	
	log.Printf("Login attempt for user: %s", username)

	authenticated, err := h.db.AuthenticateUser(username, password)
	if err != nil {
		log.Printf("Authentication error: %v", err)
		c.Redirect(302, "/?error=Internal+server+error")
		return
	}

	if !authenticated {
		log.Printf("Failed login attempt for user: %s", username)
		c.Redirect(302, "/?error=Invalid+credentials")
		return
	}

	// Generate session token
	sessionToken, err := generateSessionToken()
	if err != nil {
		log.Printf("Session token generation error: %v", err)
		c.Redirect(302, "/?error=Internal+server+error")
		return
	}

	// Get user data and store in session
	user, err := h.db.GetUser(username)
	if err != nil {
		log.Printf("Error getting user data: %v", err)
		c.Redirect(302, "/?error=Internal+server+error")
		return
	}

	h.sessions[sessionToken] = user
	log.Printf("Successful login for user: %s", username)

	// Set session cookie
	c.SetCookie("session_token", sessionToken, 3600, "/", "", false, true)
	c.Redirect(302, "/dashboard")
}

func (h *AuthHandler) HandleDashboard(c *gin.Context) {
	user := h.getCurrentUser(c)
	if user == nil {
		log.Println("Unauthorized access attempt to dashboard")
		c.Redirect(302, "/?error=Please+login")
		return
	}

	log.Printf("Dashboard access by user: %s", user.Username)
	c.HTML(200, "dashboard.html", gin.H{
		"username": user.Username,
	})
}

func (h *AuthHandler) HandleLogout(c *gin.Context) {
	sessionToken, err := c.Cookie("session_token")
	if err == nil {
		if user := h.sessions[sessionToken]; user != nil {
			log.Printf("Logout for user: %s", user.Username)
		}
		delete(h.sessions, sessionToken)
	}
	c.SetCookie("session_token", "", -1, "/", "", false, true)
	c.Redirect(302, "/")
}

func (h *AuthHandler) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user := h.getCurrentUser(c); user == nil {
			log.Println("Authentication required - redirecting to login")
			c.Redirect(302, "/?error=Please+login")
			c.Abort()
			return
		}
		c.Next()
	}
}

func (h *AuthHandler) getCurrentUser(c *gin.Context) *models.User {
	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		return nil
	}
	return h.sessions[sessionToken]
}

func generateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
