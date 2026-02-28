package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"mysocialai/model"
	"mysocialai/service"

	jwt "github.com/form3tech-oss/jwt-go"
)

func signinHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one signin request")
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var user model.User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
		fmt.Printf("Cannot decode user data from client: %v\n", err)
		return
	}

	success, err := service.CheckUser(user.Username, user.Password)
	if err != nil {
		http.Error(w, "Failed to read user from ElasticSearch", http.StatusInternalServerError)
		fmt.Printf("Failed to read user from ElasticSearch: %v\n", err)
		return
	}

	if !success {
		http.Error(w, "User does not exist or password is wrong", http.StatusUnauthorized)
		fmt.Printf("User does not exist or password is wrong: %s\n", user.Username)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		fmt.Printf("Failed to generate token: %v\n", err)
		return
	}
	w.Write([]byte(tokenString))
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one signup request")
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var user model.User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
		fmt.Printf("Cannot decode user data from client: %v\n", err)
		return
	}

	if user.Username == "" || user.Password == "" ||
		regexp.MustCompile(`^[a-z0-9]$`).MatchString(user.Username) {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		fmt.Printf("Invalid username or password: %s\n", user.Username)
		return
	}

	success, err := service.AddUser(&user)
	if err != nil {
		http.Error(w, "Failed to save user to ElasticSearch", http.StatusInternalServerError)
		fmt.Printf("Failed to save user to ElasticSearch: %v\n", err)
		return
	}

	if !success {
		http.Error(w, "User already exists", http.StatusConflict)
		fmt.Printf("User already exists: %s\n", user.Username)
		return
	}

	fmt.Printf("User is added: %s\n", user.Username)
	w.WriteHeader(http.StatusCreated)
}
