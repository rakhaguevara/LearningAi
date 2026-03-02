package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/adaptive-ai-learn/backend/config"
)

func main() {
	// Connect ke database
	config.ConnectDB()

	// Buat router
	r := gin.Default()

	// Test endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Backend jalan 🚀",
		})
	})

	// ======================
	// REGISTER
	// ======================
	r.POST("/api/register", func(c *gin.Context) {
		var input struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		query := "INSERT INTO users (name, email, password) VALUES ($1, $2, $3)"
		_, err := config.DB.Exec(query, input.Name, input.Email, input.Password)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"message": "Register berhasil & tersimpan 🎉",
		})
	})

	// ======================
	// LOGIN
	// ======================
	r.POST("/api/login", func(c *gin.Context) {
		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		var storedPassword string
		var name string

		query := "SELECT name, password FROM users WHERE email=$1"
		err := config.DB.QueryRow(query, input.Email).Scan(&name, &storedPassword)

		if err != nil {
			c.JSON(401, gin.H{"error": "Email tidak ditemukan"})
			return
		}

		if input.Password != storedPassword {
			c.JSON(401, gin.H{"error": "Password salah"})
			return
		}

		c.JSON(200, gin.H{
			"message": "Login berhasil 🎉",
			"name":    name,
		})
	})

	log.Println("Server running on port 8080")
	r.Run(":8080")
}