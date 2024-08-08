package controller

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/temuka-api-service/internal/model"
	"github.com/temuka-api-service/internal/repository"
	httputil "github.com/temuka-api-service/pkg/http"
	"golang.org/x/crypto/bcrypt"
)

type AuthController interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	ResetPassword(w http.ResponseWriter, r *http.Request)
}

type AuthControllerImpl struct {
	UserRepository repository.UserRepository
}

func NewAuthController(userRepository repository.UserRepository) AuthController {
	return &AuthControllerImpl{
		UserRepository: userRepository,
	}
}

func (c *AuthControllerImpl) Register(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := httputil.ReadRequest(r, &requestBody); err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
	if err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error hashing password"})
		return
	}

	newUser := model.User{
		Username:       requestBody.Username,
		Email:          requestBody.Email,
		Password:       string(hashedPwd),
		ProfilePicture: "",
		CoverPicture:   "",
	}

	if err := c.UserRepository.CreateUser(context.Background(), &newUser); err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error creating user"})
		return
	}

	response := struct {
		Message string     `json:"message"`
		Data    model.User `json:"data"`
	}{
		Message: "New user has been registered",
		Data:    newUser,
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}

func (c *AuthControllerImpl) Login(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := httputil.ReadRequest(r, &requestBody); err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	user, err := c.UserRepository.GetUserByEmail(context.Background(), requestBody.Email)
	if err != nil {
		httputil.WriteResponse(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password)); err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Wrong password"})
		return
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})

	tokenString, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error creating token"})
		return
	}

	response := struct {
		Message string `json:"message"`
		Token   string `json:"token"`
	}{
		Message: "User has login successfully",
		Token:   tokenString,
	}

	httputil.WriteResponse(w, http.StatusOK, response)
}

func (c *AuthControllerImpl) ResetPassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDstr := vars["id"]

	var requestBody struct {
		ResetToken              string `json:"reset_token"`
		Email                   string `json:"email"`
		NewPassword             string `json:"new_password"`
		NewPasswordConfirmation string `json:"new_password_confirmation"`
	}

	if err := httputil.ReadRequest(r, &requestBody); err != nil {
		httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	token, err := jwt.Parse(requestBody.ResetToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		httputil.WriteResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid or expired token"})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["email"] != requestBody.Email {
			httputil.WriteResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid token for the provided email"})
			return
		}

		if requestBody.NewPassword != requestBody.NewPasswordConfirmation {
			httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Password and password confirmation do not match"})
			return
		}
		hashedNewPwd, err := bcrypt.GenerateFromPassword([]byte(requestBody.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error hashing password"})
			return
		}

		userID, err := strconv.Atoi(userIDstr)
		if err != nil {
			httputil.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid user id"})
			return
		}

		user, err := c.UserRepository.GetUserByID(context.Background(), userID)
		if err != nil {
			httputil.WriteResponse(w, http.StatusNotFound, map[string]string{"error": "User not found"})
			return
		}

		user.Password = string(hashedNewPwd)
		if err := c.UserRepository.UpdateUser(context.Background(), userID, user); err != nil {
			httputil.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error updating new password"})
			return
		}

		response := struct {
			Message string `json:"message"`
		}{
			Message: "Password was reset successfully",
		}
		httputil.WriteResponse(w, http.StatusOK, response)
	} else {
		httputil.WriteResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
	}
}
