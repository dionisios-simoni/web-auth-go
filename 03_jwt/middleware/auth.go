package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/web-auth-go/03_jwt/initialisers"
	"github.com/web-auth-go/03_jwt/models"
	"log"
	"net/http"
	"os"
	"time"
)

func Auth(c *gin.Context) {
	tokenString, err := c.Cookie("authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var u models.User
		initialisers.DB.First(&u, claims["sub"])
		if u.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user", u)

		c.Next()

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
