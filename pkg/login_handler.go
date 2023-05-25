package pkg

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("SECRET_CODE")

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorization"})
			c.Abort()
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) { return jwtKey, nil })

		if !token.Valid || err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorization"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("claims", claims)

		c.Next()
	}
}

type body struct {
	username string
	password string
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(c *gin.Context) {

	var user domain.User

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	userName := getUserByUsername(username)
	//logic authentication(compare username and password)
	if userName != nil && user.Password == password {
		//bikin code untuk generate token
		token := jwt.New(jwt.SigningMethodHS256)

		claims := token.Claims.(jwt.MapClaims) // ini map

		claims["username"] = user.Name
		claims["exp"] = time.Now().Add(time.Minute * 1).Unix() //token akan expired dalam 1 menit

		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}

func ProfileHandler(c *gin.Context) {
	// Ambil username dari JWT token
	claims := c.MustGet("claims").(jwt.MapClaims)
	username := claims["username"].(string)

	// seharusnya return user  dari database , tapi di contoh ini kita pakai code block
	c.JSON(http.StatusOK, gin.H{"username": username})
}

func LoginGPTHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("name")
		password := c.PostForm("password")

		// Cek kecocokan pengguna di database
		query := "SELECT id FROM users WHERE name = $1 AND password = $2"
		var userID int
		err := db.QueryRow(query, username, password).Scan(&userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		// Login berhasil
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Login successful for user with ID %d", userID)})
	}
}

func getUserByUsername(username string) *domain.User {
	// Mengirimkan kueri ke database untuk mendapatkan informasi pengguna berdasarkan username
	var db *sql.DB
	row := db.QueryRow("SELECT name, password FROM users WHERE name = $1", username)

	// Menginisialisasi variabel untuk menyimpan hasil kueri
	var user domain.User

	// Mengisi nilai-nilai pengguna dari hasil kueri
	err := row.Scan(&user.Name, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil // Username tidak ditemukan
		}
		panic(err)
	}

	return &user
}
