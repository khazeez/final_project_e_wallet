package pkg

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/KhoirulAziz99/final_project_e_wallet/internal/domain"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// type body struct {
// 	username string
// 	password string
// }

// type User struct {
// 	ID       int    `json:"id"`
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }

var jwtKey = []byte("SECRET_CODE")

func LoginHandler(c *gin.Context) {

	// var body struct {
	// 	Name     string
	// 	Password string
	// }

	// if c.Bind(&body) != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
	// 	return
	// }

	// var DB *gorm.DB
	var user domain.User
	// DB.First(&user, "name = ?", body.Name)

	// if user.ID == 0 {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user or password and id"})
	// 	return
	// }

	// if err := c.ShouldBind(&user); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	// , _ := FindOne(username)
	var password_tes = bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password))

	//logic authentication(compare username and password)
	if password_tes != nil {

		token := jwt.New(jwt.SigningMethodHS256)

		claims := token.Claims.(jwt.MapClaims)

		claims["username"] = user.Name
		claims["exp"] = time.Now().Add(time.Minute * 1).Unix()
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
		var user domain.User
		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Dapatkan pengguna dari database berdasarkan username
		query := "SELECT id, password FROM users WHERE username = $1"
		err = db.QueryRow(query, user.Name).Scan(&user.ID, &user.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		// Memeriksa kecocokan password
		password := c.Request.FormValue("password")
		err = bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		// Berhasil login, buat dan kirimkan token akses
		token := jwt.New(jwt.SigningMethodHS256)

		claims := token.Claims.(jwt.MapClaims)

		claims["username"] = user.Name
		claims["exp"] = time.Now().Add(time.Minute * 1).Unix()
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": tokenString})

	}
}

// func getUserByUsername(username string) *domain.User {
// 	// Mengirimkan kueri ke database untuk mendapatkan informasi pengguna berdasarkan username
// 	var db *sql.DB
// 	row := db.QueryRow("SELECT name, password FROM users WHERE name = $1", username)

// 	// Menginisialisasi variabel untuk menyimpan hasil kueri
// 	var user domain.User

// 	// Mengisi nilai-nilai pengguna dari hasil kueri
// 	err := row.Scan(&user.Name, &user.Password)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil // Username tidak ditemukan
// 		}
// 		panic(err)
// 	}

// 	return &user
// }

func FindOne(name string) (*domain.User, error) {
	var db *sql.DB
	query := `SELECT name, email, password, profile_picture, is_deleted FROM users WHERE name=$1;`
	row := db.QueryRow(query, name)
	var user domain.User
	err := row.Scan(&user.Name, &user.Email, &user.Password, &user.ProfilePicture, &user.IsDeleted)
	if err != nil {
		panic(err)
	}
	user.Name = name
	return &user, nil
}
