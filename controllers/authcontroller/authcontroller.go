package authcontroller

import (
	"encoding/json"
	"golang-jwt-practice/config"
	"golang-jwt-practice/helpers"
	"golang-jwt-practice/models"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(w http.ResponseWriter, r *http.Request) {
	// Ambil inputan JSON
	var userInput models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helpers.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// Cek user ke DB berdasarkan username
	var user models.User
	if err := models.DB.Where("username = ?", userInput.Username).First(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			response := map[string]string{"message": "Username atau password salah"}
			helpers.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		default:
			response := map[string]string{"message": err.Error()}
			helpers.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}
	}

	// Cek Password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
		response := map[string]string{"message": "Username atau password salah"}
		helpers.ResponseJSON(w, http.StatusUnauthorized, response)
		return
	}

	// proses pembuatan token JWT
	expTime := time.Now().Add(time.Minute * 1)
	claims := &config.JWTClaim{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "golang-jwt-practice",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	// Declare algoritma untuk signin
	tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign token
	token, err := tokenAlgo.SignedString(config.JWT_KEY)
	if err != nil {
		response := map[string]string{"message": err.Error()}
		helpers.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	// Set token ke cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    token,
		HttpOnly: true,
	})

	response := map[string]string{"message": "Login Success"}
	helpers.ResponseJSON(w, http.StatusOK, response)
}

func Register(w http.ResponseWriter, r *http.Request) {

	// Ambil inputan JSON
	var userInput models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helpers.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}
	defer r.Body.Close()

	// Hash pass menggunakan bcrypt
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	userInput.Password = string(hashPassword)

	// Insert to DB
	if err := models.DB.Create(&userInput).Error; err != nil {
		log.Fatal("Failed to create user")
	}

	response := map[string]string{"message": "Success"}
	helpers.ResponseJSON(w, http.StatusOK, response)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// delete token di cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})

	response := map[string]string{"message": "Logout Success"}
	helpers.ResponseJSON(w, http.StatusOK, response)
}
