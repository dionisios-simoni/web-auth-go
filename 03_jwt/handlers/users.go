package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/web-auth-go/03_jwt/initialisers"
	"github.com/web-auth-go/03_jwt/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

func Signup(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "could not hash password",
		})
		return
	}

	u := models.User{
		Email:    body.Email,
		Password: string(hashed),
	}

	result := initialisers.DB.Create(&u)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "could not create user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "user created"})
}

func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}

	var u models.User
	initialisers.DB.First(&u, "email = ?", body.Email)
	if u.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "could not find user",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": u.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	s, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("authorization", s, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"success": "logged in"})
}

func Validate(c *gin.Context) {
	u, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"user": u})
}
