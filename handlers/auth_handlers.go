package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/temuka-api-service/models"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), 10)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	newUser := models.User{
		Username:       requestBody.Username,
		Email:          requestBody.Email,
		Password:       string(hashedPwd),
		ProfilePicture: "",
		CoverPicture:   "",
	}

	if err := db.Create(&newUser).Error; err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	response := struct {
		Message string      `json:"message"`
		Data    models.User `json:"data"`
	}{
		Message: "New user has been registered",
		Data:    newUser,
	}
	json.NewEncoder(w).Encode(response)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := db.Where("email = ?", requestBody.Email).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password)); err != nil {
		http.Error(w, "Wrong password", http.StatusBadRequest)
		return
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})

	tokenString, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		http.Error(w, "Error creating token", http.StatusInternalServerError)
		return
	}

	response := struct {
		Message string `json:"message"`
		Token   string `json:"token"`
	}{
		Message: "User has login successfully",
		Token:   tokenString,
	}
	json.NewEncoder(w).Encode(response)
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDstr := vars["id"]

	var requestBody struct {
		ResetToken              string `json:"reset_token"`
		Email                   string `json:"email"`
		NewPassword             string `json:"new_password"`
		NewPasswordConfirmation string `json:"new_password_confirmation"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := jwt.Parse(requestBody.ResetToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["email"] != requestBody.Email {
			http.Error(w, "Invalid token for the provided email", http.StatusUnauthorized)
			return
		}

		if requestBody.NewPassword != requestBody.NewPasswordConfirmation {
			http.Error(w, "Password and password confirmation do not match", http.StatusBadRequest)
			return
		}
		hashedNewPwd, err := bcrypt.GenerateFromPassword([]byte(requestBody.NewPassword), 10)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		userID, err := strconv.Atoi(userIDstr)
		if err != nil {
			http.Error(w, "Invalid user id", http.StatusBadRequest)
			return
		}

		var user models.User
		res := db.First(&user, userID)
		if res.Error != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		user.Password = string(hashedNewPwd)
		if err := db.Save(&user).Error; err != nil {
			http.Error(w, "Error updating new password", http.StatusInternalServerError)
			return
		}

		response := struct {
			Message string `json:"message"`
		}{
			Message: "Password was reset successfully",
		}
		json.NewEncoder(w).Encode(response)
	}

}
